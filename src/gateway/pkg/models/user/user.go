package user

import (
	"errors"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// func (us *User) SetUpdated(val uint32) {
// 	us.updated = sql.NullString{String: strconv.Itoa(int(val))}
// }

//go:generate mockgen -source=user.go -destination=repo_mock.go -package=user UserRepo
type UserRepo interface {
	Register(login, pass string) (*User, error)
	Authorize(login, pass string) (*User, error)
}

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invald password")
)
