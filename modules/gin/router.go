package gin

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	app "github.com/venusforest2013/config/application"
	"net/http/pprof"
	"os"
)

var (
	defaultHealthUriCheck  = false
	defaultHealthUri       = "/checkpreload.htm"
	defaultHealthFileCheck = false
	defaultHealthFileUri   = "/status.taobao"
	defaultHealthFilePath  = "/home/admin/cai/htdocs/status.taobao"
)

type Router struct {
	*gin.Engine
	controllers map[string]Controller
}

func NewRouter() *Router {
	engine := gin.New()
	router := &Router{
		Engine:      engine,
		controllers: map[string]Controller{},
	}
	return router
}

func (r *Router) RegisterController(path string, c Controller) {
	r.controllers[path] = c
	r.GET(path, r.Dispatch(c, c.Get))
	r.POST(path, r.Dispatch(c, c.Post))
}

func (r *Router) Dispatch(c Controller, handler HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		c.Before(ctx)
		handler(ctx)
		c.After(ctx)
	}
}

func (r *Router) Controllers() map[string]Controller {
	return r.controllers
}

func (r *Router) configure(cfg *app.Config) error {

	if cfg.Bool("pprof_enabled", false) {
		r.pprof(r.Engine.Group(cfg.Str("pprof_path", "/debug/pprof")))
	}

	if cfg.Bool("swagger_enabled", false) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	if cfg.Bool("health_uri_check", defaultHealthUriCheck) {
		r.GET(cfg.Str("health_uri", defaultHealthUri), func(c *gin.Context) {
			c.String(200, "success")
		})
	}

	if cfg.Bool("health_file_check", defaultHealthFileCheck) {
		uri := cfg.Str("health_file_uri", defaultHealthFileUri)
		r.GET(uri, func(c *gin.Context) {
			path := cfg.Str("health_file_path", defaultHealthFilePath)
			if _, err := os.Stat(path); err == nil {
				c.String(200, "success")
			} else {
				c.Status(404)
			}
		})
	}

	return nil
}

func (r *Router) pprof(rg *gin.RouterGroup) {
	rg.GET("/", gin.WrapF(pprof.Index))
	rg.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	rg.GET("/profile", gin.WrapF(pprof.Profile))
	rg.POST("/symbol", gin.WrapF(pprof.Symbol))
	rg.GET("/symbol", gin.WrapF(pprof.Symbol))
	rg.GET("/trace", gin.WrapF(pprof.Trace))
	rg.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	rg.GET("/block", gin.WrapH(pprof.Handler("block")))
	rg.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	rg.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	rg.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	rg.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
}
