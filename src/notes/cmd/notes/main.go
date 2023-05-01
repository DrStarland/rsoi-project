package main

import (
	"net/http"
	"notes/pkg/initialization"
	"os"

	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	noteApp := initialization.NewApp(logger)
	defer func() {
		noteApp.Stop()
	}()

	logger.Infow("starting server",
		"type", "START",
		"addr", noteApp.ServerAddress,
	)

	err := http.ListenAndServe(noteApp.ServerAddress, noteApp.Router)
	if err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}
