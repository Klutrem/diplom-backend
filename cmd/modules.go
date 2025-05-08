package cmd

import (
	"main/internal/application/api"
	"main/internal/config"
	"main/internal/infrastructure/kubernetes"
	"main/pkg"
	"main/pkg/handler"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,
	handler.Module,
	kubernetes.Module,
	api.Module,
)
