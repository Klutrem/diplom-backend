package cmd

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/pkg"
	"main/pkg/handler"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func Run() any {
	return func(
		env config.Env,
		logger pkg.Logger,
		handler handler.RequestHandler,
	) {
		logger.Fatal(handler.Run(fmt.Sprint(env.ServerAddress, ":", env.Port)))
	}
}

func StartApp() error {
	logger := pkg.GetLogger(config.NewEnv())
	opts := fx.Options(
		fx.WithLogger(func() fxevent.Logger {
			return logger.GetFxLogger()
		}),
		fx.Invoke(Run()),
	)
	ctx := context.Background()
	app := fx.New(CommonModules, opts)
	err := app.Start(ctx)
	return err
}
