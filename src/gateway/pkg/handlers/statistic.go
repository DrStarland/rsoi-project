package handlers

import (
	"context"
	"fmt"
	"gateway/pkg/models/statistic"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StatisticHandler interface {
	List(w http.ResponseWriter, r *http.Request)
}

type StatisticMainHandler struct {
	Logger  *zap.SugaredLogger
	Service services.StatisticService
}

func NewStatisticsHandler(logger *zap.SugaredLogger) (h *StatisticMainHandler) {
	client := &http.Client{}

	ctrl := services.NewStatisticService(client, logger)
	h = &StatisticMainHandler{
		Logger:  logger,
		Service: ctrl,
	}
	return h
}

type FetchResponse struct {
	Reqests []statistic.RequestStat `json:"requests"`
}

func (h *StatisticMainHandler) List(w http.ResponseWriter, r *http.Request) {
	if role := r.Header.Get("X-User-Role"); role != "admin" {
		myjson.JSONError(w, http.StatusMethodNotAllowed, fmt.Sprintf("not allowed for %s role", role))
		return
	}

	queryParams := r.URL.Query()
	log.Println(queryParams.Get("begin_time"))
	log.Println(queryParams.Get("end_time"))
	begin_time, err := time.Parse(time.RFC3339, queryParams.Get("begin_time"))
	if err != nil {
		http.Error(w, fmt.Sprintf("bad begin_time format: %s", err.Error()), http.StatusBadRequest)
		return
	}
	end_time, err := time.Parse(time.RFC3339, queryParams.Get("end_time"))
	if err != nil {
		http.Error(w, "bad end_time format", http.StatusBadRequest)
		return
	}

	data, err := h.Service.Query(
		context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")),
		begin_time, end_time,
	)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, "in stat handler").Error())
	}

	myjson.JSONResponce(w, http.StatusOK, data)
}

// func (h *StatisticMainHandler) List(w http.ResponseWriter, r *http.Request) {
// 	queryParams := r.URL.Query()
// 	begin_time, err := time.Parse(time.RFC3339, queryParams.Get("begin_time"))
// 	if err != nil {
// 		http.Error(w, "bad begin_time format", http.StatusBadRequest)
// 		return
// 	}
// 	end_time, err := time.Parse(time.RFC3339, queryParams.Get("end_time"))
// 	if err != nil {
// 		http.Error(w, "bad end_time format", http.StatusBadRequest)
// 		return
// 	}

// 	requests, err := h.Service.Query(ctx, begin_time, end_time)
// 	if err != nil {
// 		http.Error(w, "failed to fetch statistics", http.StatusInternalServerError)
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, FetchResponse{Reqests: requests})
// }
