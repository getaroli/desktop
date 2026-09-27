# Changelog, Aroli Desktop

All notable changes to this rice are documented here, newest first.
The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and version numbers follow [Semantic Versioning](https://semver.org/):

- `MAJOR`: breaking change (new Hyprland requirement, removed phase,
  renamed deployed path, config format that old installs cannot read).
- `MINOR`: new feature that stays compatible (new panel, script, package,
  translation, optional).
- `PATCH`: fix or docs with no behavior change for existing installs.

## [Unreleased]

## [3.4.6] - 2026-09-26

### Fixed

- No user-facing changes. Carries the `main` sync (rebase-copies from
  the v3.4.4 promotion) so fresh tags descend from `main` again and the
  automatic promotion assert passes.

## [3.4.5] - 2026-09-26

### Fixed

- No user-facing changes. Validation release exercising the automatic
  tag promotion (`promote-main`) after enabling Actions PR creation in
  the organization.

## [3.4.4] - 2026-09-26

### Fixed

- Delivery CI lints QML through `.qmllint.ini` at the repo root instead
  of version-specific flags (Qt 6.4 knows neither `--max-warnings` nor
  the needed `--bare` behavior). Semantic warnings are informational
  because `Quickshell` and `QtQuick.Effects` cannot resolve on noble;
  syntax errors still fail the gate.

## [3.4.3] - 2026-09-26

### Fixed

- Delivery CI lints QML with Qt's default imports instead of `--bare`
  (which hid even `QtQuick` and the builtins, failing with exit 255),
  installs the `QtQuick` QML modules, and tolerates import warnings for
  modules unavailable on noble's Qt 6.4 (`Quickshell`,
  `QtQuick.Effects`) via `--max-warnings -1`. Syntax errors still fail.

## [3.4.2] - 2026-09-26

### Fixed

- Delivery CI resolves `qmllint`/`qml6lint` from the Qt bin directories
  (`/usr/lib/qt6/bin`), where Ubuntu installs it outside `PATH`, instead
  of failing with `xargs: qmllint: No such file or directory`.

## [3.4.1] - 2026-09-26

### Fixed

- Delivery CI now runs profile dry-runs with an explicit pacman shim, supplies
  the container user expected by the installer, and limits ShellCheck to
  errors so existing advisory warnings remain visible without blocking a
  release.

## [3.4.0] - 2026-09-26

### Added

- Delivery contract: `aroli materialize` re-lays managed configuration,
  preserves local overrides, records the delivered inventory and backs up
  retired managed files. `aroli update` now captures pre/post snapshots and
  records delivery progress in `update.json`.
- Release channels: stable follows signed release tags; `unstable-dev` follows
  its integration branch. The package manifest is generated from
  `packages/*.txt`; `aroli verify` compares it to the installed box.
- Promotion gate: releases open a fast-forward-only, rebase PR to `main` and
  require the release tag check. Branch protection rejects direct, non-linear
  and force updates.
- Delivery CI: ShellCheck, QML lint, profile dry-run matrix, Arch container
  interface check, artifact verification, and the delivery-contract gate.
- `aroli doctor --explain`, delivery progress in `aroli status`, and local
  translation overrides persisted outside the managed shell tree.

### Added

- Delivery-contract foundations: `docs/STRUCTURE.md` (repo map, one job per
  path), `docs/DEVELOPMENT.md` (dev loop, gates, how to add things) and
  `docs/UPDATES.md` (how a change reaches installed machines, lanes, contract).
  Local gates in `.githooks/` (`commit-msg` enforces Conventional Commits,
  `pre-commit` checks shell syntax plus whitespace, `pre-push` runs shellcheck
  when present; activate with `git config core.hooksPath .githooks`) and
  `scripts/verify-delivery.sh` (installer phase parity, no orphan helper
  scripts or package lists), wired into `pre-commit` and the next CI pass.
- CLI split into `cmd/aroli/internal/`, one package per concern (`sys`,
  `install`, `plugins`, `selfupdate`, `logs`, `profiles`, `snapshot`,
  `wallpaper`, `gaming`, `reading`, `power`, `session`, `recover`, `prefs`,
  `doctor`, `status`, `tui`); `cmd/aroli/main.go` is now only the dispatcher,
  help text, and version. No behavior change; tests moved next to their
  packages. Removed the dead `guidedInstall`/`guidedPlugins` helpers the TUI
  had replaced.
- Cardinal rules for scale and stability added to `LLMS.md` (repository
  contract: organization, source of truth, delivery, hooks, safety).
- Aroli para o Brave Origin (`~/.config/aroli-newtab`): tela inicial
  (hero, busca, zonas de atalhos, tarefas e atividade GitHub) com o
  Encaixe e o avatar oficiais, aba de ajustes oculta na engrenagem
  (tema, buscador, nome, GitHub, atalho) e dropdown próprio nos tokens
  Aroli. A chave Sistema / Aroli Dark / Aroli Black repinta o navegador
  inteiro via `chrome.theme` (frame, toolbar, abas, omnibox, incógnito).
  O Encaixe acompanha o acento do tema; animações (entrada, relógio,
  ajustes, itens) via anime.js vendorizado, sem CDN. Sem seleção de
  texto com o mouse; atividade em cache (1 fetch/hora, sem token). O tema
  é derivado uma vez por wallpaper por `aroli-newtab-pywal.sh` e
  registrado via `brave-origin-flags.conf` pelo `install.sh config`.
- Quadrado Máquina ao lado do ano (CPU, memória, disco, temperatura,
  bateria, uptime via host nativo `aroli-sys.py`); o Brave abre direto
  na página Aroli por policy gerenciada.
- `scripts/sync-aroli-fonts.sh` (`--check`/`--apply`/`--install`) com
  verificação (allowlist, magic OTTO, piso de tamanho, família Aroli via
  fc-scan) e lock em `scripts/aroli-fonts.lock`; workflow `fonts`
  valida os OTFs em todo PR e abre PR de revisão quando o upstream
  publica builds novos (nunca push direto na main).

### Changed

- Aroli Mono NF e Aroli Sans sincronizadas com o upstream `aroli@main`
  (0.905: família óptica de pontos unificada, rabicho da vírgula/ponto
  e vírgula redesenhados). `scripts/sync-aroli-fonts.sh --install`
  aplica no `$HOME`; `./install.sh config` + `final` também.
- `brave-origin-bin` replaces `google-chrome` in the optional plugins
  and becomes the default handler for web links.

## [3.3.0] - 2026-09-21

### Changed

- The CLI is now `aroli`, the product name. `rice` stays as a
  deprecation shim (warning + handover), same treatment the
  `umbra-cursor-size` rename got. Backend dispatcher renamed to
  `aroli-backend` with a silent repo-level `rice` shim so installers
  predating the rename keep updating; the updater fetches
  `aroli-linux-*` assets while releases keep publishing the
  `rice-linux-*` names through the transition. Units, shell helpers,
  diagnose, docs and help speak `aroli`.

## [3.2.0] - 2026-09-21

### Added

- Aroli Sans and Aroli Mono NF ship with the rice
  (`home/.local/share/fonts/`, installed by the `config` phase) and
  Aroli Sans is a UI font choice in Settings › Appearance. Both are
  original Aroli drawings from the upstream `aroli` repo.
- Launcher uses the Encaixe product mark (U+100000 glyph from
  Aroli Mono NF) tinted with the contrast-guaranteed `onWallAccent`
  ink instead of the distro logo. Measured 5.95:1 against the
  wallpaper (house guarantee is 3.5:1). Image colorization was
  discarded: it keeps the source lightness and can never hit the
  guaranteed ink.
- Lyrics in MediaPanel (Super+D): synced follow from lrclib with a
  calm centered glide, spectrum as loading/fallback.
- Tokens: `fsCaption`/`fsTitle`, `mFollow`/`mTint`, `hitPad`,
  `inkHi`/`inkMid`/`inkLo`; type, color, radius, spacing and motion
  migrated to tokens across control-center panels and Settings.

### Changed

- Aroli for Spotify theme replaces termspot (own theme: rounded
  cards and covers, hairlines, system type, wallpaper accent).
- Product naming: "Rice de Eduardo" strings become Aroli Desktop.

### Fixed

- `fmt()` prints `--:--` for bogus durations instead of digit walls.

## [3.1.0] - 2026-09-21

### Added

- Spotify follows the wallpaper again, end to end. Menu launches get
  `--remote-debugging-port` from the new
  `home/.local/share/applications/spotify.desktop` override (shipped like
  `vlc.desktop`, laid down by the `config` phase), and `launch-spotify.sh`
  carries the same flags for keybind/terminal launches. Without the port,
  every wallpaper change failed its live push silently and Spotify looked
  stuck on a stale scheme.
- `spicetify-pywal.sh` no longer swallows failures: a stale backup
  (`refresh` exiting 0 while compiling nothing) log to
  `~/.cache/spicetify-pywal.log` and raises a critical notification, and a
  failed live push is logged and notified before the restart fallback.
  `phase_spicetify` also verifies the desktop override carries the port.
- Notch artwork is a true circle in both modes (compact 22px, expanded
  34px). `clip: true` crops to a rectangle even with radius, so both
  artworks now use a real `MultiEffect` mask, the same pattern as
  `MediaPanel`. The compact media row also keeps the battery row's 15px
  right inset, so the edge stops jumping when music starts.

### Fixed

- `install.sh` re-rates pacman mirrors at most weekly instead of on every
  run (re-rating only burned minutes on timing-out mirrors).
- `install.sh config` seeds `poke-theme/state` creating its directory
  first, and no longer copies `__pycache__` into `~/.local/bin`.

## [3.0.0] - 2026-09-21

### Changed

- Product rebrand: Umbra Noctis is now **Aroli Desktop**
  (`github.com/eduardoaugustolb/aroli-desktop`), the desktop product of the
  Aroli family. The `rice` binary name is unchanged. `install.sh config`
  migrates `~/.local/share/umbra-noctis`, `~/.local/state/umbra-noctis` and
  `~/.cache/umbra-noctis` to their `aroli-desktop` counterparts (nothing is
  copied when the new location already exists), and `rice export`/`import`
  reads the legacy `umbra-preferences.json` when `aroli-preferences.json`
  is absent. `AROLI_DESKTOP_REPO` replaces `UMBRA_RICE_REPO` (still honored
  as a fallback).

Historical entries below keep the Umbra Noctis name as published.

- Follow the upstream `umbra` → `aroli` rename
  (`github.com/eduardoaugustolb/aroli`): the `cursor` phase clones the new
  repo, builds `themes/cursor/aroli` and installs `~/.local/share/icons/Aroli`
  (`Aroli Pointer`). Configs (`hyprland.lua`, `hyprland.conf`, GTK 2/3/4)
  point at `Aroli`; the legacy `~/.local/share/icons/Umbra` is moved to
  `Umbra.bak` after a successful install.
- Bundled wallpapers renamed to Aroli Backdrops following upstream:
  `aroli-ember-coast.png`, `aroli-obsidian-dunes.png`,
  `aroli-silent-threshold.png`, `aroli-black-mountains.png`
  (was `umbra-ink-mountains.png`). `install.sh config` renames the live
  copies in `~/Pictures/wallpapers` and fixes the `~/.cache/wal/wal` pointer
  instead of duplicating files.
- `home/.local/bin/aroli-cursor-size` is the canonical cursor-size tool;
  `umbra-cursor-size` remains as a deprecated shim. `build-umbra-hyprcursor.sh`
  is now `build-aroli-hyprcursor.sh` (with fallback to the old `Umbra` source
  dir when present).

## [2.3.2] - 2026-09-21

### Fixed

- Shell tool PATHs no longer vanish on every rice install (`bun: command
  not found`). The deployed rcs only carried `~/.local/bin`, and being
  symlinks into the repo, any `export PATH` an installer (bun, rustup,
  mise) appended to them was lost on the next `install.sh config` /
  `rice update`. `.zshrc`, `.profile`, `.bashrc` and `.zprofile` now build
  the standard tool PATHs themselves (`~/.bun/bin` honoring `$BUN_INSTALL`,
  `~/go/bin`, `~/.cargo/bin`, mise shims last like upstream Omarchy),
  `zsh`/`bash` activate mise when present, and personal overrides live in
  `~/.zshrc.local` / `~/.bashrc.local` / `~/.zprofile.local` /
  `~/.profile.local`, seeded once by the installer, never touched by
  updates. `install.sh` also rescues surviving tool-PATH lines from a
  replaced rc file into its `*.local` counterpart instead of dropping them.

## [2.3.1] - 2026-09-20

### Fixed

- `rice update` hands over to the target release's own updater (extracted
  from git without touching the worktree), so updater fixes apply to the
  very update that delivers them, an old backend can no longer block its
  own replacement. Tags without a usable updater are refused, and an
  unreachable remote falls back to the checkout's updater.

## [2.3.0] - 2026-09-20

### Fixed

- `rice update` and `rice rollback` no longer refuse to run over a dirty
  checkout. Files the desktop rewrites on its own (`kdeglobals`, spicetify
  colors) are reset, they rebuild on the next wallpaper/theme change, and
  anything else is stashed automatically and restored after the checkout.
  Changes that do not re-apply cleanly stay in the stash, never deleted.
- The installer reloads a live Hyprland session at the end of the `final`
  phase, so an update can no longer leave a stale config-error banner (e.g.
  "cannot open hyprland.lua" caught mid-checkout) parked on screen.

### Added

- Login update notice (`rice-update-notify.service`, enabled by default):
  reads only the update cache on graphical login, no network, oneshot,
  idle priority, and shows one desktop notification per pending release.
  Disable with `systemctl --user disable rice-update-notify.service`; the
  daily check is now also nag-once-per-release instead of once per day.
- The Omarchy mark (bar logo and Settings › About header) is tinted with
  the wallpaper accent via `MultiEffect` colorization, so it follows the
  palette like the Arch glyph already did. Untinted images stay untouched:
  the tint only applies when `BarItem.imageTint` is set.

### Changed

- Faster terminal startup: the first-kitty detection is one `hyprctl`
  call instead of two (workspace derived from `$PPID`, no `ps` fork),
  sprite transcoding moved to the background on cache miss, and
  OMZ/plugins/theme are byte-compiled (`zcompile`, refreshed only when
  sources change).

## [2.2.0] - 2026-09-20

### Added

- Quick-app keybinds with `launch-*.sh` helpers honoring system defaults
  with rice fallbacks: `Super+T` / `Super+Alt+T` (terminal), `Super+Alt+M`
  (Spotify), `Super+Alt+D` (Discord), and `Super+Alt+C` (default text
  editor).
- Update-safe Hyprland user override layer (`user.lua` / `user.conf` seeded
  once from `.example` templates and loaded last), with dirty-tree guard
  messaging and prune protection in `install.sh`.
- Settings shortcuts panel now parses keybinds from both the base and user
  override files with shared variables, plus pt-BR/es translation updates.

## [2.1.0] - 2026-09-19

### Added

- Configurable wallpaper palette intensity (0-4 slider): level 0 keeps the
  stable Umbra surfaces, level 4 reproduces the previous raw pywal output,
  and existing `paletteMode` settings migrate to the matching endpoint so
  no one's desktop changes on upgrade.
- Full configurable Umbra palette system: `Wallpaper`, `Umbra`, `Hybrid`,
  and `Manual` presets, manual canvas/surface/text/accent tokens with
  capture and restore, saturation and minimum-contrast guards, and
  independent scopes for shell, terminal, GTK/Qt, and Hyprland.
- Quickshell threshold lock screen (WlSessionLock + PAM card with morph
  badge, pill input, caps/num states, and pt-BR/es translations) replacing
  the hyprlock session island. Lock entry points (idle, power menu,
  launcher, `Super+L`) go through `loginctl`; the hyprlock binary stays
  installed for the remote/RustDesk profile.
- Consolidated adaptive desktop controls: unified hypridle normal/remote
  profiles, Umbra XCursor build with text-cursor aliases and live apply,
  and expanded network, Bluetooth, and system settings panels.

### Fixed

- Declared the runtime shell dependencies in `packages/pacman.txt`.
- Kept the Umbra threshold treatment out of expanded panels and removed
  the broken threshold notch treatment.

### Changed

- The remote profile now runs hyprlock with its built-in defaults (the
  session lock-preference sync script left with the hyprlock island).

## [2.0.0] - 2026-09-17

### Added

- Coordinated session profiles, recovery, preference export/import, and
  explicit rules between Game Mode, Reading Mode, and Power Saver.
- A reversible Game Mode toggle in the Control Panel and `rice gaming`.
- A Power Profile selector in the Control Panel, System Settings, and `rice
  battery` for power-saver, balanced, and supported performance modes.
- A read-only smoke-test release gate and a 2.0 compatibility guide.

## [1.1.7] - 2026-09-17

### Fixed

- Release asset uploads now remove stale assets before re-uploading, avoiding
  GitHub's duplicate asset name validation errors.

## [1.1.6] - 2026-09-17

### Fixed

- TUI actions now suspend Bubble Tea while running, allowing native `sudo`
  password prompts and other interactive commands to receive terminal input.
- Release publication is idempotent and no longer relies on a hanging upload
  action.

## [1.1.5] - 2026-09-17

### Fixed

- Plugin discovery now runs asynchronously with a loading state, cancellation
  path, and a 45-second network timeout instead of blocking the TUI.

## [1.1.4] - 2026-09-17

### Fixed

- Suppressed temporary repository clone output at the subprocess level so it
  cannot leak into the Bubble Tea screen.

### Changed

- Added screen titles and breadcrumbs throughout the TUI navigation flow.

## [1.1.3] - 2026-09-17

### Fixed

- Kept plugin catalog bootstrap and repository cloning out of the Bubble Tea
  screen while opening the plugin manager.

## [1.1.2] - 2026-09-17

### Fixed

- Prevented installer and diagnostic output from corrupting the Bubble Tea
  alternate screen while actions run inside the TUI.

## [1.1.1] - 2026-09-17

### Fixed

- `rice cli update` now replaces the executable actually selected by the
  current shell, including installations managed through Go or mise.

## [1.1.0] - 2026-09-17

### Added

- Go `rice` CLI with a guided terminal interface, individual optional-package management (`rice plugins list` and `rice plugins install`), verified self-updates, resumable installation checkpoints, and persistent install logs.

- Release and update system: `VERSION` as the single source of truth,
  `rice` CLI (`version`, `status`, `check`, `update`, `rollback`, `prune`),
  background check via systemd user timer (24 h, idle priority, zero
  resident processes), and update badge in the shell (reads a local cache
  file only, no extra polling).

### Changed

- Replaced the numeric first-run menu with a keyboard-driven Bubble Tea v2
  interface using Charm styling, plugin selection, confirmations, and action
  feedback.
- Added Linux amd64 and arm64 CLI release binaries with SHA-256 verification.

## [1.0.0] - 2026-09-17

### Added

- First versioned release of Umbra Noctis. Everything before this tag is
  unversioned history; from here on every user-visible change lands under
  a `vX.Y.Z` tag with release notes.
