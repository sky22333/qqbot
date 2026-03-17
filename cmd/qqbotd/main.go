package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sky22333/qqbot/config"
	"github.com/sky22333/qqbot/internal/bootstrap"
	"github.com/sky22333/qqbot/internal/httpserver"
)

func main() {
	configPath := flag.String("config", "configs/config.toml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	logLevel := slog.LevelInfo
	if cfg.App.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	components, err := bootstrap.New(cfg, logger, bootstrap.Options{StartCollector: true})
	if err != nil {
		panic(err)
	}
	logger.Info("采集器已启动")

	server, err := httpserver.New(cfg, logger, components.Notifier, components.Targets)
	if err != nil {
		panic(err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-stopCh:
		logger.Info("收到退出信号，停止程序", "signal", sig.String())
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("服务异常退出", "error", err)
		}
	}

	shutdownTimeout, err := cfg.ShutdownTimeout()
	if err != nil {
		shutdownTimeout = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	_ = server.Shutdown(ctx)
	components.Close()
	logger.Info("服务已退出")
}
