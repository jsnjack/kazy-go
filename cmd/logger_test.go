package cmd

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stderr = old

	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func TestInitLoggerDefaultLevel(t *testing.T) {
	out := captureStderr(t, func() {
		cleanup := initLogger("", "")
		defer cleanup()
		slog.Warn("visible warning")
		slog.Debug("hidden debug")
	})

	if !strings.Contains(out, "visible warning") {
		t.Errorf("expected warning to be logged, got: %q", out)
	}
	if strings.Contains(out, "hidden debug") {
		t.Errorf("expected debug to be filtered out by default, got: %q", out)
	}
}

func TestInitLoggerDebugLevel(t *testing.T) {
	out := captureStderr(t, func() {
		cleanup := initLogger("", "debug")
		defer cleanup()
		slog.Debug("now visible debug")
	})

	if !strings.Contains(out, "now visible debug") {
		t.Errorf("expected debug to be logged, got: %q", out)
	}
}

func TestInitLoggerTraceFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trace.log")
	if err := os.WriteFile(path, []byte("stale content\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cleanup := initLogger(path, "trace")
	slog.Log(context.Background(), LevelTrace, "trace message")
	cleanup()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "stale content") {
		t.Errorf("expected trace file to be truncated, got: %q", content)
	}
	if !strings.Contains(string(content), "trace message") {
		t.Errorf("expected trace message to be written, got: %q", content)
	}
}
