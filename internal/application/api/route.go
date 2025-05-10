package api

import (
	"main/pkg/handler"

	"go.uber.org/fx"
)

func SetupRoutes(handler handler.RequestHandler, nodeController *NodeController, namespaceController *NamespaceController, podController *PodController) {
	nodeGroup := handler.Group("/api/nodes")
	nodeGroup.GET("", nodeController.GetNodes)

	namespaceGroup := handler.Group("/api/namespaces")
	namespaceGroup.GET("", namespaceController.GetNamespaces)

	podGroup := handler.Group("/api/pods")
	podGroup.GET("", podController.GetPods)
}

var Module = fx.Module("api", fx.Provide(NewNamespaceController), fx.Provide(NewNodeController), fx.Provide(NewPodController), fx.Invoke(SetupRoutes))
