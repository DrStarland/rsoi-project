package note

import (
	"context"
	"fmt"
	"gateway/pkg/models/scope"
	"gateway/pkg/models/timestamp"
	"gateway/pkg/models/user"
	"log"

	// валидатор
	"github.com/asaskevich/govalidator"
)

type NoteCreationRequest struct {
	Title           string
	Content         string
	Tags            []string
	VisibilityScope scope.Scope
}

type Note struct {
	ID     int       `json:"id"`
	Author user.User `json:"author"` // author user.User
	// VisibilityScope scope.Scope
	// Tags            []string            // []tag.Tag
	CreatedAt timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
	UpdatedAt timestamp.Timestamp
	NoteCreationRequest
	// Title           string
	// Content         string
	// mu              *sync.Mutex
}

func (Note) TableName() string {
	return "note"
}

func (p *Note) Validate() error {
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
	Get(ctx context.Context, id int) (Note, error)
	// Count returns the number of albums.
	Count(ctx context.Context) (int, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Note, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, note *Note) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, note Note) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id int) error
}
