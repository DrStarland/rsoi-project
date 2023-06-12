package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/pkg/models/statistic"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"time"
)

type FetchResponse struct {
	Reqests []statistic.RequestStat `json:"requests"`
}

type StatisticsM struct {
	client *http.Client
}

func NewStatisticsM(client *http.Client) *StatisticsM {
	return &StatisticsM{client: client}
}

func (model *StatisticsM) Fetch(beginTime time.Time, endTime time.Time, authHeader string) *statistic.FetchResponse {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/requests", utils.Config.StatEndpoint), nil)
	q := req.URL.Query()
	q.Add("begin_time", beginTime.Format(time.RFC3339))
	q.Add("end_time", endTime.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", authHeader)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("client: error making http request\n")
	}

	data := &statistic.FetchResponse{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, data)
	return data
}

type statisticsCtrl struct {
	statistics *StatisticsM
}

// func InitStatistics(r *mux.Router, statistics *StatisticsM) {
// 	ctrl := &statisticsCtrl{statistics}
// 	r.HandleFunc("/requests", ctrl.fetch).Methods("GET")
// }

// func (ctrl *statisticsCtrl) fetch(w http.ResponseWriter, r *http.Request) {
// 	if role := r.Header.Get("X-User-Role"); role != "admin" {
// 		responses.ForbiddenMsg(w, fmt.Sprintf("not allowed for %s role", role))
// 		return
// 	}

// 	queryParams := r.URL.Query()
// 	log.Println(queryParams.Get("begin_time"))
// 	log.Println(queryParams.Get("end_time"))
// 	begin_time, err := time.Parse(time.RFC3339, queryParams.Get("begin_time"))
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("bad begin_time format: %s", err.Error()), http.StatusBadRequest)
// 		return
// 	}
// 	end_time, err := time.Parse(time.RFC3339, queryParams.Get("end_time"))
// 	if err != nil {
// 		http.Error(w, "bad end_time format", http.StatusBadRequest)
// 		return
// 	}

// 	data := ctrl.statistics.Fetch(begin_time, end_time, r.Header.Get("Authorization"))
// 	responses.JsonSuccess(w, data)
// }
