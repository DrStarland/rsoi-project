package handlers

import (
	"context"
	"sort"
	obj "tasks/pkg/models/task"
	"tasks/pkg/models/timestamp"
	"tasks/pkg/models/user"

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
}

// // CreateAlbumRequest represents an album creation request.
// type CreateAlbumRequest struct {
// 	Name string `json:"name"`
// }

// // Validate validates the CreateAlbumRequest fields.
// func (m CreateAlbumRequest) Validate() error {
// 	return validation.ValidateStruct(&m,
// 		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
// 	)
// }

// // UpdateAlbumRequest represents an album update request.
// type UpdateAlbumRequest struct {
// 	Name string `json:"name"`
// }

// // Validate validates the CreateAlbumRequest fields.
// func (m UpdateAlbumRequest) Validate() error {
// 	return validation.ValidateStruct(&m,
// 		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
// 	)
// }

type taskService struct {
	repo   obj.Repository
	logger zap.SugaredLogger
}

// NewService creates a new album service.
func NewTaskService(repo obj.Repository, logger zap.SugaredLogger) taskService {
	return taskService{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s taskService) Get(ctx context.Context, id int) (obj.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "service work failed")
	}
	return task, nil
}

// type CommentCreationRequest struct {
// 	Content string
// }

// type Comment struct {
// 	Author    user.User
// 	CreatedAt timestamp.Timestamp
// 	CommentCreationRequest
// }

// type TaskCreationRequest struct {
// 	Title           string
// 	Description     string
// 	NoteURL         string
// 	VisibilityScope scope.Scope
// }

// type Task struct {
// 	ID        int       `json:"id"`
// 	Author    user.User `json:"author"` // author user.User
// 	Status    string
// 	CreatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
// 	UpdatedAt timestamp.Timestamp
// 	Comments  []Comment

// 	TaskCreationRequest
// }

func (s taskService) Create(ctx context.Context, req *obj.TaskCreationRequest) (obj.Task, error) {
	// if err := req.Validate(); err != nil {
	// 	return obj.Task{}, err
	// }

	now := timestamp.Now()

	insertion := obj.Task{
		Author: user.User{
			ID:       ctx.Value("X-UID").(string),
			Username: ctx.Value("X-User-Name").(string),
		},
		Status: obj.StatusQueued,
		// VisibilityScope:     req.VisibilityScope,
		CreatedAt:           timestamp.Timestamp(now),
		UpdatedAt:           timestamp.Timestamp(now),
		TaskCreationRequest: *req,
		Comments:            make([]obj.Comment, 0),
	}

	err := s.repo.Create(ctx, &insertion)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "service work failed")
	}
	return insertion, nil
}

// Update updates the album with the specified ID.
func (s taskService) Update(ctx context.Context, id int, req *obj.TaskCreationRequest) (obj.Task, error) {
	now := timestamp.Timestamp(timestamp.Now())

	task, _ := s.Get(ctx, id)

	task.UpdatedAt = now
	task.TaskCreationRequest = *req

	err := s.repo.Update(ctx, task)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "service work failed")
	}

	return task, nil
}

// Delete deletes the album with the specified ID.
func (s taskService) Delete(ctx context.Context, id int) (obj.Task, error) {
	album, err := s.Get(ctx, id)
	if err != nil {
		return obj.Task{}, errors.New("not exist")
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return obj.Task{}, err
	}
	return album, nil
}

// Count returns the number of albums.
func (s taskService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

type ByUpdatedAt []obj.Task

func (slice ByUpdatedAt) Len() int           { return len(slice) }
func (slice ByUpdatedAt) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice ByUpdatedAt) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s taskService) Query(ctx context.Context, offset, limit int) ([]obj.Task, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	sort.Stable(ByUpdatedAt(items))
	return items, nil
}

// Update updates the album with the specified ID.
func (s taskService) AddComment(ctx context.Context, id int, req obj.CommentCreationRequest) (obj.Task, error) {
	now := timestamp.Timestamp(timestamp.Now())

	// task, _ := s.Get(ctx, id)

	com := obj.Comment{
		CommentCreationRequest: req,
		CreatedAt:              now,
		Author: user.User{
			ID:       ctx.Value("X-UID").(string),
			Username: ctx.Value("X-User-Name").(string),
		},
	}

	//task.Comments = append(task.Comments, com)
	// s.logger.Infoln("!!!! ", com)
	err := s.repo.AddComment(ctx, id, com)

	// err := s.repo.Update(ctx, task)
	if err != nil {
		return obj.Task{}, errors.Wrap(err, "service work failed")
	}

	return s.repo.Get(ctx, id)
}
