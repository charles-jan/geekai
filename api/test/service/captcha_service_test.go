package service

import (
	"context"
	"geekai/core/types"
	"geekai/service"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestCaptchaService_Get(t *testing.T) {

	// 初始化服务
	captchaService := service.NewCaptchaService(redis.NewClient())

	// 调用 Get 方法生成验证码
	result, err := captchaService.Get(context.Background())
	if err != nil {
		t.Errorf("Get() error = %v", err)
		return
	}

	// 检查返回结果
	if result == nil {
		t.Errorf("Get() returned nil result")
	} else {
		t.Logf("Get() success: %+v", result)
	}
}

func TestCaptchaService_GenUniqueId(t *testing.T) {
	t.Log(service.GenUniqueId())
}
