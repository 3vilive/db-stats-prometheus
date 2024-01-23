package tracer

import (
	"errors"
	"time"
)

type Config struct {
	CheckInterval time.Duration `json:"check_internal"`
	Labels        []string      `json:"labels"`
}

func DefaultConfig() Config {
	return Config{
		CheckInterval: 5 * time.Second,
		Labels:        []string{"name"},
	}
}

func (c *Config) Check() error {
	if c.CheckInterval == 0 {
		return errors.New("invalid check interval")
	}

	return nil
}

type ApplyConfig func(*Config)

func WithCheckInterval(checkInterval time.Duration) ApplyConfig {
	return func(config *Config) {
		config.CheckInterval = checkInterval
	}
}

func WithLabels(labels ...string) ApplyConfig {
	return func(config *Config) {
		config.Labels = append(config.Labels, labels...)
	}
}
