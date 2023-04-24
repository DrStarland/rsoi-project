package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"

	"gateway/pkg/handlers"
	"gateway/pkg/session"

	mid "gateway/pkg/middleware"

	"go.uber.org/zap"
)

func HealthOK(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	router := httprouter.New()
	router.PanicHandler = mid.PanicHandler

	sm := session.NewSessionsManager()

	gs := &handlers.GatewayHandler{
		TicketServiceAddress: "http://testum-tickets:8070",
		FlightServiceAddress: "http://testum-flights:8060",
		BonusServiceAddress:  "http://testum-bonus:8050",
		Logger:               logger,
	}

	router.GET("/api/v1/flights", mid.AccessLog(mid.Auth(gs.GetAllFlights, sm, logger), logger))
	router.GET("/api/v1/me", mid.AccessLog(mid.Auth(gs.GetUserInfo, sm, logger), logger))
	router.GET("/api/v1/tickets", mid.AccessLog(mid.Auth(gs.GetUserTickets, sm, logger), logger))
	router.GET("/api/v1/tickets/:ticketUID", mid.AccessLog(mid.Auth(gs.GetUserTicket, sm, logger), logger))
	router.POST("/api/v1/tickets", mid.AccessLog(mid.Auth(gs.BuyTicket, sm, logger), logger))
	router.DELETE("/api/v1/tickets/:ticketUID", mid.AccessLog(mid.Auth(gs.CancelTicket, sm, logger), logger))
	router.GET("/api/v1/privilege", mid.AccessLog(mid.Auth(gs.GetPrivilege, sm, logger), logger))

	router.GET("/manage/health", HealthOK)

	ServerAddress := os.Getenv("PORT")
	if ServerAddress == "" || ServerAddress == ":80" {
		ServerAddress = ":8080"
	} else if !strings.HasPrefix(ServerAddress, ":") {
		ServerAddress = ":" + ServerAddress
	}

	logger.Infow("starting server",
		"type", "START",
		"addr", ServerAddress,
	)

	err := http.ListenAndServe(ServerAddress, router)
	if err != nil {
		log.Panicln(err.Error())
	}
}
