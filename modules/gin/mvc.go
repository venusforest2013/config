package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type Context = gin.Context
type HandlerFunc = gin.HandlerFunc
type Model = gin.H
type View = render.Render

type Controller interface {
	Name() string
	Init(ctx context.Context) error
	Destroy(ctx context.Context) error
	Before(ctx *Context)
	After(ctx *Context)
	Get(ctx *Context)
	Post(ctx *Context)
	//Put(ctx *Context)
	//Patch(ctx *Context)
	//Head(ctx *Context)
	//Options(ctx *Context)
	//Delete(ctx *Context)
	//Connect(ctx *Context)
	//Trace(ctx *Context)
}

type ViewResolver interface {
	Resolve(view string) (View, error)
}
