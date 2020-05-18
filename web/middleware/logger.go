package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/adeki/go-utils/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		xForwarded := strings.TrimSpace((strings.Split(r.Header.Get("X-Forwarded-For"), ","))[0])
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		fields := logger.Fields{
			"ip":              ip,
			"method":          r.Method,
			"proto":           r.Proto,
			"path":            r.URL.String(),
			"status":          lrw.status,
			"ua":              r.UserAgent(),
			"referer":         r.Referer(),
			"elapsed":         time.Since(start).Milliseconds(),
			"x-forwarded-for": xForwarded,
		}

		logger.WithFields(fields).Info("access")
	}
	return http.HandlerFunc(fn)
}
