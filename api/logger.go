package main

import "go.uber.org/zap"
import "go.uber.org/zap/zapcore"

func NewLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()               // console, human-readable [web:21]
    cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // colored levels [web:16]
    logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}