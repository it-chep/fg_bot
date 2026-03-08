# Repository Guidelines

## Project Structure & Module Organization
- `cmd/main.go`: service entrypoint.
- `internal/app.go` and `internal/init.go`: app wiring, dependency initialization, startup flow.
- `internal/modules/bot`: bot domain logic.
- `internal/modules/bot/action`: command handlers (`/init_fg`, stats, report save, ping toggle).
- `internal/modules/bot/worker`: background jobs (including report reminders via cron-like worker pool).
- `internal/server`: HTTP/webhook server and handlers.
- `internal/pkg`: shared infrastructure (logger, telegram client, transactions, worker pool).
- `tools/migrations`: Goose SQL migrations.
- `bin/`: local tooling binaries (for example `goose`).

## Build, Test, and Development Commands
- `make deps`: install local tools into `./bin` (notably `goose`).
- `make infra`: start local infrastructure with Docker Compose.
- `make minfra-up`: apply DB migrations to local DB.
- `make minfra-down`: rollback/reset migrations.
- `make build`: build binary to `bin/app`.
- `go test ./...`: quick default test/build verification.
- `make test` and `make test-cover`: run tests with race/coverage using `gotestsum` (if installed).

Example local bootstrap:
```bash
make deps
make infra
make minfra-up
go test ./...
```

## Coding Style & Naming Conventions
- Language: Go 1.24+.
- Always run `gofmt` (and preferably `go test ./...`) before opening a PR.
- Package names are lower_snake/lowercase; exported identifiers use `CamelCase`.
- Keep handlers/actions thin; move SQL to `dal` packages and domain behavior to module-level services/workers.
- Prefer explicit, context-aware logging through `internal/pkg/logger`.

## Testing Guidelines
- Primary framework: Go `testing` package.
- Place tests next to code as `*_test.go`.
- Name tests clearly by behavior, e.g. `TestGetDailyStats_EmptyReports`.
- For DB-dependent behavior, cover query edge cases (date boundaries, duplicates, idempotency).
- Minimum pre-merge check: `go test ./...` must pass.

## Commit & Pull Request Guidelines
- This checkout does not include `.git` history; use Conventional Commit style:
  - `feat(bot): add 21:00 reminder worker`
  - `fix(dal): remove timezone from SQL layer`
- PR should include:
  - short problem statement,
  - summary of changes by module/path,
  - migration notes (`tools/migrations/*`) if schema changed,
  - test evidence (command + result).

## Security & Configuration Tips
- Configuration is loaded from `.env` (`internal/config/config.go`).
- Never commit real tokens/passwords.
- Required runtime vars include bot token, DB credentials, webhook flags, and port settings.
