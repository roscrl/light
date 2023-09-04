package rlog

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/utils/contextutil"
)

var (
	resultLog *slog.Logger
	resultCtx context.Context
)

func BenchmarkNewLoggerEveryRequest(b *testing.B) {
	b.ReportAllocs()

	var (
		logger *slog.Logger
		ctx    context.Context
	)

	for i := 0; i < b.N; i++ {
		ctx = context.Background()
		requestID := "123"

		textHandler := slog.NewTextHandler(os.Stdout, nil)
		requestContextHandler := ContextRequestHandler{Handler: textHandler}

		logger = slog.New(&requestContextHandler)

		ctx = context.WithValue(ctx, contextutil.RequestLoggerKey{}, logger)
		ctx = context.WithValue(ctx, contextutil.RequestIDKey{}, requestID)
	}

	resultCtx = ctx
	resultLog = logger
}

func BenchmarkCloningWithLoggerEveryRequest(b *testing.B) {
	b.ReportAllocs()

	var (
		logger *slog.Logger
		ctx    context.Context
	)

	textHandler := slog.NewTextHandler(os.Stdout, nil)
	requestContextHandler := ContextRequestHandler{Handler: textHandler}
	logger = slog.New(&requestContextHandler)

	for i := 0; i < b.N; i++ {
		ctx = context.Background()
		requestID := "123"

		logger2 := logger.With(slog.String(key.RequestID, requestID))

		ctx = context.WithValue(ctx, contextutil.RequestLoggerKey{}, logger2)
	}

	resultCtx = ctx
	resultLog = logger
}
