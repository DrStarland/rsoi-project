package main

import (
	"net/http"
	"os"
	"time"

	"gateway/pkg/initialization"

	"go.uber.org/zap"
)

func main() {
	time.Sleep(5 * time.Second)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	noteApp := initialization.NewApp(logger)
	defer noteApp.Stop()

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
