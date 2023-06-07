package scope

// валидатор

type AffiliationType int

const (
	User AffiliationType = iota + 1
	Group
)

type Scope struct {
	ScopeAffiliation AffiliationType
	OwnerID          int
}

// type Repository interface {
// 	GetByUsername(flightNumber string) ([]*Ticket, error)
// 	Add(*Ticket) error
// 	Delete(ticketUID string) error
// }
