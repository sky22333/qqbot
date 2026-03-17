package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSDKAllowsServerOnlyConstraintsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sdk.toml")
	content := `
[qqbot]
app_id = "123"
client_secret = "secret"

[server]
listen_addr = ""
max_body_bytes = 0
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	if _, err := LoadSDK(path); err != nil {
		t.Fatalf("LoadSDK 不应因 server 字段失败: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("Load 应校验 server 字段并返回错误")
	}
}
