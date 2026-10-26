---
name: go-modern
description: Write modern, idiomatic Go code. Use whenever writing, modifying, fixing, refactoring, or reviewing Go code. Detects the project's Go version from go.mod and applies version-aware modern idioms (Go 1.0–1.27, based on the JetBrains Modern Go Guidelines) — prefers slices/maps/cmp helpers, range-over-int, min/max builtins, errors.AsType, and other features available in the target version over outdated patterns, even when surrounding code uses legacy style.
---

# Modern Go

Write modern, idiomatic Go. Models generate outdated Go because of training-data lag and frequency bias; this skill is the explicit reference that fixes both.

## Workflow

1. Detect the Go version: read the `go` directive in `go.mod` (fallback: `go version`). Use only language features and stdlib additions available in that version or earlier.
2. Before editing Go code, scan `references/features.md` for guidelines relevant to the code at hand. It lists every guideline ordered newest first, grouped by category (Loops, Collections, Strings, Errors, JSON, Sync, Context, Types, Testing, HTTP, Time, Fmt, URL, Utilities), each with impact rating and before/after examples.
3. Apply every guideline whose Go version is ≤ the project version. Prefer modern idioms even when nearby code uses an older pattern. Skip a guideline only if it would not compile, would change behavior, or clearly does not apply.
4. Respect project conventions found in `AGENTS.md` / `README.md` (module path, lint config, layout) — they override defaults, not this skill's idioms.

## High-impact idioms (quick reference)

- `for i := range n` instead of `for i := 0; i < n; i++` (1.22)
- `min`/`max` builtins instead of if-else chains (1.21)
- `slices.Contains`, `slices.IndexFunc`, `slices.SortFunc`, `slices.Clone` instead of manual loops and copies (1.21)
- `slices.Collect`, `slices.Sorted`, `maps.Keys`/`maps.Values` iterators (1.23)
- `any` instead of `interface{}` (1.18)
- `errors.Is`/`errors.As` for matching; `errors.AsType[T]` for type-safe matching (1.26); wrap with `%w`
- `strings.Cut`, `strings.CutPrefix`, `strings.CutSuffix` instead of `Index`+slicing (1.18/1.20)
- `cmp.Or` instead of chains of zero-value checks; `cmp.Compare` for ordering (1.22)
- `new(expr)` to take a pointer to a value (1.26)
- `sync.WaitGroup.Go` instead of `Add`/`go func`/`Done` (1.25); `sync.OnceFunc`/`OnceValue` (1.21)
- `t.Context()` in tests instead of `context.Background()` (1.24); `t.Chdir`, `t.TempDir` instead of `os` calls
- `fmt.Appendf` instead of `append([]byte(fmt.Sprintf(...)))` (1.19)

Full details, rationale, and counter-examples: `references/features.md`.

## General style

- Format with gofumpt (gofmt-compatible, stricter); group imports with goimports, project packages last.
- Doc comment on every exported identifier, starting with the identifier's name.
- Errors are values: package-level sentinel errors + `errors.Is`/`As`; no `panic` outside `main`.
- Table-driven tests in an external `_test` package; `t.Parallel()` where safe.
- After edits, run the project's checks — typically `golangci-lint fmt`, `golangci-lint run`, `go test ./...`, or `make check` — and fix all findings instead of suppressing them.

## Attribution

Based on the JetBrains Modern Go Guidelines (Apache-2.0). Full text: `references/features.md`; license: `LICENSE.jetbrains`.
