package middleware

import (
	"context"
	"net/http"

	"gateway/pkg/services"

	"go.uber.org/zap"
)

func Auth(next http.HandlerFunc, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := session.ContextWithSession(r.Context(), sess)
		ctx := context.TODO()
		// next(w, r.WithContext(ctx))
		if token, err := services.RetrieveToken(w, r); err == nil {
			r.Header.Set("X-User-Name", token.Subject)
			next(w, r.WithContext(ctx))
		}
	}
}
