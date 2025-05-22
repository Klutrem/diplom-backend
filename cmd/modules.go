package cmd

import (
	"main/internal/application/api"
	"main/internal/config"
	"main/internal/domain/alerts"
	"main/internal/domain/events"
	"main/internal/infrastructure/database"
	"main/internal/infrastructure/kubernetes"
	"main/internal/infrastructure/prometheus"
	"main/pkg"
	"main/pkg/handler"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,
	handler.Module,
	kubernetes.Module,
	prometheus.Module,
	api.Module,
	events.Module,
	database.Module,
	alerts.Module,
)
