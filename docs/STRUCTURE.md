# Repository structure

One job per path. If what you need already exists somewhere, reference
it; do not keep two copies. See `DEVELOPMENT.md` for the workflow and
`UPDATES.md` for how a change reaches an installed machine.

## Top level

| Path | One job |
| --- | --- |
| `install.sh` | Phased installer (default `--link`); owns `home/` deployment, packages, services |
| `diagnose` | Read-only diagnostics; never gains write behavior |
| `cmd/aroli/` | `aroli` Go CLI: public interface over the installer backend |
| `home/` | Deployed tree, mirrored into `$HOME` by `install.sh config` |
| `packages/` | Package lists (`pacman`, `aur`, `gpu`, `extra`, `services`, `optional-*`, `git-repos`) |
| `system/` | System-level files outside `$HOME` (currently `system/etc/`) |
| `scripts/` | Standalone helpers (smoke test, font sync, delivery check) |
| `bin/` | Go build output (gitignored; `make build`) |
| `docs/` | Contributor and user guides |
| `agents/skills/` | AI agent skills, deployed to `~/.agents/skills` by `install.sh` |
| `.github/workflows/` | CI (`ci`, `fonts`, `release`, `auto-assign`) |
| `.githooks/` | Local commit/push gates; activate with `git config core.hooksPath .githooks` |
| `VERSION` / `CHANGELOG.md` | Single source of truth for the version and release history |
| `Makefile` | `build`, `test`, `smoke`, `install-cli`, `clean` |
| `LLMS.md` | Operating contract for automations and AI agents |
| `OMARCHY.md` | Compatibility contract with Omarchy |
| `CONTRIBUTING.md` | Collaboration language, commit style, testing checklist |

## `home/` (the deployed tree)

Mirrors `$HOME`. `install.sh config` lays it via symlinks (default) or
copies (`--copy`), backing up replaced files to `~/.dotfiles-backup`.

- `home/.config/hypr/` — compositor config (`hyprland.lua` entry point, not
  `hyprland.conf`); `user.lua` / `user.conf` are seeded once from
  `*.example` and then user-owned (untracked, never re-laid).
- `home/.config/quickshell/` — bar, notch, panels, Settings UI, translations
  (`translations-es.js`, `translations-pt-BR.js`; en-US source in QML).
- `home/.local/bin/` — user-level helpers (`aroli-backend`,
  wallpaper/palette/battery/gaming scripts). No root installs.
- `home/.config/systemd/user/` — `aroli-update-check.{service,timer}` and
  `aroli-update-notify.service` (idle, oneshot, cache-read only).
- `home/Pictures/wallpapers/` — bundled Aroli Backdrops.

Runtime state that the rice rewrites on its own (palettes, effects toggles,
derived new-tab theme, `user.lua`/`user.conf`, lock-screen language) is
untracked in `.gitignore` with a reason comment. Versioned files next to
them are seeds and defaults, never live state.

## `cmd/aroli/` (the CLI)

`main.go` is only the dispatcher, the help text, and the version
(`-X main.version=` stays working). Command implementations live in
`internal/`, one concern per package: `sys` (shared primitives),
`install`, `plugins`, `selfupdate`, `logs`, `profiles`, `snapshot`,
`wallpaper`, `gaming`, `reading`, `power`, `session`, `recover`,
`prefs`, `doctor`, `status`, `tui`. Tests live next to their package.
The CLI owns the public interface and the beginner-friendly menu; the
installer backend keeps the idempotency and backup rules.

## Rules

1. **The repo is the source of truth.** Deployment is one way: repo into
   `$HOME`. Never hand-copy a live tweak back; change the repo and redeploy.
2. **User-owned files are never re-laid.** Seeds (`*.example`) copy once;
   runtime state is untracked. Tag checkouts and `aroli update` must keep
   the tree green for users with local edits.
3. **No second copies.** One package list per purpose, one helper per job,
   one language table per language. Search (`rg --files | rg <name>`)
   before adding.
4. **Never ship into a path the user's own tools write.** App defaults that
   an app rewrites (e.g. `mimeapps.list`) go one layer down or become seeds;
   the shipped tree must not clobber user choice on update.
