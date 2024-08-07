package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	app "gitlab.alibaba-inc.com/amap-aos-go/application"
)

var (
	defaultServerName              = "default"
	defaultServerAddr              = ":29292"
	defaultServerReadTimeout       = 1 * time.Second
	defaultServerReadHeaderTimeout = 1 * time.Second
	defaultServerWriteTimeout      = 1 * time.Second
	defaultServerIdleTimeout       = 0 * time.Second
	defaultServerMaxHeaderBytes    = int64(0)
)

func init() {
	app.Register(newModule())
}

type Module struct {
	config   *app.Config
	handlers map[string]http.Handler
	servers  map[string]*http.Server
}

func newModule() *Module {
	return &Module{
		handlers: map[string]http.Handler{},
		servers:  map[string]*http.Server{},
	}
}

func (m *Module) RegisterHandler(name string, handler http.Handler) {
	m.handlers[name] = handler
}

func (m *Module) Name() string {
	return "http"
}

func (m *Module) Configure(ctx context.Context, cfg *app.Config) error {
	serversConf := cfg.ValueArray("servers")
	if serversConf == nil {
		return errors.New("http.servers not found")
	}
	for _, conf := range serversConf {
		name := conf.Str("name", defaultServerName)
		server := &http.Server{
			Addr:              conf.Str("addr", defaultServerAddr),
			ReadTimeout:       conf.Duration("read_timeout", defaultServerReadTimeout),
			ReadHeaderTimeout: conf.Duration("read_header_timeout", defaultServerReadHeaderTimeout),
			WriteTimeout:      conf.Duration("write_timeout", defaultServerWriteTimeout),
			IdleTimeout:       conf.Duration("idle_timeout", defaultServerIdleTimeout),
			MaxHeaderBytes:    int(conf.Int64("max_header_bytes", defaultServerMaxHeaderBytes)),
		}
		if _, ok := m.servers[name]; ok {
			return fmt.Errorf("http.server '%s' already exists", name)
		}
		m.servers[name] = server
	}
	m.config = cfg
	return nil
}

func (m *Module) Start(ctx context.Context) error {
	for name, server := range m.servers {
		var err error
		ch := make(chan error)
		go func() {
			var err error
			defer func() {
				if e := recover(); e != nil {
					ch <- fmt.Errorf("%v", e)
				}
			}()
			if handler, ok := m.handlers[name]; ok {
				if server.Handler != nil {
					ch <- fmt.Errorf("server.%s: handler already exists", name)
					return
				}
				server.Handler = handler
				delete(m.handlers, name)
			}
			if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				ch <- err
			}
		}()
		select {
		case err = <-ch:
			close(ch)
		case <-time.After(100 * time.Millisecond):
			log.Printf("http: server '%s' started, serving on %s", name, server.Addr)
		case <-ctx.Done():
			err = ctx.Err()
		}
		if err != nil {
			return fmt.Errorf("failed to start http.server '%s': %v", name, err)
		}
	}
	if len(m.handlers) > 0 {
		names := []string{}
		for k, _ := range m.handlers {
			names = append(names, k)
		}
		return fmt.Errorf("servers not found: %v", names)
	}
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	for name, server := range m.servers {
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			err = server.Shutdown(ctx)
		}
		if err != nil {
			return fmt.Errorf("server '%s' forced to shutdown: %v", name, err)
		}
		log.Printf("http: server '%s' shutdown: ok", name)
	}
	return nil
}

func (m *Module) Stats() map[string]interface{} {
	return nil
}
