# Contributing to olap-sql

Thank you for your interest in contributing! This guide covers everything you need to get up and running, understand the project, and submit quality changes.

---

## Table of Contents

- [Development environment](#development-environment)
- [Project layout](#project-layout)
- [Running tests](#running-tests)
- [Making changes](#making-changes)
- [Pull request checklist](#pull-request-checklist)
- [Commit message style](#commit-message-style)
- [Code style](#code-style)
- [Adding a new database backend](#adding-a-new-database-backend)
- [Reporting bugs and requesting features](#reporting-bugs-and-requesting-features)

---

## Development environment

### Prerequisites

| Tool | Minimum version | Notes |
|------|----------------|-------|
| Go   | 1.22           | Project uses range-over-integer (Go 1.22) and built-in `min`/`max` (Go 1.21). |
| Git  | any recent     |       |
| SQLite (optional) | — | Only needed to run the SQLite-backed integration tests locally. |
| ClickHouse / MySQL / PostgreSQL (optional) | — | Only needed for the respective backend integration tests. |

### Clone and build

```bash
git clone https://github.com/AWaterColorPen/olap-sql.git
cd olap-sql
go mod tidy
go build ./...
```

No additional setup steps are required for the core library. The project has no generated code and no `Makefile` targets.

---

## Project layout

```
.
├── api/
│   ├── models/         # TOML schema structs (Dictionary data model)
│   └── types/          # Public query/result/filter types used by callers
├── docs/               # User-facing documentation (Markdown)
│   └── superpowers/    # Internal design specs and iteration plans
├── test/               # Integration test fixtures (TOML configs, SQL seeds)
├── client.go           # Database client abstraction and GORM wiring
├── configuration.go    # Configuration and DBOption types
├── database.go         # Low-level query execution helpers (RunSync, RunChan)
├── dependency_graph.go # Metric dependency resolution (METRIC_DIVIDE, etc.)
├── dictionary.go       # Dictionary: loads and caches the TOML schema
├── dictionary_adapter.go   # Adapter layer (TOML file → in-memory model)
├── dictionary_column.go    # Column resolution helpers
├── dictionary_splitter.go  # Multi-source JOIN splitting
├── dictionary_translator.go # Query → Clause translation (core logic)
├── manager.go          # Public Manager API
└── run.go              # Result assembly helpers
```

The **core translation pipeline** lives in `dictionary_translator.go` → `clause.go` (in `api/types`). If you want to understand how a `Query` becomes SQL, start there and follow the `Translate` → `Statement` call chain.

---

## Running tests

### Unit and SQLite integration tests (no external DB required)

```bash
go test ./...
```

All tests that use SQLite run with an in-memory database; no setup is needed.

### Running a specific test

```bash
go test -run TestManager ./...
```

### Running tests with verbose output

```bash
go test -v ./...
```

### Integration tests with other backends

Tests that target ClickHouse, MySQL, or PostgreSQL are skipped automatically when the corresponding DSN environment variable is not set. To run them, set the relevant variable before running `go test`:

```bash
# ClickHouse example
CLICKHOUSE_DSN="clickhouse://localhost:9000/default" go test ./...
```

Check the test file comments in `manager_test.go` for the exact env var names per backend.

---

## Making changes

1. **Fork** the repository and create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-change-description
   ```
2. Make your changes. Keep each PR **focused on a single concern** — one PR per feature, refactor, or bug fix.
3. Add or update **tests** for your change. New functionality without tests will not be merged.
4. Run the test suite and confirm it passes:
   ```bash
   go test ./...
   ```
5. Run `gofmt` and `go vet`:
   ```bash
   gofmt -w .
   go vet ./...
   ```
6. Open a pull request against `main`.

### When to open an issue first

For **large or design-changing** contributions (new backends, query language extensions, API changes), please open an issue to discuss the approach before writing code. This avoids wasted effort if the direction doesn't align with the project.

For **small, self-contained** changes (typo fixes, documentation improvements, minor bug fixes), a PR without a prior issue is fine.

---

## Pull request checklist

Before marking a PR ready for review, confirm the following:

- [ ] All existing tests pass (`go test ./...`).
- [ ] New tests cover the added or changed behaviour.
- [ ] `gofmt` has been run; no formatting changes in the diff.
- [ ] `go vet` reports no issues.
- [ ] Public types and functions have godoc comments.
- [ ] The PR description explains *what* changed and *why*.
- [ ] For breaking changes: a `CHANGELOG.md` entry is included.

---

## Commit message style

Follow the [Conventional Commits](https://www.conventionalcommits.org/) convention:

```
<type>(<optional scope>): <short description>

[optional body]
```

Common types:

| Type       | Use for                                          |
|------------|--------------------------------------------------|
| `feat`     | New features                                     |
| `fix`      | Bug fixes                                        |
| `refactor` | Code changes that neither add features nor fix bugs |
| `docs`     | Documentation only                              |
| `test`     | Adding or fixing tests                           |
| `chore`    | Dependency updates, build changes, CI config     |

**Examples:**

```
feat(filter): add FILTER_OPERATOR_HAS for ClickHouse array columns
fix(translator): handle nil TimeInterval without panicking
docs: add API reference page
chore: upgrade gorm to v1.25
```

---

## Code style

- Follow standard Go conventions as enforced by `gofmt` and `go vet`.
- Prefer **explicit error handling** over panics.
- Use the `any` alias instead of `interface{}`.
- Prefer stdlib (`slices`, `maps`) over third-party utility packages for simple helpers.
- Keep package names short and lowercase; avoid underscores in package names.
- Exported symbols must have godoc comments; unexported helpers are encouraged but not required to.

---

## Adding a new database backend

olap-sql currently supports ClickHouse, MySQL, PostgreSQL, and SQLite. If you want to add a new backend:

1. Add a new `DBType` constant in `api/types/db_type.go`.
2. Implement the GORM driver initialisation in `client.go` (see the existing `newGormDB` switch statement).
3. If the new database uses dialect-specific SQL functions (e.g. array operations, date truncation), add handling in `dictionary_translator.go` where existing dialect branches exist.
4. Add integration test coverage under `test/` with a representative fixture.
5. Update `README.md` to list the new backend under **Requirements**.

Please open an issue first to discuss the new backend — we want to ensure test infrastructure and CI are set up correctly.

---

## Reporting bugs and requesting features

- **Bugs:** Open a GitHub issue with a minimal reproducible example — include your TOML schema, the `Query` you built, and the SQL or error you got vs. what you expected.
- **Features:** Open a GitHub issue describing the use case and your proposed API. For query language extensions, include the SQL you'd like to generate.
