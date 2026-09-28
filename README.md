<h1 align="center">Aroli Desktop</h1>

<p align="center">
  A premium visual layer for <strong>Omarchy</strong>, Arch Linux, Hyprland, and Quickshell.<br>
  Dark by nature. Personal by wallpaper. Yours by choice.
</p>

<p align="center">
  <a href="https://github.com/getaroli/desktop/stargazers"><img src="https://img.shields.io/github/stars/getaroli/desktop?style=flat-square&color=8b7cff&label=stars" alt="GitHub stars"></a>
  <a href="https://github.com/getaroli/desktop/releases"><img src="https://img.shields.io/github/v/release/getaroli/desktop?display_name=tag&style=flat-square&color=8b7cff&label=release" alt="Latest release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/getaroli/desktop?style=flat-square&color=8b7cff" alt="GPL-3.0 license"></a>
  <a href="https://github.com/getaroli/desktop/actions"><img src="https://img.shields.io/github/actions/workflow/status/getaroli/desktop/release.yml?style=flat-square&label=build" alt="Build status"></a>
</p>

<p align="center">
  <a href="#installation">Install</a> ·
  <a href="#aroli-desktop">Aroli Desktop</a> ·
  <a href="#privacy-and-control">Privacy</a> ·
  <a href="OMARCHY.md">Omarchy</a> ·
  <a href="LLMS.md">LLMS</a>
</p>

---

## Aroli Desktop

**Aroli Desktop** is the ready-to-use Aroli workspace on top of the
stable Omarchy base: the desktop product of the Aroli family
(Aroli Themes, Aroli Pointer, Aroli Backdrops). Aroli Desktop keeps the
structure, notch, dark surfaces, cutouts, and luminance bands,
while each wallpaper sets the mood with Pywal.

The dark foundation is fixed: depth, negative space, and contrast stay
Aroli Dark. Your wallpaper supplies the accent palette, so the result changes with
your image without turning every application surface into its dominant colour.

| Base | Aroli surface | Your choice |
| --- | --- | --- |
| Omarchy · Arch · Hyprland | Quickshell, notch, and dynamic palette | Wallpaper, language, and optionals |

<p align="center">
  <a href="https://github.com/getaroli/desktop/commits/main"><img src="https://img.shields.io/github/commit-activity/m/getaroli/desktop?style=flat-square&color=8b7cff&label=community%20activity" alt="Monthly commit activity"></a>
  <a href="https://github.com/getaroli/desktop/issues"><img src="https://img.shields.io/github/issues/getaroli/desktop?style=flat-square&color=8b7cff&label=issues" alt="Open issues"></a>
</p>

<p align="center">
  <img src="home/Pictures/wallpapers/aroli-ember-coast.png" alt="Ember Coast, bundled wallpaper" width="49%">
  <img src="home/Pictures/wallpapers/aroli-obsidian-dunes.png" alt="Obsidian Dunes, bundled wallpaper" width="49%">
</p>

## What ships with Aroli Desktop

- Complete Quickshell interface, with notch, launcher, panels, and overview.
- Interface translation in **Brazilian Portuguese** and ABNT2 keyboard layout.
- **Aroli Pointer** cursor and visual identification of the active platform: Omarchy on
  Omarchy, Arch on Arch.
- Pywal accents extracted from the current wallpaper, applied over stable dark
  Aroli Dark surfaces. Settings › Appearance › Colour system can enable
  wallpaper-tinted application backgrounds.
- Four Aroli Backdrops installed with the aroli: Ember Coast, Silent
  Threshold, Obsidian Dunes, and Black Mountains.
- Aroli New Tab for Brave Origin: search, clock, shortcuts and tasks on the
  wallpaper palette, with a Sistema / Aroli Dark / Aroli Black switch that
  also repaints the browser frame. No network, no polling.

Change the cursor size consistently across Hyprland and GTK with
`aroli-cursor-size 40`. Values from 16 to 96 are accepted; relog afterwards
so already-running applications reload the XCursor.

<details>
<summary><strong>See the aroli in motion</strong></summary>
<br>

[Open video demo](https://github.com/user-attachments/assets/2b35a6fb-5a08-4539-99a9-7c525eb463b3)

</details>

## Installation

> [!IMPORTANT]
> Read the plan before writing to the system. Dry-run mode changes no files,
> installs no packages, and asks for no privileges.

> [!TIP]
> `aroli` is the recommended installation path. It gives first-time users a guided
> terminal assistant, while the explicit subcommands remain suitable for automation.

```sh
# Beta CLI binary from GitHub Releases
arch="$(uname -m)"; case "$arch" in x86_64) arch=amd64;; aarch64|arm64) arch=arm64;; esac
curl -fL "https://github.com/getaroli/desktop/releases/download/v0.1.0-beta.1/aroli-linux-$arch" -o ~/.local/bin/aroli
chmod +x ~/.local/bin/aroli

# Or install from the current main branch (Go 1.27+, development version)
go install github.com/getaroli/desktop/cmd/aroli@main

# Opens a guided terminal interface for first-time users
aroli

# Or use the non-interactive commands
aroli install --dry-run --lang pt-BR
aroli install --lang pt-BR
```

Requires **Hyprland 0.56+**. Aroli Desktop uses `hyprland.lua`, not
`hyprland.conf`.

| Command | Result |
| --- | --- |
| `aroli` | Opens the guided terminal interface. |
| `aroli install --dry-run --lang pt-BR` | Shows the plan without changing anything. |
| `aroli install --lang pt-BR` | Installs the aroli in Brazilian Portuguese. |
| `aroli install restore` | Restores the previous configuration files. |
| `aroli diagnose` | Read-only diagnostics. |
| `aroli doctor` | Checks the local desktop integration; `--fix` offers safe session repairs. |
| `aroli status` | Installed version, last check, and rollback ref. |
| `aroli profile list` | Lists curated, transparent installation profiles. |
| `aroli profile install creator` | Installs a named set of phases; its optional apps remain visible and confirmable. |
| `aroli snapshot create before-tweaks` | Saves the supported visual configuration paths locally. |
| `aroli snapshot restore before-tweaks` | Restores a configuration snapshot while retaining the current files as backups. |
| `aroli wallpaper list` | Lists bundled and imported wallpapers. |
| `aroli wallpaper set aroli-ember-coast.png` | Applies a wallpaper and its dynamic Pywal palette. |
| `aroli gaming on` | Temporarily disables expensive compositor effects, preserving their exact previous state. |
| `aroli gaming launch steam` | Runs a game under GameMode when the optional `gamemode` package is installed. |
| `aroli battery set power-saver` | Selects the Power Profiles Daemon energy-saver profile. |
| `aroli battery balanced` | Selects the balanced profile; `performance` is available when supported by the hardware. |
| `aroli session apply laptop` | Applies a coordinated daily-use profile. |
| `aroli recover --snapshot latest` | Recovers stuck modes and can restore the newest local snapshot. |
| `aroli export ~/aroli-preferences` | Exports portable visual preferences without credentials or personal data. |
| `aroli update --dry-run` | Shows the update plan without changing anything. |
| `aroli prune` | Lists files retired by the latest release (nothing is deleted without `--apply`). |
| `aroli plugins list` | Lists optional applications and tools that can be installed later. |
| `aroli plugins install btop` | Installs only the selected optional item, after confirmation. |
| `aroli cli update --dry-run` | Checks the CLI release and its verified binary update plan. |
| `aroli cli update` | Updates only the CLI binary after verifying SHA-256. |
| `aroli install --resume` | Continues from the last successfully completed installation phase. |
| `aroli logs` | Shows the last 100 lines of the latest installation log. |

See [supported environments and release smoke checks](docs/SUPPORT.md) and the
[beta release checklist](docs/RELEASE-CHECKLIST.md). Report security issues
privately using [SECURITY.md](SECURITY.md).

## Updates

Aroli Desktop is in beta and follows **beta tags** (`v0.1.0-beta.N`), never
`main`. Beta releases are pre-releases and updates require your confirmation:

- A daily timer (`aroli-update-check.timer`, idle priority, no
  resident process) transfers a few KB and writes
  `~/.cache/aroli-desktop/update.json`. The bar/notch shows a dot
  and `Settings > About` shows the version seen.
- A login notice (`aroli-update-notify.service`) reads only that cache file,
  no network, exits in milliseconds, and shows one desktop notification per
  release when an update is pending.
- Updating is always your own act: `aroli update --dry-run` shows the plan,
  `aroli update` asks for confirmation, hands over to the target release's
  own updater, records the rollback point, and re-runs
  the installer. Local changes are stashed automatically and restored
  afterwards, so files the desktop rewrote on its own never block the update.
  `aroli rollback` undoes it. `aroli prune --apply` moves
  leftovers to the Trash (with backup), never deletes directly, and never
  touches files you modified.
- To turn off the notices (the cache, the dot, and `aroli status` keep working):
  `systemctl --user disable aroli-update-check.timer aroli-update-notify.service`.
Details and history in [CHANGELOG.md](CHANGELOG.md).

## Privacy and control

Aroli Desktop does not decide your desktop for you.

- Extra apps never come by default: each item in
  `packages/optional-*.txt` requires individual confirmation.
- Runs with `--yes`, CI, or no terminal **do not install unselected
  optionals**; named items and documented profiles remain explicit choices.
- The project adds no telemetry, analytics, or remote services.
- The installer preserves backups of replaced configurations.
- Snapshots are stored locally in `~/.local/share/aroli-desktop/snapshots` and
  never include credentials, browser data, or personal documents.
- The project does not touch `/boot`, the bootloader, or partitions.

For automations and AI agents, [LLMS.md](LLMS.md) is the operating contract:
audit before changing, ask for consent for material actions, and never
collect or expose credentials, tokens, history, or personal profiles.

## Omarchy-friendly

Aroli Desktop works on top of Omarchy without replacing its foundation. It does
not edit `/usr/share/omarchy/` and keeps customizations in the appropriate
user paths. See [OMARCHY.md](OMARCHY.md) for compatibility, limits, and
safe diagnostics.

---

<p align="center">
  <sub>
    Aroli Desktop · the Aroli workspace on open platforms<br>
    Code under <a href="LICENSE">GPL-3.0</a> · details in <a href="docs/IDENTITY.md">Identity</a>
  </sub>
</p>
