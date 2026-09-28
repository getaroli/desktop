# Operating guide for AI agents, Aroli Desktop

This document defines how agents should install, audit, and fix
Aroli Desktop without taking control away from the user. Safety and privacy take
precedence over convenience.

## Non-negotiable principles

1. **Inspect before changing.** Run diagnostics and dry-run mode before
   installing, migrating, enabling services, or editing configurations.
2. **Explicit consent for material effects.** Ask for authorization before
   using sudo/pkexec, installing or removing packages, enabling services,
   writing to `/etc`, changing the wallpaper, re-logging/restarting the
   session, or sending any data over the network.
3. **No automatic bloatware.** The `packages/optional-*.txt` files
   contain optional apps and tools. Each package requires individual
   confirmation; `--yes`, CI, and non-interactive runs must skip them unless
   the user explicitly names each package (or invokes a documented profile
   whose optional package list is shown in full).
4. **Privacy by default.** Do not read, print, store, send, or version
   keys, tokens, passwords, cookies, history, browser profiles, messaging
   data, or files in `~/.ssh`, `~/.gnupg`, `~/.config/gh`, and the like.
5. **Never edit `/usr/share/omarchy/`.** Use only user paths and
   official Omarchy mechanisms.
6. **No destructive actions without a confirmed target and reason.** Prefer
   verifiable backup, Trash, and reversible operations.

## Relevant structure

| Path | Purpose |
| --- | --- |
| `install.sh` | Phased installer; use `--dry-run` first. |
| `diagnose` | Read-only diagnostics. |
| `packages/pacman.txt` | Required aroli base. |
| `packages/optional-*.txt` | Opt-in items, never default. |
| `home/` | Files to be linked or copied into `$HOME`. |
| `home/.local/bin/aroli` | Version, update, rollback, and cleanup CLI (`aroli --help`). |
| `VERSION` / `CHANGELOG.md` | Single source of truth for the version and per-release history. |
| `home/Pictures/wallpapers/` | Default wallpapers shipped by the aroli. |
| `agents/skills/` | AI agent skills, deployed to `~/.agents/skills` by `install.sh`. |
| `OMARCHY.md` | Compatibility contract with Omarchy. |

## Safe install flow

```sh
# 1. Audit with no writes.
./diagnose
./install.sh --dry-run --lang pt-BR

# 2. Only after the user's explicit approval.
./install.sh --lang pt-BR
```

Do not use `--yes` as a substitute for human choice for optional packages. If
the user wants an optional package, show the name, purpose, origin (official
repository or AUR), relevant dependencies, and estimated size before
installing.

## Auditing and fixing

Start with write-free commands:

```sh
omarchy debug --no-sudo --print
hyprctl configerrors
systemctl --user is-active quickshell
./diagnose
git status --short
```

Before fixing, explain the likely cause, target files, expected effect, and
how to revert. After a Hyprland change, validate with `hyprctl reload` and
`hyprctl configerrors`. After changing the shell, check that `quickshell`
is still active. Do not use `omarchy refresh` without confirmation: it
replaces user configurations, although it creates a backup.

## Network and data

- Wallpaper downloads, updates, clones, and AUR queries require
  explicit consent, with ONE declared exception: the daily timer
  `aroli-update-check.timer` (enabled by the installer) transfers a few
  KB (`git ls-remote --tags` or a conditional GET against the GitHub
  releases API, with ETag) once a day and only writes
  `~/.cache/aroli-desktop/update.json`. A login notice
  (`aroli-update-notify.service`, also enabled by the installer) only reads
  that cache file, no network, oneshot, exits in milliseconds, and shows
  at most one desktop notification per release. No telemetry, no identifiers,
  no cached response body. To turn them off:
  `systemctl --user disable aroli-update-check.timer aroli-update-notify.service`.
  `aroli update`, `aroli rollback`, and `aroli prune --apply` never run on their own: they
  require the command (and confirmation, except for the user's explicit
  `--yes`).
- Do not send complete diagnostic reports to external services; strip
  usernames, personal paths, IP addresses, SSIDs, and identifiers.
- Do not add telemetry, analytics, remote plugins, or background
  processes without a clear, reversible user choice.

## Contributions

Preserve [LICENSE](LICENSE), the upstream attribution, and the
Aroli Desktop identity. Every new package must be classified as essential or
optional; when in doubt, treat it as optional.

## Recommended CLI-first workflow

`aroli` is the supported interface for people and agents. Do not invoke
`install.sh` directly unless debugging the installer implementation or the CLI
is unavailable. The CLI can bootstrap the checkout in
`~/.local/share/aroli-desktop`, keeps its binary in `~/.local/bin/aroli`, and
presents the same safe installer phases.

```sh
# First run: inspect only. This does not install packages or edit files.
aroli install --dry-run --lang pt-BR

# After explicit user authorization for package, sudo, and configuration work.
aroli install --lang pt-BR

# Continue an interrupted installation; completed checkpointed phases are skipped.
aroli install --resume

# Read the most recent local install record, or enumerate older records.
aroli logs
aroli logs --list
```

The no-argument `aroli` terminal assistant is appropriate for a user who wants
help choosing an action. For an agent or an unattended instruction, use an
explicit subcommand so the requested side effect is auditable.

### Package-manager policy

On Omarchy, use its package interface only: `omarchy-pkg-add` for official
packages, `omarchy-pkg-aur-add` for AUR packages, and `omarchy-update` for a
full system update. `aroli` detects these commands and delegates to them. Do
not bypass Omarchy with direct `pacman -Syu`, do not edit
`/usr/share/omarchy/`, and do not disable its package hooks. On plain Arch,
`aroli` falls back to `pacman` and `yay`.

Optional software is never part of the normal install result. First show the
catalog, then obtain approval for each selected package:

```sh
aroli plugins list
aroli plugins install --dry-run btop
aroli plugins install btop
```

Treat AUR software as a materially different trust decision. State that it is
from the AUR before executing it, even when `aroli` has already classified it.

## Agent decision procedure

1. Establish whether the aroli is installed with `aroli status`; if it is not,
   use `aroli install --dry-run` and do not clone or install until authorized.
2. Run `aroli diagnose` and the dry-run before a first installation, repair, or
   post-update reapply. Read the plan for disk, network, sudo, package source,
   and `/etc` effects.
3. Ask for a single explicit confirmation that names the material effects:
   package installation, sudo writes, services, network download, and any
   selected optional/AUR package. Never treat “continue” as permission for an
   optional package the user did not name.
4. On success, report the installed aroli version (`aroli status`), the backup
   path printed by the installer when applicable, the install-log path, and
   the required reboot/log-out. Do not reboot, reload Hyprland, or enable a
   new service without separate approval.
5. On failure, stop at the failed phase. Preserve the checkout, backup,
   checkpoint, and log; do not retry a large package transaction blindly.
   Report the exact failed command/phase and offer the smallest safe next
   command below.

## Failure playbook

| Situation | Read-only evidence | Safe next action after approval | Do not do |
| --- | --- | --- | --- |
| Network or mirror download failed | `aroli logs`, `ping -c1 archlinux.org`, package-manager error | Restore connectivity or refresh mirrors, then `aroli install --resume` | Re-run all phases repeatedly or delete package caches without approval. |
| Disk-space preflight failed | `df -h "$HOME"`, `aroli logs` | Ask the user to free at least the stated headroom, then resume | Delete user files, caches, or backups on the agent's initiative. |
| AUR/yay failure | `aroli logs`, `command -v omarchy-pkg-aur-add`, `command -v yay` | Complete only the missing AUR/helper prerequisite, then resume | Replace an AUR package with a different package without consent. |
| Sudo or Omarchy package hook rejected work | Exact command output, `omarchy-version` when available | Explain the required authorization or use the Omarchy package command | Work around Omarchy hooks or write beneath `/usr/share/omarchy/`. |
| Configuration, SDDM, Hyprland, or Quickshell fails after install | `aroli diagnose`, `hyprctl configerrors`, `systemctl --user status quickshell` | Fix the smallest identified cause and re-run its named phase | Run `restore`, remove symlinks, or reboot automatically. |
| Interrupted install | `aroli logs`, `~/.local/state/aroli-desktop/install.checkpoint` | `aroli install --resume` after cause is fixed | Delete the checkpoint; it is the recovery record. |
| CLI update verification fails | `aroli cli update --dry-run`, release/checksum error | Keep the current binary and report the integrity failure | Install an unchecked binary or disable SHA-256 validation. |

`aroli update` updates the aroli checkout and reapplies it. `aroli cli update`
updates only the Go CLI binary after verifying the published SHA-256 checksum.
They are distinct operations; explain which one is needed before running either.

## Repository contract (scale and stability)

These rules keep the project shippable as it grows. See `docs/STRUCTURE.md`,
`docs/DEVELOPMENT.md`, and `docs/UPDATES.md` for the full guides.

1. **Organization is the point.** One purpose per path. Search before adding
   (`rg --files | rg <name>`); reuse, never duplicate.
2. **The repo is the source of truth.** Deployment is one way, repo into
   `$HOME`. Never hand-copy a live tweak back; change the repo and redeploy.
3. **One concern per file (and per package).** `cmd/aroli/main.go` only
   dispatches; command logic lives in `cmd/aroli/internal/<concern>`.
   An installer phase does one thing; a helper script has one job.
4. **Every change must reach users.** A shipped config needs a delivery path
   (`install.sh` phase or updater); a package needs its `packages/*.txt`
   consumed; a script needs a caller. `scripts/verify-delivery.sh` enforces
   the checkable part and runs in `pre-commit` and CI.
5. **Base vs user, separated.** Shipped files may be re-laid; user-owned
   files (`hypr/user.lua`, `kitty/user.conf`, runtime state in `.gitignore`)
   are seeded once and never touched by updates.
6. **Hooks always pass.** Activate with `git config core.hooksPath .githooks`
   and never `--no-verify`. Commits follow Conventional Commits in English.
7. **Safety culture is load-bearing.** No silent optionals, no telemetry,
   no `sudo` surprises, nothing deleted without `--apply`, backups before
   destructive acts. These are features, not friction.
8. **Primary sources first.** Arch Wiki, Hyprland wiki, Quickshell/Qt docs;
   confirm on the running system; prefer the repo's existing pattern.

## Completion criteria
A successful installation is not merely a zero exit code. Confirm all of the
following before reporting completion:

- the requested phases completed and `aroli status` locates the checkout;
- no installer checkpoint remains, unless the user deliberately stopped before
  completion;
- `aroli diagnose` has no new blocking failure attributable to the install;
- the user knows whether to reboot or log out for groups, drivers, or SDDM;
- optional packages were installed only when individually named and approved;
- no secrets, diagnostics containing identifiers, or local logs were sent off
  the machine.
