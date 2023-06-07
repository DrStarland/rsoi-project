package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func AccessLog(next http.HandlerFunc, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infoln("ЕЩЁ ЖИВ")
		start := time.Now()
		next(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	}
}
