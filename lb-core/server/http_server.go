package server

import (
	"fmt"

	"github.com/CakeForKit/rsoi-lab1/internal/common/config"
	"github.com/CakeForKit/rsoi-lab1/internal/common/controller"
	"github.com/CakeForKit/rsoi-lab1/internal/common/utils"
	"github.com/gin-gonic/gin"
)

type Runnable interface {
	Run() error
}

type httpServer struct {
	router      *gin.Engine
	startupFunc func(*gin.Engine)
}

func NewHttpServer(startupFunc func(*gin.Engine)) Runnable {
	return &httpServer{
		startupFunc: startupFunc,
	}
}

func (server *httpServer) Run() error {
	server.router = gin.New()
	server.router.Use(gin.Recovery())
	controller.NewHealthCheckController().RegisterHttpController(server.router)

	server.startupFunc(server.router)
	utils.CheckedError(server.router.Run(fmt.Sprintf(":%d", config.CoreConfig.Port)))
	return nil
}
