# Aroli Desktop Identity

Aroli Desktop is the desktop product of the Aroli family: the ready-to-use
workspace (Hyprland + Quickshell + installer) built on Aroli Themes,
Aroli Pointer, and Aroli Backdrops. It follows the Aroli direction,
“Tudo encontra seu lugar”, without adding a product logo of its own.

## Concept

The desktop holds the structure while the wallpaper sets the mood: the
system provides the dark surface and shape (Aroli Dark); the wallpaper
provides the accent colour. This keeps the workspace recognisably
Aroli even with a bright or strongly coloured image.

## Visual assets

1. Dark surfaces and negative space to hold the content.
2. Luminance bands and notch cutouts as recurring structure.
3. Accent palette derived from the wallpaper over a stable dark foundation.
4. Correct branding for the current platform: Omarchy on Omarchy, Arch on
   Arch, and a neutral symbol on other distributions.

## Name and attribution

- Full name: **Aroli Desktop**.
- Short form: **Desktop**, when the Aroli context is clear.
- Aroli product names in use: Aroli Themes (Dark / Black), Aroli Pointer,
  Aroli Backdrops. The rice adds no new brand marks.
- Omarchy and Arch are trademarks of their respective communities; they are
  shown only to identify the active platform.
- The code remains under [GPL-3.0-only](../LICENSE), with required source
  attribution.

## Lock screen

A single floating card, Caelestia-placed and Aroli-dressed: the morphing
badge (lock icon becoming the face/initial) and the PAM states are
Caelestia's ideas; the stable dark surface, the wallpaper as blurred mood
only, the accent tab, and the motion tokens are
Aroli's. Every element justifies itself:

| Element | Why it exists |
| --- | --- |
| Avatar | Answers "whose password?" before asking. Click focuses the field. |
| Time | Reassurance, not hero: tabular numerals, no layout shift. |
| Date | Rice language, never the process locale. |
| Pill | The only raised surface. Accent at rest, ok while checking, crit on failure, state is colour and text, never colour alone. |
| Status | Attempt count and PAM errors. Fixed height so nothing reflows. |
| Caps/num | Warn tone: a wrong-password mystery is usually this. |
| User | Identity anchor in mono, like the bar. Click focuses the field. |
| Battery | Charge anchor. Click toggles percent and time estimate. |

The Quickshell session lock (`Lock.qml`, `LockStage.qml`, `LockCenter.qml`,
`LockPam.qml`) is the implementation. The `hyprlock` session island, its
scripts, and its language files are gone; the `hyprlock` binary stays in
`packages/pacman.txt` for the remote/RustDesk profile only
(`hypridle.d/remote.conf`). `hypridle.conf` is generated and
gitignored, the shipped profile lives in `hypridle.d/normal.conf`.
