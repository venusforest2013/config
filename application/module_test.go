package application

import (
	"context"
	"fmt"
	"os"
	"time"
)

type testRPC struct {
	sleep time.Duration
}

func (m *testRPC) Name() string {
	return "test-rpc"
}

func (m *testRPC) Configure(ctx context.Context, cfg *Config) error {
	m.sleep = cfg.Duration("sleep", 1*time.Second)
	return nil
}

func (m *testRPC) Start(ctx context.Context) error {
	ch := make(chan error, 1)
	go func() {
		time.Sleep(m.sleep)
		fmt.Fprintf(os.Stderr, "%s started\n", m.Name())
		ch <- nil
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (m *testRPC) Shutdown(ctx context.Context) error {
	fmt.Fprintf(os.Stderr, "%s done\n", m.Name())
	return nil
}

func (m *testRPC) Stats() map[string]interface{} {
	return nil
}

type testEngine struct {
}

func (m *testEngine) Name() string {
	return "test-engine"
}

func (m *testEngine) Configure(ctx context.Context, cfg *Config) error {
	return nil
}

func (m *testEngine) Start(ctx context.Context) error {
	fmt.Fprintf(os.Stderr, "%s started\n", m.Name())
	return nil
}

func (m *testEngine) Shutdown(ctx context.Context) error {
	fmt.Fprintf(os.Stderr, "%s done\n", m.Name())
	return nil
}

func (m *testEngine) Stats() map[string]interface{} {
	return nil
}
