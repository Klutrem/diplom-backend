package pkg

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(GetLogger),
	fx.Provide(NewPostgresDatabase),
)
