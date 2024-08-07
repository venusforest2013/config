package application

import "context"

// Module is the interface implemented by an object
// that can configure/start/shutdown itself.
type Module interface {
	Name() string
	Configure(ctx context.Context, cfg *Config) error
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Stats() map[string]interface{}
}

type BaseModule struct {}

func (m *BaseModule) Name() string {
	return "UndefinedModule"
}

func (m *BaseModule) Configure(ctx context.Context, cfg *Config) error {
	return nil
}

func (m *BaseModule) Start(ctx context.Context) error {
	return nil
}

func (m *BaseModule) Shutdown(ctx context.Context) error {
	return nil
}

func (m *BaseModule) Stats() map[string]interface{} {
	return nil
}