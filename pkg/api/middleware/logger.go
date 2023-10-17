package middleware

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				status := ww.Status()
				reqLogger := logger.With(
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("requestID", middleware.GetReqID(r.Context())),
					zap.Duration("elapsed", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
				ref := ww.Header().Get("Referer")
				if ref == "" {
					ref = r.Header.Get("Referer")
				}
				if ref != "" {
					reqLogger = reqLogger.With(zap.String("ref", ref))
				}
				ua := ww.Header().Get("User-Agent")
				if ua == "" {
					ua = r.Header.Get("User-Agent")
				}
				if ua != "" {
					reqLogger = reqLogger.With(zap.String("userAgent", ua))
				}

				switch {
				case status >= 200 && status < 300:
					reqLogger.Info(fmt.Sprintf("%d OK", status))
				case status >= 300 && status < 400:
					reqLogger.Info(fmt.Sprintf("%d Redirect", status))
				case status >= 400 && status <= 499:
					reqLogger.Info(fmt.Sprintf("%d Client error", status))
				case status >= 500:
					reqLogger.Error(fmt.Sprintf("%d Server error", status))
				default:
					reqLogger.Info(fmt.Sprintf("%d Unknown", status))
				}
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
