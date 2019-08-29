package boil

import (
	"context"
	"io"
)

// WithDebug modifies a context to configure debug writing.
func WithDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, ctxDebug, debug)
}

// IsDebug returns true if the context has debugging enabled, or
// the value of DebugMode if not set.
func IsDebug(ctx context.Context) bool {
	debug, ok := ctx.Value(ctxDebug).(bool)
	if ok {
		return debug
	}
	return DebugMode
}

// WithDebugWriter modifies a context to configure the writer written to
// when debugging is enabled.
func WithDebugWriter(ctx context.Context, writer io.Writer) context.Context {
	return context.WithValue(ctx, ctxDebugWriter, writer)
}

// DebugWriterFrom returns the debug writer for the contexts, or DebugWriter
// if not set.
func DebugWriterFrom(ctx context.Context) io.Writer {
	writer, ok := ctx.Value(ctxDebugWriter).(io.Writer)
	if ok {
		return writer
	}
	return DebugWriter
}
