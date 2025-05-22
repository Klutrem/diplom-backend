package database

import "go.uber.org/fx"

var Module = fx.Module("database",
	fx.Provide(NewEventPGRepository),
	fx.Provide(NewWatchedNamespacePGRepository),
	fx.Provide(NewTelegramAlertPGRepository),
)
