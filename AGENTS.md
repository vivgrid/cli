# Repository Guidelines

## Project Structure & Module Organization
This repository hosts the `viv` command-line client for Vivgrid. Primary command wiring lives in `cmd/`, where `main.go` registers Cobra commands and helpers such as `ignore_files.go` manage packaging uploads. Shared constants and protocol structures sit in `pkg/`, while user-facing tutorials stay in `docs/` and are regenerated via the CLI. Build artifacts land in `bin/`; keep only reproducible binaries there and exclude them from commits.

## Build, Test, and Development Commands
- `make build`: Compile the CLI into `bin/viv` using the module configuration in `go.mod`.
- `go run ./cmd --help`: Run the latest sources without producing artifacts to verify flag wiring.
- `make doc`: Rebuild the binary and invoke `bin/viv doc` to refresh Markdown references in `docs/`.

## Coding Style & Naming Conventions
All Go sources must stay `gofmt`-clean; run `gofmt -w cmd pkg` or rely on your editor’s Go formatting on save. Organize imports with `goimports` to keep standard, third-party, and local packages grouped. Use PascalCase for exported identifiers, lowerCamelCase for locals, and kebab-case for CLI command names to match existing verbs. Reuse existing Cobra command patterns when introducing new subcommands or flags.

## Testing Guidelines
Add table-driven tests alongside the code under test, e.g., `cmd/feature_test.go`. Execute `go test ./...` before pushing to confirm command behavior and ignore rules pass. When adding new functionality, cover both happy paths and failure cases that surface in deployment workflows. Target meaningful assertions over log output, and prefer fakes over network calls to keep tests hermetic.

## Commit & Pull Request Guidelines
Follow the conventional commit prefixes already in the history (`feat:`, `fix:`, etc.) and write concise, imperative descriptions. Each pull request should describe the user-facing impact, outline test evidence, and link any Vivgrid issue or ticket. Include screenshots or terminal captures only when they clarify a CLI change. Request review once CI and `go test ./...` succeed locally.

## Configuration & Security Notes
Keep secrets out of the repo; load runtime credentials through `vivgrid.yml` or the `VIV_CONFIG_FILE` environment variable. Validate zipper endpoints with `--zipper` flags during development and avoid hard-coding customer-specific hosts. Remove any transient zip bundles or temp files after manual packaging experiments.
