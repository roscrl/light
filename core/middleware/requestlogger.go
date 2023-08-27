package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/roscrl/light/core/support/rlog"
	"github.com/roscrl/light/core/util/contextutil"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		textHandler := slog.NewTextHandler(os.Stdout, nil)
		requestContextHandler := rlog.ContextRequestHandler{
			Handler: textHandler,
		}

		log := slog.New(&requestContextHandler)

		rctx := context.WithValue(r.Context(), contextutil.RequestLoggerKey{}, log)

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}
