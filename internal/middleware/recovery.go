package middleware

import (
	"TransportLayer/internal/entity"
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, entity.ErrInternal, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
