package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	app "github.com/venusforest2013/config/application"
	"github.com/venusforest2013/config/modules/http"
	"github.com/venusforest2013/config/modules/utils"
)

var (
	defaultMode       = "debug"
	defaultRouterName = "default"
	defaultServerRef  = "default"
)

func init() {
	app.Register(newModule())
}

type Module struct {
	config  *app.Config
	routers map[string]*Router
}

func newModule() *Module {
	return &Module{
		routers: map[string]*Router{},
	}
}

func (m *Module) Router(name string) *Router {
	r, ok := m.routers[name]
	if !ok {
		r = NewRouter()
		m.routers[name] = r
	}
	return r
}

func (m *Module) Register(name, path string, c Controller) {
	m.Router(name).RegisterController(path, c)
}

func (m *Module) Name() string {
	return "gin"
}

func (m *Module) Configure(ctx context.Context, cfg *app.Config) error {
	v, err := utils.Module(ctx, "http")
	if err != nil {
		return err
	}
	httpModule := v.(*http.Module)
	gin.SetMode(cfg.Str("mode", defaultMode))
	routersConf := cfg.ValueArray("routers")
	if routersConf == nil {
		return errors.New("gin.routers not found")
	}
	for _, conf := range routersConf {
		name := conf.Str("name", defaultRouterName)
		router := m.Router(name)
		if err := router.configure(conf); err != nil {
			return err
		}
		ref := conf.Str("server_ref", defaultServerRef)
		httpModule.RegisterHandler(ref, router)
		m.routers[name] = router
	}
	m.config = cfg
	return nil
}

func (m *Module) Start(ctx context.Context) error {
	for _, r := range m.routers {
		for _, c := range r.Controllers() {
			if err := c.Init(ctx); err != nil {
				return fmt.Errorf("failed to init gin.controller '%s': %v",
					c.Name(), err)
			}
		}
	}
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	for _, r := range m.routers {
		for _, c := range r.Controllers() {
			if err := c.Destroy(ctx); err != nil {
				return fmt.Errorf("failed to destroy gin.controller '%s': %v",
					c.Name(), err)
			}
		}
	}
	return nil
}

func (m *Module) Stats() map[string]interface{} {
	return nil
}
