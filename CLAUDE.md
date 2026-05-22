# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development workflow

For every task, always follow this git workflow:

1. Create a new branch from `main` with a descriptive name (`feature/`, `fix/`, `chore/` prefix).
2. Implement the changes with one or more commits.
3. Push the branch to the remote repository (`git push -u origin <branch>`).
4. Open a pull request targeting `main` via `gh pr create`.

Never commit directly to `main`.

## Commands

```bash
make run          # run locally (requires .env)
make test         # run all tests
go test ./internal/application/greet/...  # run a single package's tests

make build-init   # build Docker image + create network (first-time setup)
make http-up      # start service via Docker
make http-down    # stop Docker container
```

Copy `.env.example` to `.env` before running. GitHub credentials are required to pull the private `go-devkit` module.

## Architecture

This service follows Clean Architecture with three layers:

**Domain** (`internal/domain/`) — pure business entities, no dependencies.

**Application** (`internal/application/`) — use cases that orchestrate domain logic. Each use case receives a `DTO`, runs a `validator.Validator[DTO]` before executing, and returns an output DTO. Validators are composed of individual rule functions and built with `validator.New[T](rules...)`.

**Infrastructure** (`internal/infra/`) — adapters for the outside world:
- `http/handlers/` — HTTP handlers implementing `httpserver.Handler` (returns `error`). Handlers call use cases and write JSON via `infrahttp.WriteJson`.
- `http/routes/` — groups routes by feature and wires up handler factories.
- `http/` (package level) — error handling pipeline. `ErrorHandler` is a chain of `ErrorWriter` implementations; unmatched errors fall back to 500. Validation errors (`validator.Error`) map to 422.
- `clock/` — `Clock` interface abstracts `time.Now()` for testability.

**Entry point** (`cmd/http/main.go`) — bootstraps config (viper/.env), logger, metrics (StatsD/Datadog), and the HTTP server from `go-devkit/pkg/httpserver` (chi-based). Routes are registered once at startup.

## Adding a new endpoint

1. Define a `DTO` and output struct in `internal/application/<feature>/`.
2. Implement `UseCase` interface with a `Run(ctx, dto)` method; add validators via `newXxxValidator()`.
3. Create a handler in `internal/infra/http/handlers/` that reads path params with `httpserver.GetParam(req, "param")`.
4. Register the route in `internal/infra/http/routes/` using `httpserver.NewGet/Post/...`.
5. Append the route group in `cmd/http/main.go` `bootstrapRoutes()`.

## Error handling

Return errors from handlers — do not write the response directly on error. The error handler chain in `infrahttp.ErrorHandler` matches and serializes them. To handle a new error type, implement `ErrorWriter` and register it in `error_handler_config.go`.

## Testing conventions

- **One test file per source file** — `foo.go` gets `foo_test.go`. Never group tests from multiple source files into one test file.
- **External test package** — use `package foo_test`, not `package foo`.
- **Table-driven tests** — use a slice of structs with named subtests via `t.Run()`.
- **Assertions** — use `testify/assert` and `testify/require`.
- **Mocks and helpers** — define them in the same test file that uses them, not in shared files.