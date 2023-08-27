package rlog

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/roscrl/light/core/util/contextutil"
)

// L returns the logger from the context along with given context or a default logger and the given context.
func L(r *http.Request) (*slog.Logger, context.Context) {
	rctx := r.Context()

	if logger, ok := rctx.Value(contextutil.RequestLoggerKey{}).(*slog.Logger); ok {
		return logger, rctx
	}

	return NewDefaultLogger(), rctx
}

// LW log with adds new key value arguments to the request scoped logger.
// Updates the given request with the updated logger via context.
// Returns the new logger, the updated request and request context containing the updated logger.
func LW(r *http.Request, args ...any) (*slog.Logger, *http.Request, context.Context) {
	log, rctx := L(r)

	log = log.With(args...)
	rctx = context.WithValue(rctx, contextutil.RequestLoggerKey{}, log)

	r = r.WithContext(rctx)

	return log, r, rctx
}

func NewDefaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
