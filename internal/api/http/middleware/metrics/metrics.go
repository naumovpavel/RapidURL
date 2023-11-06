package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func New(
	totalRequests *prometheus.CounterVec,
	responseTime *prometheus.GaugeVec,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t0 := time.Now()

			lrw := NewLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)

			deltaT := time.Now().Sub(t0)
			totalRequests.With(prometheus.Labels{
				"method":      r.Method,
				"url":         r.URL.String(),
				"status_code": strconv.Itoa(lrw.statusCode),
			}).Add(1)
			responseTime.With(prometheus.Labels{
				"method":      r.Method,
				"url":         r.URL.String(),
				"status_code": strconv.Itoa(lrw.statusCode),
			}).Set(float64(deltaT))
		}

		return http.HandlerFunc(fn)
	}
}
