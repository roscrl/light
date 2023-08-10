package rlog

import (
	"context"
	"log/slog"
	"sync"

	"github.com/roscrl/light/core/support/contexthelp"
)

const (
	RequestPathLogKey = "request_path"
	RequestIDLogKey   = "request_id"
)

type ContextRequestHandler struct {
	Handler slog.Handler

	attrs []slog.Attr
	mu    sync.Mutex
}

func (h *ContextRequestHandler) Handle(ctx context.Context, record slog.Record) error {
	if path, ok := ctx.Value(contexthelp.RequestPathKey{}).(string); ok {
		record.AddAttrs(slog.String(RequestPathLogKey, path))
	}

	if rid, ok := ctx.Value(contexthelp.RequestIDKey{}).(string); ok {
		record.AddAttrs(slog.String(RequestIDLogKey, rid))
	}

	return h.Handler.Handle(ctx, record)
}

func (h *ContextRequestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *ContextRequestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	h.attrs = append(h.attrs, attrs...)
	h.mu.Unlock()

	return &ContextRequestHandler{
		Handler: h.Handler.WithAttrs(h.attrs),
	}
}

func (h *ContextRequestHandler) WithGroup(name string) slog.Handler {
	return &ContextRequestHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
