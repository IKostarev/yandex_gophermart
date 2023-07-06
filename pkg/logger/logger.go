package logger

import "go.uber.org/zap"

func New(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()

	cfg.Level = lvl

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger, nil
}
