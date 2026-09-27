# AGENTS.md

Guidance for AI agents (and humans) working in this repository.

## Project overview

`envd` is a command-line tool that builds container-based development
environments. Users describe their environment in a `build.envd` file written
in a Python-like DSL ([Starlark](https://github.com/bazelbuild/starlark)).
envd interprets the file, registers the declared state into an in-memory
graph, compiles the graph into [buildkit LLB](https://github.com/moby/buildkit#exploring-llb),
and builds the image via buildkitd or Docker (moby).

The full documentation lives in a separate repo:
[github.com/tensorchord/envd-docs](https://github.com/tensorchord/envd-docs)
(see `docs/developers/development.md` there for the long-form development
tutorial this file summarizes).

## Repository layout

- `cmd/envd/` — the main CLI entry point.
- `cmd/envd-sshd/` — the sshd binary injected into built environments; not
  used directly by end users.
- `pkg/` — all of the implementation:
  - `app/` — CLI command definitions and wiring (urfave/cli).
  - `autocomplete/` — bash/zsh completion.
  - `builder/` — the core builder: interprets the Starlark manifest and
    compiles the resulting graph to LLB, then drives the build.
  - `buildkitd/` — client code for connecting to the buildkitd container.
  - `config/` — envd configuration handling.
  - `data/` — data-source support (e.g. downloading remote files into images).
  - `driver/` — abstraction over container engines (Docker, nerdctl).
  - `editor/` — VS Code / Jupyter integration.
  - `envd/` — the envd engine: abstraction to create, start, and manage envd
    environments.
  - `flag/` — shared viper/CLI flag definitions.
  - `home/` — management of XDG directories (`~/.config/envd`,
    `~/.cache/envd`).
  - `lang/` — the build language. `frontend/starlark/` interprets
    `build.envd` (v1 syntax under `v1/`); `ir/` holds the intermediate
    graph representation (`ir/v1`). **Add new `build.envd` functions here**:
    the Starlark builtin in `lang/frontend/starlark/v1/<module>/`, the IR
    entry point in `lang/ir/v1/interface.go`, and the compile step in
    `lang/ir/v1/`.
  - `metrics/` — build progress metrics widgets.
  - `progress/` — build progress UI.
  - `remote/` — sshd implementation for remote development.
  - `shell/` — shell (zsh) integration.
  - `ssh/` — SSH client used by `envd up` to attach to environments.
  - `syncthing/` — syncthing-based file synchronization.
  - `types/` — shared types used across packages.
  - `util/` — utility helpers (fileutil, netutil, starlarkutil, ...).
  - `version/` — version variables populated at build time via ldflags.
- `envd/` — a fake Python package (`envd/api/v0`, `envd/api/v1`) whose
  docstrings generate the `build.envd` API reference. **Update the docstring
  here whenever you change a `build.envd` function's signature or behavior**;
  these files are linted/formatted with ruff, not compiled.
- `base-images/` — Dockerfiles and build scripts for the base images
  (`envd`, `envd-sshd`, `remote-cache`).
- `examples/` — example `build.envd` projects.
- `e2e/` — end-to-end tests: `cli/`, `language/`, `docs/` (doc examples).
  They build real images and require a Docker daemon.
- `docs/` — not user docs; only design proposals and README assets.

## Makefile targets

The Makefile builds both binaries (`envd`, `envd-sshd`) and stamps version
info via ldflags. Common targets:

| Target | Purpose |
| --- | --- |
| `make` / `make build-local` | Build both binaries into `./bin` (default) |
| `make dev` | Build and reinstall the local Python wheel for debugging |
| `make debug` / `make debug-local` | Build debug binaries into `./debug-bin` (for VS Code attach) |
| `make test` | Generate mocks, then run unit tests with race detector + coverage |
| `make test-local` | Unit tests without mock generation |
| `make lint` | Run golangci-lint (auto-installs it into GOPATH/bin) |
| `make vet` / `make fmt` | `go vet` / `go fmt` |
| `make e2e-test` | All e2e tests (needs Docker; ~20 min timeout) |
| `make e2e-cli-test` / `make e2e-lang-test` / `make e2e-doc-test` | e2e subsets |
| `make envd-lint` / `make envd-fmt` | Ruff lint/format for `envd/api` Python stubs |
| `make generate` | Regenerate mocks (mockgen) |
| `make addlicense` | Add the Apache license header to Go files |
| `make clean` | Remove `bin/`, `debug-bin/`, `dist/`, etc. |

## Development workflow

1. Prerequisites: Docker (20.10+) and a recent Go toolchain (see the `go`
   directive in `go.mod`).
2. Build and smoke-test:
   ```bash
   go mod tidy
   make
   ./bin/envd bootstrap
   ./bin/envd version
   ```
3. Before sending a PR: `make lint`, `make test`, and `make envd-lint` if you
   touched `envd/api`.

## Conventions

- **Commits and PR titles** follow [Conventional Commits](https://www.conventionalcommits.org/)
  (`feat:`, `fix:`, `chore(deps):`, `docs:`, ...) — enforced by CI
  (`.github/semantic.yml`).
- **DCO sign-off is required**: commit with `git commit -s`.
- New Go files need the Apache license header (`make addlicense`).
- Tests: unit tests live next to the code (`*_test.go`); anything that needs
  a Docker daemon belongs in `e2e/`.
- Dependency updates: dependabot runs monthly with a cooldown (7 days
  default, 3 days for patch releases — see `.github/dependabot.yml`).
