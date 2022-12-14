package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log      *zap.Logger
	logLevel zap.AtomicLevel
)

func InitTestLogger() {
	os.Setenv("CONSOLE_LOGGING", "enabled")
	InitLogger("json", "debug")
}

func InitLogger(format, level string) {
	// Log to the console by default.
	logLevel = zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	var core zapcore.Core
	switch format {
	case "json":
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			logLevel)
	case "text":
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			logLevel)
	default:
		panic("Invalid log format chosen: " + format)
	}
	Log = zap.New(core, zap.AddCaller())

	SetLogLevel(level)
}

func ShutdownLogger() {
	_ = Log.Sync()
}

func SetLogLevel(level string) {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		// Fallback to logging at the info level.
		fmt.Printf("Falling back to the info log level. You specified: %s.\n",
			level)
		logLevel.SetLevel(zapcore.InfoLevel)
	} else {
		logLevel.SetLevel(parsedLevel)
	}
}
