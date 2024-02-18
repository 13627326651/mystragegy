package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var defaultLogger zap.Logger

func init() {
	logger, err := NewLogger(DefaultConfig())
	if err != nil {
		panic(fmt.Errorf("init logger: %v", err))
	}

	defaultLogger = *logger
}

func InitDefaultLogger(config Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}

	defaultLogger = *logger

	return nil
}

func DefaultLogger() *zap.Logger {
	return &defaultLogger
}
