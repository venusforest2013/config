package gin

import (
	"context"
	"net/http"
)

type FrontController struct {
}

func (c FrontController) Name() string {
	return "FrontController"
}

func (c FrontController) Init(ctx context.Context) error {
	return nil
}

func (c FrontController) Destroy(ctx context.Context) error {
	return nil
}

func (c FrontController) Before(ctx *Context) {}

func (c FrontController) After(ctx *Context) {}

func (c FrontController) Get(ctx *Context) {
	MethodNotAllowed(ctx)
}

func (c FrontController) Post(ctx *Context) {
	MethodNotAllowed(ctx)
}

func MethodNotAllowed(ctx *Context) {
	ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
