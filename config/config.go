package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App       AppConfig       `toml:"app"`
	QQBot     QQBotConfig     `toml:"qqbot"`
	Server    ServerConfig    `toml:"server"`
	Dispatch  DispatchConfig  `toml:"dispatch"`
	Runtime   RuntimeConfig   `toml:"runtime"`
	Collector CollectorConfig `toml:"collector"`
	Targets   TargetsConfig   `toml:"targets"`
}

type AppConfig struct {
	LogLevel string `toml:"log_level"`
}

type QQBotConfig struct {
	AppID         string `toml:"app_id"`
	ClientSecret  string `toml:"client_secret"`
	TokenURL      string `toml:"token_url"`
	APIBase       string `toml:"api_base"`
	RequestTimout string `toml:"request_timeout"`
	Markdown      bool   `toml:"markdown"`
}

type ServerConfig struct {
	ListenAddr      string `toml:"listen_addr"`
	APIToken        string `toml:"api_token"`
	ReadTimeout     string `toml:"read_timeout"`
	WriteTimeout    string `toml:"write_timeout"`
	ShutdownTimeout string `toml:"shutdown_timeout"`
	MaxBodyBytes    int64  `toml:"max_body_bytes"`
}

type DispatchConfig struct {
	QueueSize      int    `toml:"queue_size"`
	Workers        int    `toml:"workers"`
	RetryMax       int    `toml:"retry_max"`
	RetryBackoffMS int    `toml:"retry_backoff_ms"`
	EnqueueTimeout string `toml:"enqueue_timeout"`
}

type RuntimeConfig struct {
	StatusTTLSeconds      int `toml:"status_ttl_seconds"`
	IdempotencyTTLSeconds int `toml:"idempotency_ttl_seconds"`
	CleanupIntervalSec    int `toml:"cleanup_interval_seconds"`
}

type CollectorConfig struct {
	ReconnectDelay string `toml:"reconnect_delay"`
}

type TargetsConfig struct {
	FilePath      string `toml:"file_path"`
	MaxRecords    int    `toml:"max_records"`
	FlushInterval string `toml:"flush_interval"`
}

func Default() Config {
	return Config{
		App: AppConfig{
			LogLevel: "info",
		},
		QQBot: QQBotConfig{
			TokenURL:      "https://bots.qq.com/app/getAppAccessToken",
			APIBase:       "https://api.sgroup.qq.com",
			RequestTimout: "10s",
			Markdown:      false,
		},
		Server: ServerConfig{
			ListenAddr:      ":8080",
			ReadTimeout:     "10s",
			WriteTimeout:    "15s",
			ShutdownTimeout: "10s",
			MaxBodyBytes:    1 << 20,
		},
		Dispatch: DispatchConfig{
			QueueSize:      2048,
			Workers:        8,
			RetryMax:       3,
			RetryBackoffMS: 200,
			EnqueueTimeout: "3s",
		},
		Runtime: RuntimeConfig{
			StatusTTLSeconds:      86400,
			IdempotencyTTLSeconds: 3600,
			CleanupIntervalSec:    60,
		},
		Collector: CollectorConfig{
			ReconnectDelay: "3s",
		},
		Targets: TargetsConfig{
			FilePath:      "data/targets.json",
			MaxRecords:    1000,
			FlushInterval: "2s",
		},
	}
}

func Load(path string) (Config, error) {
	return loadWithValidator(path, Config.ValidateForServer)
}

func LoadSDK(path string) (Config, error) {
	return loadWithValidator(path, Config.ValidateForSDK)
}

func loadWithValidator(path string, validator func(Config) error) (Config, error) {
	cfg := Default()
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, err
	}
	if err := validator(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) ValidateForSDK() error {
	return c.validateCommon()
}

func (c Config) ValidateForServer() error {
	if err := c.validateCommon(); err != nil {
		return err
	}
	if strings.TrimSpace(c.Server.ListenAddr) == "" {
		return errors.New("server.listen_addr 不能为空")
	}
	if c.Server.MaxBodyBytes <= 0 {
		return errors.New("server.max_body_bytes 必须大于 0")
	}
	if _, err := c.ReadTimeout(); err != nil {
		return fmt.Errorf("server.read_timeout 无效: %w", err)
	}
	if _, err := c.WriteTimeout(); err != nil {
		return fmt.Errorf("server.write_timeout 无效: %w", err)
	}
	if _, err := c.ShutdownTimeout(); err != nil {
		return fmt.Errorf("server.shutdown_timeout 无效: %w", err)
	}
	return nil
}

func (c Config) validateCommon() error {
	if strings.TrimSpace(c.QQBot.AppID) == "" {
		return errors.New("qqbot.app_id 不能为空")
	}
	if strings.TrimSpace(c.QQBot.ClientSecret) == "" {
		return errors.New("qqbot.client_secret 不能为空")
	}
	if c.Dispatch.QueueSize <= 0 {
		return errors.New("dispatch.queue_size 必须大于 0")
	}
	if c.Dispatch.Workers <= 0 {
		return errors.New("dispatch.workers 必须大于 0")
	}
	if c.Dispatch.RetryMax < 0 {
		return errors.New("dispatch.retry_max 不能小于 0")
	}
	if c.Dispatch.RetryBackoffMS <= 0 {
		return errors.New("dispatch.retry_backoff_ms 必须大于 0")
	}
	if _, err := c.RequestTimeout(); err != nil {
		return fmt.Errorf("qqbot.request_timeout 无效: %w", err)
	}
	if _, err := c.EnqueueTimeout(); err != nil {
		return fmt.Errorf("dispatch.enqueue_timeout 无效: %w", err)
	}
	if _, err := c.CollectorReconnectDelay(); err != nil {
		return fmt.Errorf("collector.reconnect_delay 无效: %w", err)
	}
	if strings.TrimSpace(c.Targets.FilePath) == "" {
		return errors.New("targets.file_path 不能为空")
	}
	if c.Targets.MaxRecords <= 0 {
		return errors.New("targets.max_records 必须大于 0")
	}
	if _, err := c.TargetsFlushInterval(); err != nil {
		return fmt.Errorf("targets.flush_interval 无效: %w", err)
	}
	return nil
}

func (c Config) RequestTimeout() (time.Duration, error) {
	return time.ParseDuration(c.QQBot.RequestTimout)
}

func (c Config) ReadTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Server.ReadTimeout)
}

func (c Config) WriteTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Server.WriteTimeout)
}

func (c Config) ShutdownTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Server.ShutdownTimeout)
}

func (c Config) EnqueueTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Dispatch.EnqueueTimeout)
}

func (c Config) CollectorReconnectDelay() (time.Duration, error) {
	return time.ParseDuration(c.Collector.ReconnectDelay)
}

func (c Config) TargetsFlushInterval() (time.Duration, error) {
	return time.ParseDuration(c.Targets.FlushInterval)
}
