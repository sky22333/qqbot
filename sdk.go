package qqbot

import (
	"context"
	"log/slog"
	"os"

	"github.com/sky22333/qqbot/config"
	"github.com/sky22333/qqbot/internal/bootstrap"
)

type Client struct {
	components *bootstrap.Components
}

type ClientOptions struct {
	StartCollector bool
}

func New(cfg Config) (*Client, error) {
	return NewWithOptions(cfg, ClientOptions{
		StartCollector: true,
	})
}

func NewWithOptions(cfg Config, opts ClientOptions) (*Client, error) {
	if err := cfg.ValidateForSDK(); err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	components, err := bootstrap.New(cfg, logger, bootstrap.Options{StartCollector: opts.StartCollector})
	if err != nil {
		return nil, err
	}
	return &Client{components: components}, nil
}

func NewFromConfigFile(path string) (*Client, error) {
	cfg, err := config.LoadSDK(path)
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

func (c *Client) Send(ctx context.Context, req PushRequest) (PushResult, error) {
	return c.components.Notifier.Send(ctx, req)
}

func (c *Client) Enqueue(ctx context.Context, req PushRequest) (string, error) {
	return c.components.Notifier.Enqueue(ctx, req)
}

func (c *Client) GetStatus(requestID string) (DeliveryStatus, bool) {
	return c.components.Notifier.GetStatus(requestID)
}

func (c *Client) Close() {
	c.components.Close()
}
