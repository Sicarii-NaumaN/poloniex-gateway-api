package middleware

import (
	"net/http"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"

	"go.uber.org/zap"
)

// PanicMid logges error if handler errors
func PanicMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Error(err.(error).Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
