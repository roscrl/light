package middleware

import (
	"fmt"
	"net/http"

	"github.com/roscrl/light/core/support/contexthelp"
	"github.com/roscrl/light/core/support/rlog"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				var err error
				switch panicType := recovery.(type) {
				case string:
					err = fmt.Errorf(panicType)
				case error:
					err = panicType
				default:
					err = fmt.Errorf("unknown panic: %v", panicType)
				}

				log, rctx := rlog.L(r.Context())
				log.ErrorContext(rctx, "panic", "err", err)

				var requestID string
				if rid, ok := r.Context().Value(contexthelp.RequestIDKey{}).(string); ok {
					requestID = rid
				}

				http.Error(w, "internal server error "+requestID, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
