package src

import (
	"go.uber.org/zap"
)

func TestZapLogger() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	sugar.Infof("Starting server on port %d", 8080) // want `log message should start with a lowercase letter`
	sugar.Infof("starting server on port %d", 8080) // OK
}
