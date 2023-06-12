package middleware

import (
	"gateway/pkg/models/statistic"
	"gateway/pkg/services"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

func AccessLog(next http.HandlerFunc, logger *zap.SugaredLogger, topic string, producer sarama.SyncProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)

		start := time.Now()
		next(lrw, r)
		duration := time.Since(start)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", duration,
		)

		stat := &statistic.RequestStat{}
		stat.Path = r.URL.Path
		stat.Method = r.Method
		stat.UserName = r.Header.Get("X-User-Name")
		stat.StartedAt = start
		stat.Duration = duration
		stat.ResponceCode = lrw.Status()

		go services.SendRequestStatToKafka(stat, topic, producer)
	}
}
