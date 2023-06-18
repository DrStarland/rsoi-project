package services

import (
	"context"
	"fmt"
	obj "gateway/pkg/models/income"
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
func NewIncomeService(client *http.Client, logger *zap.SugaredLogger) incomeService {
	return incomeService{client, logger}
}

type incomeService struct {
	Client *http.Client
	logger *zap.SugaredLogger
}

func (s incomeService) Create(ctx context.Context, req *obj.IncomeCreationRequest) (obj.Income, error) {
	return obj.Income{}, errors.New("not implemented")
}

// Update updates the album with the specified ID.
func (s incomeService) Update(ctx context.Context, id int, req *obj.IncomeCreationRequest) (obj.Income, error) {
	return obj.Income{}, errors.New("not implemented")
}

// Delete deletes the album with the specified ID.
func (s incomeService) Delete(ctx context.Context, id int) (obj.Income, error) {
	return obj.Income{}, errors.New("not implemented")
}

// Count returns the number of albums.
func (s incomeService) Count(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
}

type ByUpdatedAtIncome []obj.Income

func (slice ByUpdatedAtIncome) Len() int           { return len(slice) }
func (slice ByUpdatedAtIncome) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice ByUpdatedAtIncome) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s incomeService) Query(ctx context.Context, offset, limit int) ([]obj.Income, error) {
	requestURL := fmt.Sprintf("%s/api/v1/incomes", utils.Config.CostsEndpoint)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return []obj.Income{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return []obj.Income{}, errors.Wrap(err, "failed request to incomes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	items := []obj.Income{}
	if err = myjson.From(body, &items); err != nil {
		return []obj.Income{}, errors.Wrap(err, "failed to decode response")
	}

	sort.Stable(ByUpdatedAtIncome(items))
	return items, nil
}

func (s incomeService) Get(ctx context.Context, id int) (obj.Income, error) {
	requestURL := fmt.Sprintf("%s/api/v1/incomes/%d", utils.Config.CostsEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return obj.Income{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "failed request to incomes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	income := obj.Income{}
	if err = myjson.From(body, &income); err != nil {
		return obj.Income{}, errors.Wrap(err, "failed to decode response")
	}

	return income, nil
}

func (s incomeService) CreateFromRequest(ctx context.Context, r *http.Request) (obj.Income, error) {
	requestURL := fmt.Sprintf("%s/api/v1/incomes", utils.Config.CostsEndpoint)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""
	// if err != nil {
	// 	s.logger.Errorln("failed to create an http request")
	// 	return obj.Income{}, errors.Wrap(err, "failed to create an http request")
	// }

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "failed request to incomes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	income := obj.Income{}
	if err = myjson.From(body, &income); err != nil {
		return obj.Income{}, errors.Wrap(err, "failed to decode response")
	}

	return income, nil
}

func (s incomeService) UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Income, error) {
	requestURL := fmt.Sprintf("%s/api/v1/incomes/%d", utils.Config.CostsEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "failed request to incomes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	income := obj.Income{}
	if err = myjson.From(body, &income); err != nil {
		return obj.Income{}, errors.Wrap(err, "failed to decode response")
	}

	return income, nil
}

func (s incomeService) DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Income, error) {
	requestURL := fmt.Sprintf("%s/api/v1/incomes/%d", utils.Config.CostsEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "failed request to incomes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	income := obj.Income{}
	if err = myjson.From(body, &income); err != nil {
		return obj.Income{}, errors.Wrap(err, "failed to decode response")
	}

	return income, nil
}
