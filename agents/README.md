# Aroli Desktop agent skills

Autonomous skills for AI agents working on this machine, Omarchy-style:
shipped with the OS, discovered without being mentioned, grounded in
evidence, protective of system integrity.

- Source of truth: `agents/skills/<name>/SKILL.md` in this repo.
- Installed (OS level): `~/.agents/skills/<name>`, symlinked in `--link`
  mode so skill updates arrive with the rice, copied in `--copy` mode.
  Deployed by `install.sh` step 4c. Discovery also works for
  Claude-compatible agents from the same path.
- `aroli-desktop` is the hub: broad triggers route to one topic skill.
  Topic skills carry the commands; the hub carries the laws (facts
  first, integrity first, report upstream) and the communication bar.

| Skill | When |
| --- | --- |
| `aroli-desktop` | Anything Aroli, load first, then follow routing |
| `aroli-diagnose` | Something broke or looks wrong |
| `aroli-wallpaper` | Wallpaper/palette did not follow |
| `aroli-spotify` | Spotify theming, spicetify errors |
| `aroli-quickshell` | Editing QML, panels, translations |
| `aroli-install` | Install, update, cleanup |
| `aroli-report` | Filing a GitHub issue with evidence |
