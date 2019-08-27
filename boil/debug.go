package boil

import (
	"context"
	"io"
)

type debugConfig struct {
	debug  bool
	writer io.Writer
}

// ContextDebug modifies a context to configure the debug output, overriding
// the global DebugMode and DebugWriter. If writer is nil, then DebugWriter
// will still be used.
func ContextDebug(ctx context.Context, debug bool, writer io.Writer) context.Context {
	return context.WithValue(ctx, ctxDebug, &debugConfig{
		debug:  debug,
		writer: writer,
	})
}

// IsDebug checks to see if debugging is enabled, and returns the correct
// debug writer to use for outputting.
func IsDebug(ctx context.Context) (bool, io.Writer) {
	config, _ := ctx.Value(ctxDebug).(*debugConfig)

	switch {
	case config == nil:
		return DebugMode, DebugWriter
	case !config.debug:
		return false, nil
	case config.writer != nil:
		return true, config.writer
	default:
		return true, DebugWriter
	}
}
