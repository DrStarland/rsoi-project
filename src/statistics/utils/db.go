package utils

import (
	"context"
	"fmt"
	"log"
	"statistics/objects"
	"statistics/repository"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type DBConfiguration struct {
	Type string `json:"type"`
	Name string `json:"name"`

	User     string `json:"user"`
	Password string `json:"password"`

	Host string `json:"host"`
	Port string `json:"port"`
}

func InitDBConnection(cnf DBConfiguration, logger *zap.SugaredLogger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cnf.Host, cnf.User, cnf.Password, cnf.Name, cnf.Port)
	db, e := gorm.Open(cnf.Type, dsn)

	if e != nil {
		logger.Errorln("DB Connection failed: %s", e.Error())
	} else {
		logger.Infoln("DB Connection Established")
	}

	db.SingularTable(true)
	db.AutoMigrate(&objects.RequestStat{})
	return db
}

func StatWriteLoop(rep repository.StatisticsRep) {
	ctx := context.Background()
	for {
		message := GetMessage()
		rep.Create(message)

		if err := ctx.Err(); err != nil {
			log.Panic(err)
		}
	}
}
