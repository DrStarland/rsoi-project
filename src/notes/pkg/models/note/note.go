package note

import (
	"context"
	"fmt"
	"log"
	"notes/pkg/models/scope"
	"notes/pkg/models/tag"
	"notes/pkg/models/timestamp"
	"os/user"

	// валидатор
	"github.com/asaskevich/govalidator"
)

type Note struct {
	ID              int       `json:"id"`
	Author          user.User `json:"author"` // author
	VisibilityScope scope.Scope
	Tags            []tag.Tag
	Created         timestamp.Timestamp // return time.Now().Format("2006-01-02T15:04:05.000")
	Title           string
	Content         string
	// mu              *sync.Mutex
}

// type Post struct {
// 	ID    string `json:"id" bson:"id"`               // ID
// 	Score int32  `json:"score" bson:"score"`         // score
// 	Type  string `json:"type" valid:"in(link|text)"` // type
// 	Views uint32 `json:"views" bson:"views"`         // views
// 	Title string `json:"title"`                      // title
// 	URL   string `json:"url,omitempty"`              // url
// }

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
	Get(ctx context.Context, id string) (Note, error)
	// Count returns the number of albums.
	Count(ctx context.Context) (int, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Note, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, note Note) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, note Note) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id string) error
}
