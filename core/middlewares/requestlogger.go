package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/roscrl/light/core/helpers/applog"
	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/utils/contextutil"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := applog.NewDefaultLogger()
		requestContextHandler := rlog.ContextRequestHandler{
			Handler: handler.Handler(),
		}

		log := slog.New(&requestContextHandler)

		rctx := context.WithValue(r.Context(), contextutil.RequestLoggerKey{}, log)

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}
