package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.WithFields(logrus.Fields{
						"error":      err,
						"method":     r.Method,
						"path":       r.URL.Path,
						"stack":      string(debug.Stack()),
						"user_agent": r.UserAgent(),
						"remote_addr": r.RemoteAddr,
					}).Error("Panic recovered")

					// Return 500 error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					errorResponse := map[string]interface{}{
						"error": "Internal server error",
						"code":  "INTERNAL_ERROR",
					}

					json.NewEncoder(w).Encode(errorResponse)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}