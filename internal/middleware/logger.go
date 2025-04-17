package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/solumD/auth-test-task/internal/logger"
	"go.uber.org/zap"
)

// Middleware оборачивает запрос в логгер с детальной информацией
func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Info("New request",
				zap.String("user_agent", r.UserAgent()),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.Info("request completed",
					zap.Int("status", ww.Status()),
					zap.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
