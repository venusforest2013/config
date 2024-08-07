package main

import (
	"context"
	"github.com/venusforest2013/config/application"
	app "github.com/venusforest2013/config/application"
	"github.com/venusforest2013/config/cmd/controller"
	"github.com/venusforest2013/config/modules/gin"
)

func init() {
	app.Register(&AppModule{})
}

type AppModule struct {
}

func (a *AppModule) Name() string {
	return "app"
}

func (a *AppModule) Configure(ctx context.Context, cfg *app.Config) error {

	//初始化gin
	ginModule := application.GetModule("gin").(*gin.Module)
	router := ginModule.Router("default")
	router.RegisterController("ping", &controller.PingController{})

	return nil
}

func (a *AppModule) Start(ctx context.Context) error {

	return nil
}

func (a *AppModule) Shutdown(ctx context.Context) error {

	return nil
}

func (a *AppModule) Stats() map[string]interface{} {
	return nil
}
