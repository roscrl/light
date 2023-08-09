package rlog

import (
	"context"
	"log/slog"
	"os"

	"github.com/roscrl/light/core/support/contexthelp"
)

// L returns the logger from the context along with given context or a default logger and the given context.
func L(ctx context.Context) (*slog.Logger, context.Context) {
	if logger, ok := ctx.Value(contexthelp.RequestLogger{}).(*slog.Logger); ok {
		return logger, ctx
	}

	return NewDefaultLogger(), ctx
}

func NewDefaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
