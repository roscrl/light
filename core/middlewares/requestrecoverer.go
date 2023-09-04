package middlewares

import (
	"fmt"
	"net/http"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/utils/contextutil"
)

func RequestRecoverer(next http.Handler) http.Handler {
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

				log, rctx := rlog.L(r)
				log.ErrorContext(rctx, "panic", key.Err, err)

				requestID := contextutil.RequestID(r)

				http.Error(w, "internal server error "+requestID, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
