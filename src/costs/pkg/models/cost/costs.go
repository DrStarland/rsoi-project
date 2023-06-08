package cost

import (
	"context"
	"costs/pkg/models/scope"
	"costs/pkg/models/timestamp"
	"costs/pkg/models/user"
	"fmt"
	"log"

	// валидатор
	"github.com/asaskevich/govalidator"
)

type CostCreationRequest struct {
	Amount          float32
	Comment         string
	Category        string
	VisibilityScope scope.Scope
}

type Cost struct {
	ID        int       `json:"id"`
	Author    user.User `json:"author"` // author user.User
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	CostCreationRequest
}

type IncomeCreationRequest struct {
	Amount          float32
	Comment         string
	VisibilityScope scope.Scope
}

type Income struct {
	ID        int       `json:"id"`
	Author    user.User `json:"author"` // author user.User
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	IncomeCreationRequest
}

func (Cost) TableName() string {
	return "Costs"
}

func (p *Cost) Validate() error {
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
	Get(ctx context.Context, id int) (Cost, error)
	// Count returns the number of albums.
	Count(ctx context.Context) (int, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]Cost, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, cost *Cost) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, cost Cost) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id int) error
}
