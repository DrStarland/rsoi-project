package middleware

import (
	"net/http"
	"notes/pkg/services"

	"go.uber.org/zap"
)

func Auth(next http.HandlerFunc, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if token, err := services.RetrieveToken(w, r, logger); err == nil {
			r.Header.Set("X-User-Name", token.Subject)
			r.Header.Set("X-UID", token.UID)
			next(w, r)
		}
	}
}
