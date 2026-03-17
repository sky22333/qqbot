package qqbot

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sky22333/qqbot/config"
)

func TestSendNotificationWithConfig(t *testing.T) {
	enabled := strings.EqualFold(os.Getenv("QQBOT_REAL_SEND_TEST"), "true") || os.Getenv("QQBOT_REAL_SEND_TEST") == "1"
	if !enabled {
		t.Skip("QQBOT_REAL_SEND_TEST 未开启，跳过真实发送测试")
	}

	configPath := os.Getenv("QQBOT_TEST_CONFIG")
	if strings.TrimSpace(configPath) == "" {
		configPath = "configs/config.toml"
	}

	cfg, err := config.LoadSDK(configPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if strings.Contains(cfg.QQBot.AppID, "你的") || strings.Contains(cfg.QQBot.ClientSecret, "你的") {
		t.Skip("请先替换 qqbot.app_id 和 qqbot.client_secret")
	}
	targetID := strings.TrimSpace(os.Getenv("QQBOT_TEST_TARGET_ID"))
	if targetID == "" {
		t.Skip("QQBOT_TEST_TARGET_ID 为空，跳过真实发送测试")
	}
	targetType := strings.TrimSpace(strings.ToLower(os.Getenv("QQBOT_TEST_TARGET_TYPE")))
	if targetType == "" {
		targetType = string(TargetC2C)
	}
	if targetType != string(TargetC2C) && targetType != string(TargetGroup) && targetType != string(TargetChannel) {
		t.Fatalf("QQBOT_TEST_TARGET_TYPE 无效: %s", targetType)
	}
	content := strings.TrimSpace(os.Getenv("QQBOT_TEST_CONTENT"))
	if content == "" {
		content = "来自qqbot真实发送测试的通知"
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("初始化客户端失败: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	result, err := client.Send(ctx, PushRequest{
		TargetType: TargetType(targetType),
		TargetID:   targetID,
		Content:    content,
	})
	if err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	if strings.TrimSpace(result.MessageID) == "" {
		t.Fatalf("发送成功但 message_id 为空: %+v", result)
	}
}
