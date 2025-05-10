package handler

import (
	"main/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RequestHandlerImpl struct {
	gin *gin.Engine
}

func (rh *RequestHandlerImpl) Group(path string) gin.IRoutes {
	return rh.gin.Group(path).Use(gin.Logger(), gin.ErrorLogger(), gin.Recovery(), cors.Default())
}

func (rh *RequestHandlerImpl) Run(addr string) error {
	return rh.gin.Run(addr)
}

func NewRequestHandler(logger pkg.Logger) RequestHandler {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()
	return &RequestHandlerImpl{
		gin: engine,
	}
}
