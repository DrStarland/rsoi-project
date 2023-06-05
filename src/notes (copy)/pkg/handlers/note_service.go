package handlers

import (
	"context"
	obj "notes/pkg/models/note"
	"notes/pkg/models/timestamp"
	"notes/pkg/models/user"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Service encapsulates usecase logic for albums.
type NoteService interface {
	Get(ctx context.Context, id int) (obj.Note, error)
	Query(ctx context.Context, offset, limit int) ([]obj.Note, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *obj.NoteCreationRequest) (obj.Note, error)
	Update(ctx context.Context, id int, input *obj.NoteCreationRequest) (obj.Note, error)
	Delete(ctx context.Context, id int) (obj.Note, error)
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

type noteService struct {
	repo   obj.Repository
	logger zap.SugaredLogger
}

// NewService creates a new album service.
func NewNoteService(repo obj.Repository, logger zap.SugaredLogger) noteService {
	return noteService{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s noteService) Get(ctx context.Context, id int) (obj.Note, error) {
	note, err := s.repo.Get(ctx, id)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "service work failed")
	}
	return note, nil
}

// ID              int       `json:"id"`
// Author          user.User `json:"author"` // author
// VisibilityScope scope.Scope
// Tags            []tag.Tag
// CreatedAt       timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
// UpdatedAt       timestamp.Timestamp
// Title           string
// Content         string
// Create creates a new album.

// type NoteCreationRequest struct {
// 	Title           string
// 	Content         string
// 	Tags            []string
// 	VisibilityScope scope.Scope
// }

// type Note struct {
// 	ID     int       `json:"id"`
// 	Author user.User `json:"author"` // author user.User
// 	// VisibilityScope scope.Scope
// 	// Tags            []string            // []tag.Tag
// 	CreatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
// 	UpdatedAt timestamp.Timestamp
// 	NoteCreationRequest
// 	// Title           string
// 	// Content         string
// 	// mu              *sync.Mutex
// }

func (s noteService) Create(ctx context.Context, req *obj.NoteCreationRequest) (obj.Note, error) {
	// if err := req.Validate(); err != nil {
	// 	return obj.Note{}, err
	// }

	now := timestamp.Now()

	insertion := obj.Note{
		Author: user.User{
			ID:       ctx.Value("X-UID").(string),
			Username: ctx.Value("X-User-Name").(string),
		},
		// VisibilityScope:     req.VisibilityScope,
		CreatedAt:           timestamp.Timestamp(now),
		UpdatedAt:           timestamp.Timestamp(now),
		NoteCreationRequest: *req,
	}

	err := s.repo.Create(ctx, &insertion)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "service work failed")
	}
	return s.Get(ctx, insertion.ID)
}

// Update updates the album with the specified ID.
func (s noteService) Update(ctx context.Context, id int, req *obj.NoteCreationRequest) (obj.Note, error) {
	now := timestamp.Timestamp(timestamp.Now())

	note, _ := s.Get(ctx, id)

	note.UpdatedAt = now
	note.NoteCreationRequest = *req

	err := s.repo.Update(ctx, note)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "service work failed")
	}

	return note, nil
}

// Delete deletes the album with the specified ID.
func (s noteService) Delete(ctx context.Context, id int) (obj.Note, error) {
	album, err := s.Get(ctx, id)
	if err != nil {
		return obj.Note{}, errors.New("not exist")
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return obj.Note{}, err
	}
	return album, nil
}

// Count returns the number of albums.
func (s noteService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the albums with the specified offset and limit.
func (s noteService) Query(ctx context.Context, offset, limit int) ([]obj.Note, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return items, nil
}
