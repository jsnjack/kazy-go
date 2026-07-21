package cmd

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// LevelTrace is the slog level used for maximum-detail diagnostic logging.
const LevelTrace = slog.Level(-8)

// L is the package-wide structured logger, configured by initLogger.
var L *slog.Logger

// initLogger configures the package-wide logger based on the --debug and
// --trace flags. When tracePath is set, trace-level output is written to
// that file (truncated on every call). Otherwise output goes to stderr at
// the given level. It returns a cleanup function that must be called before
// the program exits.
func initLogger(tracePath, level string) func() {
	var w io.Writer = os.Stderr
	cleanup := func() {}
	if tracePath != "" {
		f, err := os.OpenFile(tracePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err == nil {
			w = f
			cleanup = func() {
				if err := f.Close(); err != nil {
					slog.Log(context.Background(), LevelTrace, "close trace file", "err", err)
				}
			}
		}
	}
	lvl := slog.LevelWarn
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "trace":
		lvl = LevelTrace
	}
	h := slog.NewTextHandler(w, &slog.HandlerOptions{Level: lvl})
	L = slog.New(h)
	slog.SetDefault(L)
	return cleanup
}
