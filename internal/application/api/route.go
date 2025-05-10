package api

import (
	"main/pkg/handler"

	"go.uber.org/fx"
)

func SetupRoutes(handler handler.RequestHandler, nodeController *NodeController, namespaceController *NamespaceController) {
	nodeGroup := handler.Group("/api/nodes")
	nodeGroup.GET("", nodeController.GetNodes)

	namespaceGroup := handler.Group("/api/namespaces")
	namespaceGroup.GET("", namespaceController.GetNamespaces)
}

var Module = fx.Module("api", fx.Provide(NewNamespaceController), fx.Provide(NewNodeController), fx.Invoke(SetupRoutes))
