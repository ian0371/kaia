# AGENTS.md

## Setup

- Go version: use the version declared in `go.mod`
- Build: `make all`
- Test: `make test`
- Lint: `go run build/ci.go lint -v --new-from-rev=dev`
- Full lint: `make lint`
- Vet: `go vet ./...`
- Dev tools for generators: `make devtools`

## Architecture

- CLI entrypoints: `cmd/`
- Consensus engines and protocols: `consensus/`
- Node runtime and services: `node/`, `networks/`
- Core blockchain, state, and VM logic: `blockchain/`, `snapshot/`, `storage/`
- Kaia feature modules: `kaiax/`
- Integration and protocol-heavy tests: `tests/`
- Contract sources and generated bindings: `contracts/`
- Do not hand-edit generated files such as files marked `Code generated`, `gen_*.go`, or generated outputs under `contracts/artifacts/` and `contracts/typechain-types/`

## Coding Rules

- Keep patches minimal and localized to the task
- Ask before adding, removing, or updating dependencies
- Preserve consensus, RPC, config, wire-format, and CLI compatibility unless the task explicitly requires a breaking change
- Prefer updating generator inputs over editing generated outputs directly
- Add or update focused tests for behavior changes
- Use commit and PR titles in the form `<area>: description`

## Validation

- Run `gofmt -w`, `gofumpt -w`, and `goimports -w` on modified Go files
- Run `make test` before finishing
- Run `go run build/ci.go lint -v --new-from-rev=dev` before finishing
- Run `make all` and `go vet ./...` when touching shared runtime, build, consensus, networking, or storage code
- Run area-specific suites when relevant: `make test-node`, `make test-networks`, `make test-datasync`, `make test-tests`, `make test-others`
- If a required check cannot be run, explain exactly why

## PR / Contribution Behavior

- Target the `dev` branch
- Link the related issue for feature work
- Keep diffs reviewable; split large changes when practical
- Summarize changed areas and rationale
- Mention tradeoffs or compatibility risks
- Add or update tests for behavior changes
- First-time contributors must sign the CLA in the PR comment as described in `CONTRIBUTING.md`
