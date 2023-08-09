package rlog

import (
	"context"
	"log/slog"

	"github.com/roscrl/light/core/support/contexthelp"
)

type ContextRequestHandler struct {
	slog.Handler
}

const (
	RequestPathLogKey = "request_path"
	RequestIDLogKey   = "request_id"
)

func (h ContextRequestHandler) Handle(ctx context.Context, record slog.Record) error {
	if path, ok := ctx.Value(contexthelp.RequestPathKey{}).(string); ok {
		record.AddAttrs(slog.String(RequestPathLogKey, path))
	}

	if rid, ok := ctx.Value(contexthelp.RequestIDKey{}).(string); ok {
		record.AddAttrs(slog.String(RequestIDLogKey, rid))
	}

	return h.Handler.Handle(ctx, record)
}
