# Updates and delivery

How a change in this repo reaches a running machine. Read before adding a
shipped config, a package, or anything a user must receive.

## Current mechanism

Releases are beta pre-release tags (`v0.1.0-beta.N`); clients never follow
`main` or stable release endpoints.

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

`aroli update` records its progress in `~/.cache/aroli-desktop/update.json`.
After the checkout changes, it creates a pre-update snapshot, runs
`aroli materialize` (re-lay base, prune only still-managed removals, then
overlay `~/.config/aroli/user_edits`), runs the read-only `doctor`, and
creates a post-update snapshot. A failed stage remains visible in
`update.json`; the checkout is retained so `aroli rollback` remains explicit.

`aroli cli update` accepts beta `v0.1.0-beta.N` release tags and checks the binary
against the release's `SHA256SUMS.txt` before an atomic replacement. The
checksum is served beside the binary by the same GitHub release, so it detects
corruption but is not independent publisher authentication. See
[`SECURITY.md`](../SECURITY.md) for the threat boundary and
[`docs/RELEASE-CHECKLIST.md`](RELEASE-CHECKLIST.md) for release verification.

`doctor --fix` owns the four idempotent local reconcilers: UserEdits,
MimeDefaults, ShellLoad and Manifest. Package convergence is separately
confirmed because it may reach pacman and sudo.
