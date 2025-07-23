package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype/truetype"
	"github.com/segmentio/ksuid"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/slide"
)

const dataStorePrefix = "captcha:data:"

type CaptchaCheckData struct {
	Key  string `json:"key"`
	Dots string `json:"dots"`
}

type SlideCheckData struct {
	Key string `json:"key"`
	X   int    `json:"x"`
}

type CaptchaService struct {
	textCapt  click.Captcha
	slideCapt slide.Captcha
	redis     *redis.Client
}

func NewCaptchaService(redis *redis.Client) *CaptchaService {
	service := &CaptchaService{redis: redis}
	service.init()
	return service
}

func (service *CaptchaService) init() {
	service.initClick()
	service.initSlide()
}

func (service *CaptchaService) initClick() {
	// 初始化验证码服务
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
	)

	// 加载字体资源
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		logger.Errorf("加载字体资源失败: %v", err)
	}

	imgs, err := imagesv2.GetImages()
	if err != nil {
		logger.Errorf("加载图片资源失败: %v", err)
	}

	// 设置验证码资源
	builder.SetResources(
		// 使用默认的中文字符集
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
	)

	service.textCapt = builder.Make()
}

func (service *CaptchaService) initSlide() {
	builder := slide.NewBuilder(
		slide.WithImageSize(option.Size{Width: 310, Height: 200}),
		slide.WithRangeDeadZoneDirections([]slide.DeadZoneDirectionType{slide.DeadZoneDirectionTypeLeft}),
		slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	imgs, err := imagesv2.GetImages()
	if err != nil {
		logger.Errorf("加载图片资源失败: %v", err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		logger.Errorf("加载拼图资源失败: %v", err)
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	service.slideCapt = builder.Make()

}

// ----------------------------
// 公开方法
// ----------------------------
func (service *CaptchaService) Get(ctx context.Context) (interface{}, error) {
	// 生成验证码数据
	captData, err := service.textCapt.Generate()
	if err != nil {
		return nil, fmt.Errorf("生成验证码失败: %v", err)
	}

	// 获取验证码数据
	dotData := captData.GetData()
	if dotData == nil {
		return nil, errors.New("验证码数据获取失败")
	}

	// 获取主图和缩略图
	mImage := captData.GetMasterImage()
	tImage := captData.GetThumbImage()

	// 转换为 Base64 编码
	mBase64, err := mImage.ToBase64()
	if err != nil {
		return nil, fmt.Errorf("转换主图到 Base64 失败：%v", err)
	}

	tBase64, err := tImage.ToBase64()
	if err != nil {
		return nil, fmt.Errorf("转换缩略图到 Base64 失败：%v", err)
	}

	key, err := service.saveCaptchaData(ctx, dotData)
	if err != nil {
		logger.Errorf("保存验证码缓存失败：%v", err)
		return nil, err
	}

	// 构建返回数据
	return map[string]string{
		"key":   key,
		"image": mBase64,
		"thumb": tBase64,
	}, nil
}

func (service *CaptchaService) Check(ctx context.Context, data CaptchaCheckData) bool {

	if data.Key == "" || data.Dots == "" {
		return false
	}

	var dct map[int]*click.Dot
	err := service.getCaptchaData(ctx, data.Key, &dct)
	if err != nil {
		logger.Debug(err)
		return false
	}

	src := strings.Split(data.Dots, ",")
	if len(dct)*2 != len(src) {
		return false
	}

	for i := 0; i < len(dct); i++ {
		dot := dct[i]
		x, _ := strconv.Atoi(src[i*2])
		y, _ := strconv.Atoi(src[i*2+1])
		if !click.Validate(x, y, dot.X, dot.Y, dot.Width, dot.Height, 5) {
			return false
		}
	}

	return true
}

func (service *CaptchaService) SlideGet(ctx context.Context) (interface{}, error) {

	captData, err := service.slideCapt.Generate()
	if err != nil {
		return nil, fmt.Errorf("生成验证码失败: %v", err)
	}

	bgImg, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, fmt.Errorf("获取验证码图片失败: %v", err)
	}

	bkImg, err := captData.GetTileImage().ToBase64()
	if err != nil {
		return nil, fmt.Errorf("获取验证码拼图失败: %v", err)
	}

	blockData := captData.GetData()
	if blockData == nil {
		return nil, errors.New("验证码数据获取失败")
	}

	key, err := service.saveCaptchaData(ctx, blockData)
	if err != nil {
		logger.Errorf("保存验证码缓存失败：%v", err)
		return nil, err
	}

	return map[string]any{
		"key":   key,
		"bgImg": bgImg,
		"bkImg": bkImg,
		"y":     blockData.Y,
	}, nil
}

func (service *CaptchaService) SlideCheck(ctx context.Context, data SlideCheckData) bool {

	if data.Key == "" || data.X < 0 {
		return false
	}

	var dct *slide.Block
	err := service.getCaptchaData(ctx, data.Key, &dct)
	if err != nil {
		logger.Debug(err)
		return false
	}

	return slide.Validate(data.X, 0, dct.X, 0, 5)
}

// ----------------------------
// 私有方法
// ----------------------------
func (service *CaptchaService) saveCaptchaData(ctx context.Context, captchaData interface{}) (string, error) {
	key := GenUniqueId()

	dotsByte, _ := json.Marshal(captchaData)
	_, err := service.redis.Set(ctx, dataStorePrefix+key, dotsByte, time.Minute*5).Result()
	if err != nil {
		return "", fmt.Errorf("保存验证码缓存失败：%v", err)
	}
	return key, nil
}

func (service *CaptchaService) getCaptchaData(ctx context.Context, key string, dest interface{}) error {
	dotsByte, err := service.redis.Get(ctx, dataStorePrefix+key).Result()
	if err != nil {
		return fmt.Errorf("获取缓存失败: %w", err)
	}

	return json.Unmarshal([]byte(dotsByte), dest)
}

// ----------------------------
// 工具函数
// ----------------------------
func GenUniqueId() string {
	return ksuid.New().String()
}
