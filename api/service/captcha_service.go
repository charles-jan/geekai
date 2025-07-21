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
	"geekai/core/types"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype/truetype"
	"github.com/imroc/req/v3"
	"github.com/segmentio/ksuid"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

type CaptchaService struct {
	config   types.ApiConfig
	client   *req.Client
	textCapt click.Captcha
	redis    *redis.Client
}

func NewCaptchaService(config types.ApiConfig, redis *redis.Client) *CaptchaService {
	service := &CaptchaService{
		config: config,
		client: req.C().SetTimeout(10 * time.Second),
		redis:  redis,
	}
	service.init()
	return service
}

func (service *CaptchaService) init() {
	// 初始化验证码服务
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
	)

	// 加载字体资源
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		log.Fatalln("加载字体资源失败:", err)
	}

	imgs, err := imagesv2.GetImages()
	if err != nil {
		log.Fatalln("加载图片资源失败:", err)
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

	key, err := service.saveDotData(ctx, dotData)
	if err != nil {
		return nil, err
	}

	// 构建返回数据
	return map[string]string{
		"key":   key,
		"image": mBase64,
		"thumb": tBase64,
	}, nil
}

const dotsStorePrefix = "captcha:dots:"

func (service *CaptchaService) saveDotData(ctx context.Context, dotData map[int]*click.Dot) (string, error) {
	key := GenUniqueId()

	dotsByte, _ := json.Marshal(dotData)
	_, err := service.redis.Set(ctx, dotsStorePrefix+key, dotsByte, time.Minute*5).Result()
	if err != nil {
		return "", fmt.Errorf("保存验证码缓存失败：%v", err)
	}
	return key, nil
}

type CaptchaCheckData struct {
	Key string `json:"key"`
	// 验证码点数据 x1,y1,x2,y2, ...
	Dots string `json:"dots"`
}

func (service *CaptchaService) Check(ctx context.Context, data CaptchaCheckData) bool {

	if data.Key == "" || data.Dots == "" {
		return false
	}

	dct, err := service.getDotData(ctx, data.Key)
	if err != nil {
		return false
	}
	logger.Debug(dct)

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

func (service *CaptchaService) getDotData(ctx context.Context, key string) (map[int]*click.Dot, error) {
	var dct map[int]*click.Dot
	dotsByte, err := service.redis.Get(ctx, dotsStorePrefix+key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(dotsByte), &dct)
	return dct, err
}

func (service *CaptchaService) SlideGet() (interface{}, error) {
	if service.config.Token == "" {
		return nil, errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s/api/captcha/slide/get", service.config.ApiURL)
	var res types.BizVo
	r, err := service.client.R().
		SetHeader("AppId", service.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", service.config.Token)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		return nil, fmt.Errorf("请求 API 失败：%v", err)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("请求 API 失败：%s", res.Message)
	}

	return res.Data, nil
}

func (service *CaptchaService) SlideCheck(data interface{}) bool {
	url := fmt.Sprintf("%s/api/captcha/slide/check", service.config.ApiURL)
	var res types.BizVo
	r, err := service.client.R().
		SetHeader("AppId", service.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", service.config.Token)).
		SetBodyJsonMarshal(data).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return false
	}

	if res.Code != types.Success {
		return false
	}

	return true
}

func GenUniqueId() string {
	return ksuid.New().String()
}
