package bootstrap

import (
	"log/slog"

	"github.com/sky22333/qqbot/config"
	"github.com/sky22333/qqbot/internal/collector"
	"github.com/sky22333/qqbot/internal/notifier"
	"github.com/sky22333/qqbot/internal/targets"
)

type Options struct {
	StartCollector bool
}

type Components struct {
	Notifier  *notifier.Notifier
	Targets   *targets.Store
	Collector *collector.Collector
}

func New(cfg config.Config, logger *slog.Logger, opts Options) (*Components, error) {
	flushInterval, err := cfg.TargetsFlushInterval()
	if err != nil {
		return nil, err
	}
	targetStore, err := targets.NewStore(cfg.Targets.FilePath, cfg.Targets.MaxRecords, flushInterval)
	if err != nil {
		return nil, err
	}
	n, err := notifier.New(cfg, logger)
	if err != nil {
		_ = targetStore.Close()
		return nil, err
	}
	n.SetTargetStore(targetStore)
	c := &Components{
		Notifier: n,
		Targets:  targetStore,
	}
	if !opts.StartCollector {
		return c, nil
	}
	targetCollector, err := collector.New(cfg, logger, targetStore)
	if err != nil {
		n.Close()
		_ = targetStore.Close()
		return nil, err
	}
	targetCollector.Start()
	c.Collector = targetCollector
	return c, nil
}

func (c *Components) Close() {
	if c == nil {
		return
	}
	if c.Collector != nil {
		c.Collector.Stop()
	}
	if c.Notifier != nil {
		c.Notifier.Close()
	}
	if c.Targets != nil {
		_ = c.Targets.Close()
	}
}
