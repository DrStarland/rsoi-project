package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func Auth(next http.HandlerFunc, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-User-Name", "mamont")
		r.Header.Set("X-UID", "2228745g")
		next(w, r)

		// if token, err := services.RetrieveToken(w, r, logger); err == nil {
		// 	r.Header.Set("X-User-Name", token.Subject)
		// 	r.Header.Set("X-UID", token.UID)
		// 	next(w, r)
		// }
	}
}
