# Aroli Desktop on Omarchy

Aroli Desktop is **Omarchy-friendly**: an Aroli layer for Arch, Hyprland, and
Quickshell that runs on top of the stack already provided by Omarchy.

## Compatibility rules

- Never modify `/usr/share/omarchy/`. The directory belongs to the package and
  will be replaced on updates.
- Keep customizations in `~/.config`, `~/.local/share`, and this clone.
- Prefer the `omarchy` commands for system operations: `omarchy update`,
  `omarchy pkg`, and `omarchy theme`.
- Aroli Desktop starts its Quickshell shell after `WAYLAND_DISPLAY` is
  available. A brief flash of the native bar may appear during login.

## What Aroli Desktop changes

| Area | Behavior |
| --- | --- |
| Shell | Aroli Desktop Quickshell notch and panels. |
| Palette | Pywal extracts wallpaper accents over stable Aroli Dark surfaces; Settings › Appearance can restore wallpaper-tinted surfaces. |
| Language | Shell, Hyprlock, and SDDM in pt-BR. |
| Cursor | Aroli Pointer (XCursor), without setting `HYPRCURSOR_THEME`. |
| Keyboard | ABNT2 (`br`) in Hyprland. |
| Optional packages | Each item requires individual confirmation. |
| Adaptive modes | Session, Game, Reading, and power profiles use user-level state only. |

## Preserved components

The project does not replace Hyprland, SDDM, PipeWire, WirePlumber,
NetworkManager, `xdg-desktop-portal-hyprland`, or the `omarchy` package.

## Safe diagnostics

```sh
omarchy debug --no-sudo --print
hyprctl configerrors
./diagnose
```

For beta releases and the update flow, see [README.md](README.md#updates) and
[docs/UPDATES.md](docs/UPDATES.md).
