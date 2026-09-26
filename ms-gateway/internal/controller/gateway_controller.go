package controller

import (
	"context"
	"net/http/httputil"
	"sync"

	coreConroller "github.com/CakeForKit/rsoi-lab2/lb-core/controller"
	"github.com/CakeForKit/rsoi-lab2/ms-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var gatewayCntrlr coreConroller.HttpController
var gatewayCntrlrMutex sync.Mutex

type gatewayController struct {
	proxy *httputil.ReverseProxy
}

func GetGatewayController() (coreConroller.HttpController, error) {
	gatewayCntrlrMutex.Lock()
	gatewayCntrlrMutex.Unlock()

	if gatewayCntrlr != nil {
		return gatewayCntrlr, nil
	}
	gatewayCntrlr = &gatewayController{proxy: proxy.NewGatewayProxy()}
	return gatewayCntrlr, nil
}

func (controller gatewayController) RegisterHttpController(router *gin.Engine) {
	router.NoRoute(controller.transfer)
}

func (controller gatewayController) transfer(ctx *gin.Context) {
	log.WithContext(ctx).Debug("GatewayController: transfer: Start")
	controller.proxy.ServeHTTP(ctx.Writer, ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), gin.ContextKey, ctx)))
	log.WithContext(ctx).Debug("GatewayController: transfer: Success")
}
