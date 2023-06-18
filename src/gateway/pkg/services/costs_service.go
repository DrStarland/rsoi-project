package services

import (
	"context"
	"fmt"
	obj "gateway/pkg/models/cost"
	"gateway/pkg/myjson"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// NewService creates a new album service.
func NewCostService(client *http.Client, logger *zap.SugaredLogger) costService {
	return costService{client, logger}
}

type costService struct {
	Client *http.Client
	logger *zap.SugaredLogger
}

func (s costService) Create(ctx context.Context, req *obj.CostCreationRequest) (obj.Cost, error) {
	return obj.Cost{}, errors.New("not implemented")
}

// Update updates the album with the specified ID.
func (s costService) Update(ctx context.Context, id int, req *obj.CostCreationRequest) (obj.Cost, error) {
	return obj.Cost{}, errors.New("not implemented")
}

// Delete deletes the album with the specified ID.
func (s costService) Delete(ctx context.Context, id int) (obj.Cost, error) {
	return obj.Cost{}, errors.New("not implemented")
}

// Count returns the number of albums.
func (s costService) Count(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
}

type ByUpdatedAtCost []obj.Cost

func (slice ByUpdatedAtCost) Len() int           { return len(slice) }
func (slice ByUpdatedAtCost) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice ByUpdatedAtCost) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s costService) Query(ctx context.Context, offset, limit int) ([]obj.Cost, error) {
	requestURL := fmt.Sprintf("%s/api/v1/costs", utils.Config.CostsEndpoint)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return []obj.Cost{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return []obj.Cost{}, errors.Wrap(err, "failed request to costs service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	items := []obj.Cost{}
	if err = myjson.From(body, &items); err != nil {
		return []obj.Cost{}, errors.Wrap(err, "failed to decode response")
	}

	sort.Stable(ByUpdatedAtCost(items))
	return items, nil
}

func (s costService) Get(ctx context.Context, id int) (obj.Cost, error) {
	requestURL := fmt.Sprintf("%s/api/v1/costs/%d", utils.Config.CostsEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return obj.Cost{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed request to costs service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	cost := obj.Cost{}
	if err = myjson.From(body, &cost); err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed to decode response")
	}

	return cost, nil
}

func (s costService) CreateFromRequest(ctx context.Context, r *http.Request) (obj.Cost, error) {
	requestURL := fmt.Sprintf("%s/api/v1/costs", utils.Config.CostsEndpoint)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""
	// if err != nil {
	// 	s.logger.Errorln("failed to create an http request")
	// 	return obj.Cost{}, errors.Wrap(err, "failed to create an http request")
	// }

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed request to costs service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	cost := obj.Cost{}
	if err = myjson.From(body, &cost); err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed to decode response")
	}

	return cost, nil
}

func (s costService) UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Cost, error) {
	requestURL := fmt.Sprintf("%s/api/v1/costs/%d", utils.Config.CostsEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed request to costs service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	cost := obj.Cost{}
	if err = myjson.From(body, &cost); err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed to decode response")
	}

	return cost, nil
}

func (s costService) DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Cost, error) {
	requestURL := fmt.Sprintf("%s/api/v1/costs/%d", utils.Config.CostsEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed request to costs service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	cost := obj.Cost{}
	if err = myjson.From(body, &cost); err != nil {
		return obj.Cost{}, errors.Wrap(err, "failed to decode response")
	}

	return cost, nil
}
