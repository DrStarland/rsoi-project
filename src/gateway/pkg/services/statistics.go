package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/pkg/models/statistic"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type StatisticService interface {
	Query(ctx context.Context, beginTime time.Time, endTime time.Time) (*statistic.FetchResponse, error)
}

// NewService creates a new album service.
func NewStatisticService(client *http.Client, logger *zap.SugaredLogger) statisticService {
	return statisticService{client, logger}
}

type statisticService struct {
	Client *http.Client
	logger *zap.SugaredLogger
}

func (model statisticService) Query(ctx context.Context, beginTime time.Time, endTime time.Time) (*statistic.FetchResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/requests", utils.Config.StatEndpoint), nil)
	q := req.URL.Query()
	q.Add("begin_time", beginTime.Format(time.RFC3339))
	q.Add("end_time", endTime.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("client: error making http request")
	}

	data := &statistic.FetchResponse{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, data)
	return data, nil
}
