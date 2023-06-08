package task

import (
	"context"
	"fmt"
	"log"
	"sync"
	"tasks/pkg/models/scope"
	"tasks/pkg/models/timestamp"
	"tasks/pkg/models/user"

	// валидатор
	"github.com/asaskevich/govalidator"
)

const (
	StatusQueued   = "В очереди"
	StatusInWork   = "В работе"
	StatusDone     = "Выполнена"
	StatusCanceled = "Отменена"
	StatusWaiting  = "Ожидание"
)

type CommentCreationRequest struct {
	Content string
}

type Comment struct {
	Author    user.User
	CreatedAt timestamp.Timestamp
	CommentCreationRequest
}

type TaskCreationRequest struct {
	Title           string
	Description     string
	NoteURL         string
	VisibilityScope scope.Scope
}

type Task struct {
	ID        int       `json:"id"`
	Author    user.User `json:"author"` // author user.User
	Status    string
	CreatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
	UpdatedAt timestamp.Timestamp
	Comments  []Comment
	mu        *sync.Mutex

	TaskCreationRequest
}

func (Task) TableName() string {
	return "task"
}

func (p *Task) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrs.Errors() {
				data := []byte(fmt.Sprintf("field: %#v\n\n", fld))
				log.Println(data)
			}
		}
	}
	return err
}

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id int) (Task, error)
	// Count returns the number of albums.
	Count(ctx context.Context) (int, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Task, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, task *Task) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, task Task) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id int) error

	AddComment(ctx context.Context, id int, comment Comment) error
}
