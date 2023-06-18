package services

import (
	"context"
	obj "costs/pkg/models/cost"
	"costs/pkg/models/timestamp"
	"costs/pkg/models/user"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Service encapsulates usecase logic for albums.
type CostService interface {
	Get(ctx context.Context, id int) (obj.Cost, error)
	Query(ctx context.Context, offset, limit int) ([]obj.Cost, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *obj.CostCreationRequest) (obj.Cost, error)
	Update(ctx context.Context, id int, input *obj.CostCreationRequest) (obj.Cost, error)
	Delete(ctx context.Context, id int) (obj.Cost, error)
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

type costService struct {
	repo   obj.Repository
	logger zap.SugaredLogger
}

// NewService creates a new album service.
func NewCostService(repo obj.Repository, logger zap.SugaredLogger) costService {
	return costService{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s costService) Get(ctx context.Context, id int) (obj.Cost, error) {
	cost, err := s.repo.Get(ctx, id)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "service work failed")
	}
	return cost, nil
}

// type CommentCreationRequest struct {
// 	Content string
// }

// type Comment struct {
// 	Author    user.User
// 	UpdatedAt timestamp.Timestamp
// 	CommentCreationRequest
// }

// type CostCreationRequest struct {
// 	Title           string
// 	Description     string
// 	NoteURL         string
// 	VisibilityScope scope.Scope
// }

// type Cost struct {
// 	ID        int       `json:"id"`
// 	Author    user.User `json:"author"` // author user.User
// 	Status    string
// 	UpdatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
// 	UpdatedAt timestamp.Timestamp
// 	Comments  []Comment

// 	CostCreationRequest
// }

func (s costService) Create(ctx context.Context, req *obj.CostCreationRequest) (obj.Cost, error) {
	// if err := req.Validate(); err != nil {
	// 	return obj.Cost{}, err
	// }

	now := timestamp.Now()

	insertion := obj.Cost{
		Author: user.User{
			ID:       ctx.Value("X-UID").(string),
			Username: ctx.Value("X-User-Name").(string),
		},
		// VisibilityScope:     req.VisibilityScope,
		UpdatedAt:           timestamp.Timestamp(now),
		CreatedAt:           timestamp.Timestamp(now),
		CostCreationRequest: *req,
	}

	err := s.repo.Create(ctx, &insertion)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "service work failed")
	}
	return insertion, nil
}

// Update updates the album with the specified ID.
func (s costService) Update(ctx context.Context, id int, req *obj.CostCreationRequest) (obj.Cost, error) {
	now := timestamp.Timestamp(timestamp.Now())

	cost, _ := s.Get(ctx, id)

	cost.UpdatedAt = now
	cost.CostCreationRequest = *req

	err := s.repo.Update(ctx, cost)
	if err != nil {
		return obj.Cost{}, errors.Wrap(err, "service work failed")
	}

	return cost, nil
}

// Delete deletes the album with the specified ID.
func (s costService) Delete(ctx context.Context, id int) (obj.Cost, error) {
	album, err := s.Get(ctx, id)
	if err != nil {
		return obj.Cost{}, errors.New("not exist")
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return obj.Cost{}, err
	}
	return album, nil
}

// Count returns the number of albums.
func (s costService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

type ByUpdatedAt []obj.Cost

func (slice ByUpdatedAt) Len() int           { return len(slice) }
func (slice ByUpdatedAt) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice ByUpdatedAt) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s costService) Query(ctx context.Context, offset, limit int) ([]obj.Cost, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	sort.Stable(ByUpdatedAt(items))
	return items, nil
}
