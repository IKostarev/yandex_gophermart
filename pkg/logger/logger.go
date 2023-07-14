package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()

	cfg.Level = lvl

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error config build is: %s", err)
	}

	return zapLogger, nil
}
