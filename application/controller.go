package application

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

var (
	errModuleControlBreak = errors.New("break")
	defaultControlTimeout = 1 * time.Second
)

type module struct {
	Module
	enabled bool
}

type moduleHandler func(ctx context.Context, m *module) error

type Controller struct {
	ctx     context.Context
	cfg     *Config
	timeout time.Duration
	mutex   *sync.RWMutex
	modules []*module
	dict    map[string]*module
}

func NewController(ctx context.Context) *Controller {
	return &Controller{
		ctx:     ctx,
		mutex:   &sync.RWMutex{},
		modules: []*module{},
		dict:    map[string]*module{},
	}
}

func (c *Controller) Module(name string) Module {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	m, ok := c.dict[name]
	if !ok {
		return nil
	}
	return m.Module
}

func (c *Controller) Modules() []Module {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	modules := []Module{}
	for _, m := range c.modules {
		modules = append(modules, m.Module)
	}
	return modules
}

func (c *Controller) Has(name string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.dict[name]
	return ok
}

func (c *Controller) Enabled(name string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if m, ok := c.dict[name]; ok {
		return m.enabled
	}
	return false
}

func (c *Controller) Register(m Module) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.dict[m.Name()]; ok {
		return fmt.Errorf("Module '%s' already exists\n", m.Name())
	}
	entry := &module{m, true}
	c.modules = append(c.modules, entry)
	c.dict[m.Name()] = entry
	return nil
}

func (c *Controller) configure(cfg *Config) error {
	c.cfg = cfg
	c.timeout = cfg.Duration("control_timeout", defaultControlTimeout)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.forEach(func(ctx context.Context, m *module) error {
		cfg := c.moduleConfig(m.Name())
		if cfg == nil {
			m.enabled = false
			return nil
		}
		return m.Configure(ctx, cfg)
	}, false)
}

func (c *Controller) start() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.forEach(func(ctx context.Context, m *module) error {
		if !m.enabled {
			return nil
		}
		return m.Start(ctx)
	}, false)
}

func (c *Controller) shutdown() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.forEach(func(ctx context.Context, m *module) error {
		if m.enabled {
			m.Shutdown(ctx)
			m.enabled = false
		}
		return nil
	}, true)
}

func (c *Controller) forEach(handler moduleHandler, desc bool) error {
	modules := c.modules
	// 逆序关闭解决模块依赖问题
	if desc {
		modules = moduleReverse(modules)
	}
	for _, m := range modules {
		// 模块控制逻辑在单独线程异步执行，兜底崩溃的情况
		if err := c.exec(c.ctx, m, handler); err != nil {
			return fmt.Errorf("%s: %v", m.Name(), err)
		}
	}
	return nil
}

func (c *Controller) exec(ctx context.Context, m *module, run moduleHandler) error {
	var err error
	ch := make(chan error)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	go func(ctx context.Context, m *module) {
		defer func() {
			if e := recover(); e != nil {

				ch <- fmt.Errorf("panic: %v\nstack: %v", e, string(debug.Stack()))
			}
		}()
		ch <- run(ctx, m)
	}(ctx, m)
	select {
	case <-ctx.Done():
		return fmt.Errorf("%v, timeout = %v", ctx.Err(), c.timeout)
	case err = <-ch:
		close(ch)
	}
	return err
}

func (c *Controller) moduleConfig(name string) *Config {
	node := c.cfg.Access(name)
	if node != nil {
		return &Config{node}
	}
	return nil
}

func moduleReverse(m []*module) []*module {
	if m == nil || len(m) <= 1 {
		return m
	}
	modules := []*module{}
	for i := len(m) - 1; i >= 0; i-- {
		modules = append(modules, m[i])
	}
	return modules
}
