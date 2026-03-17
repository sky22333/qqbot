package notifier

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sky22333/qqbot/config"
	"github.com/sky22333/qqbot/internal/targets"
	"github.com/sky22333/qqbot/message"
)

func TestEnqueueFillTargetFromLatestByType(t *testing.T) {
	cfg := config.Default()
	cfg.Dispatch.Workers = 0
	cfg.Targets.FilePath = filepath.Join(t.TempDir(), "targets.json")
	store, err := targets.NewStore(cfg.Targets.FilePath, cfg.Targets.MaxRecords, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("创建目标存储失败: %v", err)
	}
	defer store.Close()
	if err := store.Upsert(message.TargetC2C, "user-001", "m1", "hello"); err != nil {
		t.Fatalf("写入目标失败: %v", err)
	}

	n, err := New(cfg, nil)
	if err != nil {
		t.Fatalf("创建 notifier 失败: %v", err)
	}
	defer n.Close()
	n.SetTargetStore(store)

	_, err = n.Enqueue(context.Background(), message.PushRequest{
		TargetType: message.TargetC2C,
		Content:    "通知",
	})
	if err != nil {
		t.Fatalf("应自动补全 target_id，但返回错误: %v", err)
	}
}

func TestEnqueueWithoutStoreReturnsTargetIDError(t *testing.T) {
	cfg := config.Default()
	cfg.Dispatch.Workers = 0

	n, err := New(cfg, nil)
	if err != nil {
		t.Fatalf("创建 notifier 失败: %v", err)
	}
	defer n.Close()

	_, err = n.Enqueue(context.Background(), message.PushRequest{
		TargetType: message.TargetC2C,
		Content:    "通知",
	})
	if err == nil {
		t.Fatalf("预期返回 target_id 不能为空")
	}
	if !strings.Contains(err.Error(), "target_id 不能为空") {
		t.Fatalf("错误信息不符合预期: %v", err)
	}
}
