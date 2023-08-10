package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/roscrl/light/core/support/contexthelp"
	"github.com/roscrl/light/core/support/rlog"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := r.Context()

		textHandler := slog.NewTextHandler(os.Stdout, nil)
		requestContextHandler := rlog.ContextRequestHandler{
			Handler: textHandler,
		}

		logger := slog.New(&requestContextHandler)

		rctx = context.WithValue(rctx, contexthelp.RequestLoggerKey{}, logger)

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}
