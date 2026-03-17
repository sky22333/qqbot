package qqbot

import (
	"path/filepath"
	"testing"

	"github.com/sky22333/qqbot/config"
)

func TestNewWithOptionsWithoutCollector(t *testing.T) {
	cfg := config.Default()
	cfg.QQBot.AppID = "123"
	cfg.QQBot.ClientSecret = "secret"
	cfg.Server.ListenAddr = ""
	cfg.Server.MaxBodyBytes = 0
	cfg.Targets.FilePath = filepath.Join(t.TempDir(), "targets.json")

	client, err := NewWithOptions(cfg, ClientOptions{StartCollector: false})
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	client.Close()
}
