package session

import (
	"net/http"
)

//go:generate mockgen -source=manager.go -destination=session_mock.go -package=session SessionsManager
type SessionsManager interface {
	Check(r *http.Request) (*Session, error)
	// Create(user *user.User) (*Session, error)
	// DestroyCurrent(w http.ResponseWriter, r *http.Request)
}
