package main

import (
	"net/http"
	"os"
	"tasks/pkg/initialization"

	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	app := initialization.NewApp(logger)
	defer func() {
		app.Stop()
	}()

	logger.Infow("starting server",
		"type", "START",
		"addr", app.ServerAddress,
	)

	err := http.ListenAndServe(app.ServerAddress, app.Router)
	if err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}
