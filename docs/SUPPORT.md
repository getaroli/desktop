# Support and compatibility

This page defines what the project currently supports and what its automated
checks prove. A passing syntax or container check does not mean a full desktop
session was exercised.

## Compatibility policy

| Environment | Policy | What CI checks |
| --- | --- | --- |
| Arch Linux, x86_64 | Primary desktop target; rolling release, current packages expected | Latest Arch container checks installer parsing and a disposable-HOME dry-run; it does not install a full desktop |
| Omarchy on Arch, x86_64 | Supported integration; use Omarchy package commands and preserve `/usr/share/omarchy` | Linux CI checks delivery contracts; an Omarchy graphical session is not provisioned in CI |
| Hyprland | 0.56 or newer is required | No live compositor test in CI |
| Quickshell / Qt | Use the versions supplied by the supported Arch or Omarchy system | Ubuntu CI runs `qmllint` for syntax; its Qt 6.4 environment cannot validate all runtime imports |
| AArch64 / ARM64 | The CLI binary is published; the complete desktop installation is not a supported target | Go build/tests cover the CLI architecture only when explicitly added to CI |
| Other distributions, compositors, or display servers | Best effort; not part of the supported desktop matrix | No compatibility claim |

Arch and Omarchy are rolling distributions, so the project does not promise a
fixed maximum supported package age. The minimum Hyprland requirement is the
version boundary currently relied on by the runtime controls; a future change
that raises it must update this page, `README.md`, and the release notes.

## Release smoke checklist

Before each beta release, run the following on a disposable x86_64 Arch VM. For
an Omarchy beta, repeat the session checks on a disposable Omarchy VM when
one is available. Do not run this checklist on a personal workstation.

1. Start from a fresh user account with no Aroli files or prior snapshots.
2. Run `./diagnose` and `./install.sh --dry-run --lang en`; retain the output
   with usernames and machine identifiers removed.
3. Review the plan, then run the install interactively and confirm only the
   packages and system changes needed for the test.
4. Log out and back in. Check `aroli status`, `aroli diagnose`,
   `hyprctl configerrors`, the Quickshell service, wallpaper palette updates,
   and the documented language/keyboard behavior.
5. Create a snapshot, change a supported config, and restore it. Confirm the
   previous config is retained as a `.before-snapshot-*` backup.
6. On a disposable clone, exercise `aroli update --dry-run`, a beta update,
   and rollback. Confirm a failed or interrupted update remains recoverable.
7. Record the OS image date, Omarchy version if applicable, Aroli tag, hardware
   architecture, pass/fail result, and sanitized logs in the release PR.

This manual VM check complements CI; CI remains the merge gate for source,
delivery, shell, QML syntax, and Go tests.

## Bug reports and security reports

Use the [bug report template](https://github.com/getaroli/desktop/issues/new?template=bug_report.yml)
for reproducible product bugs. See the repository [security policy](../SECURITY.md)
to report vulnerabilities privately.
