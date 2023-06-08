package services

import (
	"context"
	"fmt"
	obj "gateway/pkg/models/task"
	"gateway/pkg/myjson"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Service encapsulates usecase logic for albums.
type TaskService interface {
	Get(ctx context.Context, id int) (obj.Task, error)
	Query(ctx context.Context, offset, limit int) ([]obj.Task, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *obj.TaskCreationRequest) (obj.Task, error)
	Update(ctx context.Context, id int, input *obj.TaskCreationRequest) (obj.Task, error)
	Delete(ctx context.Context, id int) (obj.Task, error)

	AddComment(ctx context.Context, id int, com obj.CommentCreationRequest) (obj.Task, error)

	CreateFromRequest(ctx context.Context, r *http.Request) (obj.Task, error)
	UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error)
	DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error)
	AddCommentFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error)
}

type taskService struct {
	Client *http.Client
	logger *zap.SugaredLogger
}

// NewService creates a new album service.
func NewTaskService(client *http.Client, logger *zap.SugaredLogger) taskService {
	return taskService{client, logger}
}

type byUpdatedAt []obj.Task

func (slice byUpdatedAt) Len() int           { return len(slice) }
func (slice byUpdatedAt) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice byUpdatedAt) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s taskService) Query(ctx context.Context, offset, limit int) ([]obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks", utils.Config.TasksEndpoint)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return []obj.Task{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return []obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	items := []obj.Task{}
	if err = myjson.From(body, &items); err != nil {
		return []obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	sort.Stable(byUpdatedAt(items))
	return items, nil
}

// Get returns the album with the specified the album ID.
func (s taskService) Get(ctx context.Context, id int) (obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks/%d", utils.Config.TasksEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return obj.Task{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	task := obj.Task{}
	if err = myjson.From(body, &task); err != nil {
		return obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	return task, nil
}

func (s taskService) Create(ctx context.Context, req *obj.TaskCreationRequest) (obj.Task, error) {
	// if err := req.Validate(); err != nil {
	// 	return obj.Task{}, err
	// }

	// now := timestamp.Now()

	insertion := obj.Task{}
	// 	Author: user.User{
	// 		ID:       ctx.Value("X-UID").(string),
	// 		Username: ctx.Value("X-User-Name").(string),
	// 	},
	// 	Status: obj.StatusQueued,
	// 	// VisibilityScope:     req.VisibilityScope,
	// 	CreatedAt:           timestamp.Timestamp(now),
	// 	UpdatedAt:           timestamp.Timestamp(now),
	// 	TaskCreationRequest: *req,
	// 	Comments:            make([]obj.Comment, 0),
	// }

	// err := s.repo.Create(ctx, &insertion)
	// if err != nil {
	// 	return obj.Task{}, errors.Wrap(err, "service work failed")
	// }
	return insertion, nil
}

func (s taskService) CreateFromRequest(ctx context.Context, r *http.Request) (obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks", utils.Config.TasksEndpoint)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""
	// if err != nil {
	// 	s.logger.Errorln("failed to create an http request")
	// 	return obj.Task{}, errors.Wrap(err, "failed to create an http request")
	// }

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	task := obj.Task{}
	if err = myjson.From(body, &task); err != nil {
		return obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	return task, nil
}

// Update updates the album with the specified ID.
func (s taskService) Update(ctx context.Context, id int, req *obj.TaskCreationRequest) (obj.Task, error) {
	task := obj.Task{}
	return task, nil
}

func (s taskService) UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks/%d", utils.Config.TasksEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	task := obj.Task{}
	if err = myjson.From(body, &task); err != nil {
		return obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	return task, nil
}

// Delete deletes the album with the specified ID.
func (s taskService) Delete(ctx context.Context, id int) (obj.Task, error) {
	return obj.Task{}, nil
}

func (s taskService) DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks/%d", utils.Config.TasksEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	task := obj.Task{}
	if err = myjson.From(body, &task); err != nil {
		return obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	return task, nil
}

// Count returns the number of albums.
func (s taskService) Count(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
}

// Update updates the album with the specified ID.
func (s taskService) AddComment(ctx context.Context, id int, req obj.CommentCreationRequest) (obj.Task, error) {
	return obj.Task{}, errors.New("not implemented")
}

func (s taskService) AddCommentFromRequest(ctx context.Context, r *http.Request, id int) (obj.Task, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tasks/%d", utils.Config.TasksEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "failed request to tasks service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	task := obj.Task{}
	if err = myjson.From(body, &task); err != nil {
		return obj.Task{}, errors.Wrap(err, "failed to decode response")
	}

	return task, nil
}
