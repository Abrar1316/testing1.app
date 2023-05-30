package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Version = "default"

// Annoyingly, to get hold of the 'AtomicLevel' we need to
// change the logging level, we need a `Config` and, obvs.,
// the global logger `zap.L()` has no access to the internal
// `Config` it holds. Which means we need a `Config` of our
// own even though we don't really want to hold one but we
// don't have to make it exportable since we provide a helper
var zapLogConfig zap.Config

func InitZapLogger() {
	zapLogConfig = zap.NewProductionConfig()

	// Add the version of this build into the logger output.
	// Handy for when you're looking at panics and want to
	// ignore those which happened before the latest release, etc.
	options := []zap.Option{
		zap.Fields(zap.String("version", Version)),
	}

	// If this fails, do we want to warn and bail? There doesn't
	// seem to be any way this can fail with the provided default
	// configuration though...
	logr, _ := zapLogConfig.Build(options...)

	zap.ReplaceGlobals(logr)

	// Mainly for testing but a handy convenience to initialise
	// the logging level from an env.var (`$ZAPLEVEL`)
	switch os.Getenv("ZAPLEVEL") {
	case "DEBUG":
		SetZapLevel(zap.DebugLevel)
	case "INFO":
		SetZapLevel(zap.InfoLevel)
	case "ERROR":
		SetZapLevel(zap.ErrorLevel)
	}

}

func SetZapLevel(level zapcore.Level) {
	zapLogConfig.Level.SetLevel(level)
}
