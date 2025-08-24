package utils

import (
	"context"
	"fmt"
	"time"
)

type DelayConfig struct {
	Duration time.Duration
	Message  string
	Enabled  bool
}

func DefaultDelay() *DelayConfig {
	return &DelayConfig{
		Duration: 1 * time.Second,
		Enabled:  true,
	}
}

func CustomDelay(duration time.Duration) *DelayConfig {
	return &DelayConfig{
		Duration: duration,
		Enabled:  true,
	}
}

func DelayWithMessage(duration time.Duration, message string) *DelayConfig {
	return &DelayConfig{
		Duration: duration,
		Message:  message,
		Enabled:  true,
	}
}

func (d *DelayConfig) Execute() {
	if !d.Enabled {
		return
	}

	if d.Message != "" {
		fmt.Printf("⏳ %s (waiting %v)\n", d.Message, d.Duration)
	}

	time.Sleep(d.Duration)
}

func (d *DelayConfig) ExecuteWithContext(ctx context.Context) error {
	if !d.Enabled {
		return nil
	}

	if d.Message != "" {
		fmt.Printf("⏳ %s (waiting %v)\n", d.Message, d.Duration)
	}

	select {
	case <-time.After(d.Duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type DelayManager struct {
	configs map[string]*DelayConfig
	enabled bool
}

func NewDelayManager() *DelayManager {
	return &DelayManager{
		configs: make(map[string]*DelayConfig),
		enabled: true,
	}
}

func (dm *DelayManager) AddDelay(name string, config *DelayConfig) {
	dm.configs[name] = config
}

func (dm *DelayManager) SetAPIDelay() {
	dm.AddDelay("api", DelayWithMessage(1*time.Second, "Rate limiting API request"))
}

func (dm *DelayManager) SetDatabaseDelay() {
	dm.AddDelay("database", DelayWithMessage(500*time.Millisecond, "Rate limiting database operation"))
}

func (dm *DelayManager) Execute(name string) {
	if !dm.enabled {
		return
	}

	if config, exists := dm.configs[name]; exists {
		config.Execute()
	}
}

func (dm *DelayManager) ExecuteWithContext(ctx context.Context, name string) error {
	if !dm.enabled {
		return nil
	}

	if config, exists := dm.configs[name]; exists {
		return config.ExecuteWithContext(ctx)
	}

	return nil
}

func (dm *DelayManager) Disable() {
	dm.enabled = false
}

func (dm *DelayManager) Enable() {
	dm.enabled = true
}

func (dm *DelayManager) ExecuteSequential(names ...string) {
	for _, name := range names {
		dm.Execute(name)
	}
}

func (dm *DelayManager) ExecuteSequentialWithContext(ctx context.Context, names ...string) error {
	for _, name := range names {
		if err := dm.ExecuteWithContext(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

func RateLimitedOperation(operation func() error, delay *DelayConfig) error {
	err := operation()
	if delay != nil {
		delay.Execute()
	}
	return err
}

func RateLimitedOperationWithContext(ctx context.Context, operation func() error, delay *DelayConfig) error {
	err := operation()
	if delay != nil {
		if delayErr := delay.ExecuteWithContext(ctx); delayErr != nil {
			return delayErr
		}
	}
	return err
}
