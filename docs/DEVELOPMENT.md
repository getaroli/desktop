# Development

The loop, the gates, and how to add things without breaking live desktops.

## The loop

```sh
go build -o /tmp/aroli ./cmd/aroli
/tmp/aroli install --dry-run --lang pt-BR   # plan first, always
./diagnose                                   # read-only audit
# ...edit...
./scripts/smoke-test.sh                      # release gate (read-only)
./scripts/test-installer-dry-run.sh          # disposable HOME; proves no dry-run writes
```

Work on a branch (`<type>/<short-topic>`); keep `main` clean. Never run
full installs headless and never pass `--yes` for the user by default
(see `LLMS.md`). Full runs need sudo and a terminal.

## Verify before committing

| What changed | Check |
| --- | --- |
| Shell scripts | `bash -n <file>` |
| Python helpers | `python3 -m py_compile <file>` |
| Go CLI | `go test -race ./...`, `go vet ./...`, `test -z "$(gofmt -l cmd/aroli)"` |
| Installer behavior | `bash scripts/test-installer-dry-run.sh`; also inspect the plan |
| Palette / i18n | Placeholder parity (`{0}`, `{1}`...) and identical key sets across dictionaries |
| Runtime visuals | Switch wallpapers once, confirm bar, terminal, btop, Discord, Spotify follow |
| Everything | `scripts/smoke-test.sh`, `scripts/test-installer-dry-run.sh`, `git diff --check` |

Test behavior, not just syntax. Exercise the change on a running session.

## Adding things

- **A package:** the right list in `packages/` (`pacman` required base,
  `aur`, `gpu`, `extra`, `services`; `optional-*` only with per-item user
  confirmation, never default). Official repos over AUR when both have it.
- **An installer phase:** a `phase_<name>()` function in `install.sh` plus
  the name in `ALL_PHASES` (or `ON_DEMAND_PHASES`); `scripts/verify-delivery.sh`
  enforces the parity. One phase, one concern; do not grow an unrelated phase.
- **A CLI command:** a dispatcher case in `main.go` plus the implementation
  in its `internal/` package (one concern per package; create a new one
  rather than growing an unrelated package).
- **A shipped config:** a file under `home/` reached by `install.sh config`
  (or seeded once from `*.example` when user-owned). A file no install/update
  path lays reaches no user; `scripts/verify-delivery.sh` fails on orphans.
- **A helper script:** `scripts/<name>.sh` referenced from `install.sh`,
  `Makefile`, CI, the CLI, or docs. Unreferenced scripts get deleted, not kept
  "just in case".
- **A doc:** `docs/` in en-US (repo standard per `CONTRIBUTING.md`); UI strings
  stay in QML/dictionaries. One guide per topic; update the layout table in
  `CONTRIBUTING.md` when adding a file.

## Commit gates

Hooks live in `.githooks/`; activate once per checkout:

```sh
git config core.hooksPath .githooks
```

- `commit-msg`: Conventional Commits per `CONTRIBUTING.md`
  (`type(scope): summary`, English, imperative, lowercase, no trailing
  period, ~72 chars). Never `--no-verify`.
- `pre-commit`: `bash -n` on staged shell scripts, `git diff --check`,
  `scripts/verify-delivery.sh`.
- `pre-push`: shellcheck on shell sources when installed (skips cleanly
  when absent).

One logical change per commit. User-facing changes get a `CHANGELOG.md`
entry under `[Unreleased]`; beta releases bump `VERSION` and tag a signed
`v0.1.0-beta.N` (CI checks tag == `VERSION`).

## Research

Primary sources first (Arch Wiki, Hyprland wiki, Quickshell/Qt docs, each
tool's own docs), cross-check anything load-bearing, confirm on the running
system. Prefer the repo's existing pattern over introducing a new one.
