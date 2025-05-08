package handler

import "github.com/gin-gonic/gin"

type RequestHandler interface {
	Group(path string) gin.IRoutes
	Run(addr string) error
}
