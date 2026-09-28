---
name: aroli-install
description: >
  Use for installing, updating, or maintaining Aroli Desktop:
  install.sh phases, dry-run, Aroli CLI, backups, rollback, and beta releases.
  Triggers: install, update, upgrade, cleanup, leftovers, backups, rollback,
  restore, phases, dry-run, aroli CLI, install.sh, config phase, spicetify.
---

# Install and update

Entry points: `./install.sh` (phases, default `--link`), `aroli`
(`status | check | update | rollback | prune`), and `./diagnose`
(read-only).

## Operating rules

- `./install.sh -n` (dry-run) and `./diagnose` before anything that writes.
  Full runs need sudo and a terminal: never run them headless, and never pass
  `-y` for the user by default.
- `config` lays `home/` into `$HOME` (symlinks, with `save_aside` backups to
  `~/.dotfiles-backup`) and deploys `agents/skills` to `~/.agents/skills`.
- `spicetify` needs `/opt/spotify` writable (ACL, sudo once) and applies
  `current_theme aroli color_scheme pywal`. See `aroli-spotify` when it
  complains.
- `final` regenerates the palette from the first wallpaper alphabetically;
  re-apply the user's wallpaper afterwards with `set-wallpaper.sh`.
- `aroli prune` previews retired managed files; `--apply` moves them to Trash
  with backup. Never delete `~/.dotfiles-backup` dirs unasked.
- `VERSION` is the single source of truth; releases use `v0.1.0-beta.N` tags
  and GitHub pre-releases (CI attaches binaries and checksums).
