package api

import (
	"main/pkg/handler"

	"go.uber.org/fx"
)

func SetupRoutes(handler handler.RequestHandler, nodeController *NodeController, namespaceController *NamespaceController, podController *PodController, eventController *EventController, telegramAlertController *TelegramAlertController) {
	nodeGroup := handler.Group("/api/nodes")
	{
		nodeGroup.GET("", nodeController.GetNodes)

	}

	namespaceGroup := handler.Group("/api/namespaces")
	{
		namespaceGroup.GET("", namespaceController.GetNamespaces)
	}

	podGroup := handler.Group("/api/pods")
	{
		podGroup.GET("", podController.GetPods)
	}

	eventsGroup := handler.Group("/api/events")
	{
		eventsGroup.GET("", eventController.ListEvents)
	}

	watchedNamespacesGroup := handler.Group("/api/watched_namespaces")
	{
		watchedNamespacesGroup.GET("", eventController.GetWatchedNamespaces)
		watchedNamespacesGroup.POST("", eventController.AddWatchedNamespace)
		watchedNamespacesGroup.DELETE("/:namespace", eventController.RemoveWatchedNamespace)
	}
}

var Module = fx.Module("api",
	fx.Provide(NewNamespaceController),
	fx.Provide(NewNodeController),
	fx.Provide(NewPodController),
	fx.Invoke(SetupRoutes),
	fx.Provide(NewEventController),
	fx.Provide(NewTelegramAlertController),
)
