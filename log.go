package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLog() {
	logger, err := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder, // INFO

			TimeKey:    "time",
			EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),

			CallerKey:        "caller",
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
			FunctionKey:      "",
		},
	}.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}