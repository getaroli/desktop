# Contributing to Aroli Desktop

Thanks for helping out. This guide keeps the repo consistent: one language
for collaboration, one commit style, and a safe workflow for a project that
touches people's live desktops.

- [1. Language standard](#1-language-standard)
- [2. Workflow](#2-workflow)
- [3. Commit messages](#3-commit-messages)
- [4. Branches and pull requests](#4-branches-and-pull-requests)
- [5. Testing your changes](#5-testing-your-changes)
- [6. Translations (i18n)](#6-translations-i18n)
- [7. Repository layout](#7-repository-layout)

## 1. Language standard

**English (en-US) is the standard for everything collaborative:**

- Commit messages, branch names, and PR titles/descriptions.
- Code comments, script headers, and user-facing docs (`*.md`).
- Issue reports and review discussions.

Use en-US spelling: `color` (not `colour`), `organize` (not `organise`),
`behavior` (not `behaviour`).

> Why: the desktop itself ships in three UI languages (`es`, `en`, `pt-BR`),
> but the project history must stay searchable in one language. A `git log`
> mixed across languages is hard to grep, review, and automate.

The desktop UI strings are the exception: their en-US source lives in the QML
files, with Spanish and Brazilian Portuguese dictionaries (see
[Translations](#6-translations-i18n)).

## 2. Workflow

1. Build or install the CLI, then preview before touching anything:
   ```sh
   go build -o /tmp/aroli ./cmd/aroli
   /tmp/aroli install --dry-run --lang pt-BR
   ```
2. Create a worktree + branch for each change (keeps `main` clean):
   ```sh
   git worktree add /tmp/opencode/<topic> -b fix/<topic> origin/main
   ```
3. Edit, validate (see [Testing](#5-testing-your-changes)), commit, push.
4. Open a PR against `main` and wait for review.

Never commit secrets, tokens, or personal paths. Home-directory paths in
shipped configs use `$HOME` or the `eduardoaugustolb` placeholder the
installer rewrites, never another machine's real paths.

## 3. Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short summary in English>
```

- **Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`,
  `chore`, `i18n`.
- **Scope:** the area touched (`theme`, `shell`, `installer`, `docs`...).
  Omit it when the change is repo-wide.
- **Summary:** imperative mood, lowercase, no trailing period,
  max ~72 characters. **Always in English.**
- **Body (optional):** explain *why*, in English. Wrap at ~72 columns.

Good:

```
fix(theme): cap accent saturation to stop neon Spotify

docs(installer): document --copy mode in README
```

Bad:

```
fix(tema): frena la saturación neón        # Spanish, not English
Fix stuff                                   # vague, capitalized
```

## 4. Branches and pull requests

- Branch names: `<type>/<short-topic-in-english>`, e.g.
  `fix/palette-hue-fidelity`, `docs/contributing-guide`.
- One logical change per PR. Rebase onto `origin/main` before opening.
- PR title follows the commit style; the description states what changed,
  how it was validated, and any manual follow-up (e.g. "re-copy files to
  `~/.config` on machines installed with `--copy`").
- CI/review checklist: scope is clear, no secrets, installer still passes
  `--dry-run`, derived themes regenerate on wallpaper change.

## 5. Testing your changes

- Shell scripts: `bash -n <script>` must pass.
- Python helpers: `python3 -m py_compile <script>` must pass.
- Palette/i18n changes: verify placeholder parity (`{0}`, `{1}`, `{2}`)
  between dictionaries and keep key sets identical.
- Go CLI changes: run `go test -race ./...`, `go vet ./...`, and `gofmt -l cmd/aroli`.
- Installer changes: run `/tmp/aroli install --dry-run` (or `./install.sh --dry-run` when testing the backend) and read the plan.
- Runtime changes: switch wallpapers once (`Super+Shift+W`) and confirm
  the bar, terminal, and derived app themes (btop, Discord, Spotify)
  follow the new palette.
- Read-only diagnostics live in `./diagnose`; never add write behavior
  to it.

## 6. Translations (i18n)

- Source strings stay in en-US inside the `.qml` files.
- Spanish lives in `home/.config/quickshell/translations-es.js`, Brazilian
  Portuguese in `home/.config/quickshell/translations-pt-BR.js`.
- Keys must match 1:1 across dictionaries; never translate placeholders,
  commands (`rfkill unblock bluetooth`), or proper nouns
  (pywal, Hyprland, Quickshell, PipeWire).
- Keep locale conventions: `notch`/`ilha`, `papel de parede`,
  `área de trabalho`, `janela`. When in doubt, mirror the en-US entry's
  meaning, not its words.

## 7. Repository layout

| Path | What lives there |
| --- | --- |
| `home/` | Deployed configs, mirrored into `$HOME` by `install.sh` |
| `home/.local/bin/aroli` | Version/update/rollback/prune CLI |
| `home/.config/systemd/user/aroli-update-check.{service,timer}` | Daily release check (idle, oneshot) |
| `home/.config/systemd/user/aroli-update-notify.service` | Login update notice, cache-read only (idle, oneshot) |
| `VERSION` / `CHANGELOG.md` | Single source of truth for the version + release history |
| `install.sh` | Installer (symlink mode by default, `--copy` available) |
| `diagnose` | Read-only system diagnostics |
| `packages/` | Package lists (`pacman`, `aur`, `optional-*`) |
| `docs/` | Guides: `STRUCTURE.md` (repo map), `DEVELOPMENT.md` (workflow), `UPDATES.md` (delivery), support and identity notes |
| `.githooks/` + `scripts/verify-delivery.sh` | Commit/push gates (activate: `git config core.hooksPath .githooks`) and the delivery-contract check |
| `LLMS.md` | Operating contract for automations and AI agents |

## 8. Releases

- Bump `VERSION`, add a `CHANGELOG.md` entry, commit, then tag:
  `git tag -s v0.1.0-beta.N -m "Aroli Desktop v0.1.0-beta.N"` and push the tag.
  CI checks that the tag matches `VERSION` and that `aroli --help`,
  `diagnose` and the i18n dictionaries still pass.
- GitHub Releases are pre-releases cut from the tag with the matching
  `CHANGELOG.md` entry. Clients only follow `v0.1.0-beta.N` tags
  (`aroli update` refuses anything else); `main` is never auto-applied.
