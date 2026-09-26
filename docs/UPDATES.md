# Updates and delivery

How a change in this repo reaches a running machine. Read before adding a
shipped config, a package, or anything a user must receive.

## Current mechanism

Releases are stable tags (`vX.Y.Z`, SemVer); clients never follow `main`.

- **Notice (own act to update, never automatic):** daily idle timer
  `aroli-update-check.timer` writes `~/.cache/aroli-desktop/update.json`;
  login service `aroli-update-notify.service` reads only that cache file
  and shows one notification per release. The bar/notch shows a dot.
- **Update:** `aroli update --dry-run` shows the plan; `aroli update`
  confirms, hands over to the target release's own updater, records the
  rollback point, and re-runs the installer. Local changes are stashed
  automatically and restored afterwards.
- **Rollback:** `aroli rollback` undoes the update.
- **Prune:** `aroli prune --apply` moves files retired by the release to
  the Trash (with backup); never deletes directly, never touches files
  the user modified. Preview without `--apply` deletes nothing.

## Lanes

Aroli owns its own layer and nothing under it:

| Lane | Command | What moves |
| --- | --- | --- |
| Aroli | `aroli update` | rice configs, helpers, CLI, optionals on request |
| System | `pacman -Syu` (yours) | base system and kernel, when you say so |

The project never touches `/boot`, the bootloader, or partitions, and
never moves the kernel. An Aroli release must not decide when the box
rebuilds its boot image, and a box that cannot take a system upgrade
today must still be able to take an Aroli fix (and the reverse).

## The delivery contract

A user-facing file must be delivered by a path a user runs, or it reaches
nobody:

- Shipped configs under `home/` are laid by `install.sh config`
  (symlink or `--copy`) and refreshed by `aroli update`.
- User-owned files (`hypr/user.lua`, `kitty/user.conf`, runtime state in
  `.gitignore`) are seeded once and never re-laid; updates must keep them
  and the tree green.
- A package in `packages/*.txt` must be consumed by the installer phases;
  a script in `scripts/` must be referenced by the installer, `Makefile`,
  CI, CLI, or docs.
- `scripts/verify-delivery.sh` encodes the checkable part of this contract
  (phase parity, no orphan scripts/packages) and runs in `pre-commit` and CI.

Planned (next phase): a declarative `materialize` step (re-lay base,
prune dropped files, overlay user edits last), idempotent `doctor`
reconcilers for stateful drift, and a release `manifest.json` so packages
added to a set reach existing boxes on update instead of fresh installs
only. Until then: every new shipped file needs an explicit delivery path
in the installer or updater, verified by dry-run.
