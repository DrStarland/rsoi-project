package services

import (
	"context"
	obj "costs/pkg/models/income"
	"costs/pkg/models/timestamp"
	"costs/pkg/models/user"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Service encapsulates usecase logic for albums.
type IncomeService interface {
	Get(ctx context.Context, id int) (obj.Income, error)
	Query(ctx context.Context, offset, limit int) ([]obj.Income, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *obj.IncomeCreationRequest) (obj.Income, error)
	Update(ctx context.Context, id int, input *obj.IncomeCreationRequest) (obj.Income, error)
	Delete(ctx context.Context, id int) (obj.Income, error)
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

type incomeService struct {
	repo   obj.Repository
	logger zap.SugaredLogger
}

// NewService creates a new album service.
func NewIncomeService(repo obj.Repository, logger zap.SugaredLogger) incomeService {
	return incomeService{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s incomeService) Get(ctx context.Context, id int) (obj.Income, error) {
	income, err := s.repo.Get(ctx, id)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "service work failed")
	}
	return income, nil
}

// type CommentCreationRequest struct {
// 	Content string
// }

// type Comment struct {
// 	Author    user.User
// 	CreatedAt timestamp.Timestamp
// 	CommentCreationRequest
// }

// type IncomeCreationRequest struct {
// 	Title           string
// 	Description     string
// 	NoteURL         string
// 	VisibilityScope scope.Scope
// }

// type Income struct {
// 	ID        int       `json:"id"`
// 	Author    user.User `json:"author"` // author user.User
// 	Status    string
// 	CreatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
// 	CreatedAt timestamp.Timestamp
// 	Comments  []Comment

// 	IncomeCreationRequest
// }

func (s incomeService) Create(ctx context.Context, req *obj.IncomeCreationRequest) (obj.Income, error) {
	// if err := req.Validate(); err != nil {
	// 	return obj.Income{}, err
	// }

	now := timestamp.Now()

	insertion := obj.Income{
		Author: user.User{
			ID:       ctx.Value("X-UID").(string),
			Username: ctx.Value("X-User-Name").(string),
		},
		// VisibilityScope:     req.VisibilityScope,
		CreatedAt:             timestamp.Timestamp(now),
		IncomeCreationRequest: *req,
	}

	err := s.repo.Create(ctx, &insertion)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "service work failed")
	}
	return insertion, nil
}

// Update updates the album with the specified ID.
func (s incomeService) Update(ctx context.Context, id int, req *obj.IncomeCreationRequest) (obj.Income, error) {
	now := timestamp.Timestamp(timestamp.Now())

	income, _ := s.Get(ctx, id)

	income.CreatedAt = now
	income.IncomeCreationRequest = *req

	err := s.repo.Update(ctx, income)
	if err != nil {
		return obj.Income{}, errors.Wrap(err, "service work failed")
	}

	return income, nil
}

// Delete deletes the album with the specified ID.
func (s incomeService) Delete(ctx context.Context, id int) (obj.Income, error) {
	album, err := s.Get(ctx, id)
	if err != nil {
		return obj.Income{}, errors.New("not exist")
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return obj.Income{}, err
	}
	return album, nil
}

// Count returns the number of albums.
func (s incomeService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

type ByCreatedAtIn []obj.Income

func (slice ByCreatedAtIn) Len() int           { return len(slice) }
func (slice ByCreatedAtIn) Less(i, j int) bool { return slice[i].CreatedAt < slice[j].CreatedAt }
func (slice ByCreatedAtIn) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s incomeService) Query(ctx context.Context, offset, limit int) ([]obj.Income, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	sort.Stable(ByCreatedAtIn(items))
	return items, nil
}
