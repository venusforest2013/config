package controller

import (
	"github.com/venusforest2013/config/modules/gin"
	"net/http"
	"runtime"
)

type PingController struct {
	gin.FrontController
}

func (p PingController) Post(ginCtx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			ginCtx.JSON(http.StatusInternalServerError, string(buf[:n]))
		}
	}()

	ginCtx.JSON(http.StatusOK, "success")

}
