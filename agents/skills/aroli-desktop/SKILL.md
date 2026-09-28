---
name: aroli-desktop
description: >
  REQUIRED for anything about Aroli Desktop (the Hyprland + Quickshell rice).
  ALWAYS load this skill first when the user mentions the rice, the desktop
  look, wallpaper theming, Spotify theming, the bar, the notch, control
  center, Settings, launcher, lock screen, install, or update, or
  Aroli in any form. Triggers: aroli, rice, quickshell, hyprland config,
  wallpaper, pywal, spicetify, spotify theme, notch, island, control center,
  Super+D, settings window, sddm, install.sh, or aroli CLI. After loading,
  follow its routing to the topic skill.
---

# Aroli Desktop

Operating contract for AI agents working on this machine's rice. Load the
topic skill it points to; that skill carries the commands.

## The three laws (non-negotiable, every task)

1. **Facts before synthesis.** Inspect files and run commands yourself before
   answering. Never guess paths, versions, or API shapes. If evidence
   contradicts an earlier claim, say so and trust the evidence.
2. **Integrity first.** Read-only commands (`diagnose`, `--dry-run`,
   `status`, log reads) come before anything that writes. Never run sudo,
   destructive, or data-deleting commands without the user's explicit go for
   that exact action. Never delete backups (`~/.dotfiles-backup`,
   `~/.audit-backups`) without being asked. Never print secrets, tokens,
   or credentials.
3. **Report upstream.** When you hit a real bug in Aroli Desktop (not user
   config), gather the evidence bundle and offer to file it as a GitHub
   issue. See `aroli-report`. Only file with the user's confirmation.

## Routing (load ONE topic skill next)

- Something looks wrong, is slow, or broke → `aroli-diagnose`
- Wallpaper, colors, or theming did not follow → `aroli-wallpaper`
- Spotify look, spicetify errors, stuck scheme → `aroli-spotify`
- Editing/adding Quickshell QML, bar, notch, panels, Settings UI → `aroli-quickshell`
- Installing, updating, migrating, rice CLI, leftover cleanup → `aroli-install`
- Filing or triaging a GitHub issue → `aroli-report`

## Communication (UX is part of the product)

- Match the user's language (default pt-BR here). Short, direct, no emojis
  unless asked. No praise fillers, no false agreement: disagree with data.
- Show the exact commands you ran and their outcome; paste short outputs,
  never paraphrase numbers.
- Before a material action, say what it does and get approval. After it,
  say what changed and how to verify.
- Never present a plan as done. `done` means executed AND verified.
