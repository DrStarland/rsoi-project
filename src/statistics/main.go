package main

import (
	"fmt"
	"net/http"
	"os"
	"statistics/controllers"
	"statistics/repository"
	"statistics/utils"

	"time"

	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	time.Sleep(5 * time.Second)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	// rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	kafka := utils.InitKafka(logger)
	defer kafka.Close()
	go kafka.ConsumeLoop()

	db := utils.InitDBConnection(utils.Config.DB, logger)
	defer db.Close()
	go utils.StatWriteLoop(repository.NewPGStatisticsRep(db))

	r := controllers.InitRouter(db)

	utils.Logger.Print("Server started")

	logger.Infow("starting server",
		"type", "START",
		"addr", fmt.Sprintf("http://localhost:%d", utils.Config.Port),
	)

	err := controllers.RunRouter(r, utils.Config.Port)

	if err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}
