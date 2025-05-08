package api

import (
	"main/pkg/handler"

	"go.uber.org/fx"
)

func SetupRoutes(handler handler.RequestHandler, nodeController *NodeController) {
	nodeGroup := handler.Group("/nodes")
	nodeGroup.GET("", nodeController.GetNodes)
}

var Module = fx.Module("api", fx.Provide(NewNodeController), fx.Invoke(SetupRoutes))
