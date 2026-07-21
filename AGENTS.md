# AGENTS.md

> See [AGENTS.universal.md](./AGENTS.universal.md) and [AGENTS.go.md](./AGENTS.go.md) for universal conventions.
> Refresh: `make standards`

---

## Overview

kažy reads lines from STDIN and highlights, filters, or extracts substrings
matching one or more patterns. It's meant to sit in the middle of a shell
pipeline, so its own STDOUT carries only the processed data.

---

## Architecture

```
main.go                    Entry point; delegates to cmd.Execute()
cmd/
  root.go                  Cobra root command: flags, RunE, logger wiring
  root_utils.go            compileRegExp, matchesRegExpList, limitLine, processData
  root_utils_test.go       Tests and benchmarks for root_utils.go
  colours.go               ANSI colour codes cycled across patterns
  logger.go                slog setup for --debug / --trace
  logger_test.go           Tests for initLogger
```

---

## Key Flows

1. `rootCmd.RunE` compiles the positional args (highlight patterns) and the
   `--include` / `--exclude` flag values into regexp lists via
   `compileRegExp`.
2. `processData` scans STDIN line by line: applies include/exclude filters,
   an optional length limit, then either extracts one matched substring per
   line (`--extract`) or wraps every match in ANSI colour codes, and prints
   the result to STDOUT.

---

## Build & Run

```bash
make build      # bin/kazy plus per-platform binaries
make test       # go test ./...
make check      # fmt, vet, build, test, lint

echo "hello 123 world" | bin/kazy 123
echo "hello 123 world" | bin/kazy -x 123
```

---

## Design Decisions

- STDOUT carries only processed data; all diagnostics go through `slog` to
  STDERR (default/`--debug`) or the trace file (`--trace`) so they never mix
  into piped output.
- Colours cycle through a fixed ANSI palette (`terminalColours` in
  `colours.go`) indexed by pattern position, giving each pattern a stable,
  distinct colour.
- The number of patterns is capped at `len(terminalColours)` since every
  pattern needs its own colour slot.

---

## Gotchas

- `bufio.Scanner`'s buffer caps line length; `--buffer` raises it, but a line
  still longer than the buffer aborts the scan (logged via `slog`, not
  printed to STDOUT) rather than being truncated silently.
