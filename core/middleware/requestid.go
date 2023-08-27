package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/roscrl/light/core/util/contextutil"
)

func RequestID(next http.Handler) http.Handler {
	const requestIDSize = 8

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes := make([]byte, requestIDSize)
		if _, err := rand.Read(bytes); err != nil {
			bytes = []byte("00000000")
		}

		requestID := fmt.Sprintf("%X", bytes)

		rctx := context.WithValue(r.Context(), contextutil.RequestIDKey{}, requestID)

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}
