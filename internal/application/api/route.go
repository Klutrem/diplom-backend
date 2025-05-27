package api

import (
	"main/pkg/handler"

	"go.uber.org/fx"
)

func SetupRoutes(handler handler.RequestHandler, nodeController *NodeController, namespaceController *NamespaceController, podController *PodController, eventController *EventController, telegramAlertController *TelegramAlertController) {
	nodeGroup := handler.Group("/api/nodes")
	{
		nodeGroup.GET("", nodeController.GetNodes)
		nodeGroup.GET("/metrics/:node", nodeController.GetNodeMetrics)
	}

	namespaceGroup := handler.Group("/api/namespaces")
	{
		namespaceGroup.GET("", namespaceController.GetNamespaces)
	}

	podGroup := handler.Group("/api/pods")
	{
		podGroup.GET("", podController.GetPods)
		podGroup.GET("/metrics/:namespace/:pod", podController.GetPodMetrics)
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

	alertsGroup := handler.Group("/api/alerts")
	{
		alertsGroup.POST("", telegramAlertController.CreateAlert)
		alertsGroup.GET("/namespace/:namespace", telegramAlertController.GetAlertsByNamespace)
		alertsGroup.DELETE("/:id", telegramAlertController.DeleteAlert)
		alertsGroup.PUT("/:id", telegramAlertController.UpdateAlert)
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
