package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/utils/contextutil"
)

func RequestLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		textHandler := slog.NewTextHandler(os.Stdout, nil)
		requestContextHandler := rlog.ContextRequestHandler{
			Handler: textHandler,
		}

		log := slog.New(&requestContextHandler)

		rctx := context.WithValue(r.Context(), contextutil.RequestLoggerKey{}, log)

		next.ServeHTTP(w, r.WithContext(rctx))
	}
}
