#!/usr/bin/env bash
# Exercise config planning against a disposable HOME and prove dry-run writes
# neither user files nor generated identity/state files.
set -euo pipefail

root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
home="$(mktemp -d "${TMPDIR:-/tmp}/aroli-installer-home.XXXXXXXX")"
trap 'rm -rf -- "$home"' EXIT

mkdir -p "$home/.config/hypr"
mkdir -p "$home/.cache/aroli-desktop"
mkdir -p "$home/.config"
printf 'user shell content\n' > "$home/.zshrc"
printf 'user hypr content\n' > "$home/.config/hypr/user.lua"
printf '{"previous_ref":"fixture-ref","current_ref":"current-ref"}\n' \
    > "$home/.cache/aroli-desktop/state.json"
printf '{"local":"0.1.0-beta.1","remote":"v0.1.0-beta.2","update_available":true,"checked_at":%s}\n' \
    "$(date +%s)" > "$home/.cache/aroli-desktop/update.json"
: > "$home/.cache/aroli-desktop/aroli.lock"
before_files="$(find "$home" -type f -o -type l | sort)"
before_hashes="$(find "$home" -type f -print0 | sort -z | xargs -0 -r sha256sum)"

HOME="$home" XDG_CACHE_HOME="$home/.cache" XDG_CONFIG_HOME="$home/.config" \
    USER=aroli CI=1 bash "$root/install.sh" --dry-run --yes --lang en config >/dev/null

if [[ -d "$root/.git" ]]; then
    rollback_output="$(HOME="$home" XDG_CACHE_HOME="$home/.cache" \
        XDG_CONFIG_HOME="$home/.config" AROLI_DESKTOP_REPO="$root" \
        bash "$root/home/.local/bin/aroli-backend" --dry-run rollback)"
    grep -q 'would: git checkout fixture-ref' <<<"$rollback_output" || {
        echo 'rollback dry-run did not report the saved rollback target' >&2
        exit 1
    }

    update_output="$(HOME="$home" XDG_CACHE_HOME="$home/.cache" \
        XDG_CONFIG_HOME="$home/.config" AROLI_DESKTOP_REPO="$root" \
        bash "$root/home/.local/bin/aroli-backend" --dry-run update)"
    grep -q "would: hand over to v0.1.0-beta.2's own updater" <<<"$update_output" || {
        echo 'update dry-run did not report the beta-channel plan' >&2
        exit 1
    }
else
    printf 'Skipping repo-history dry-runs: source archive has no .git directory.\n'
fi

after_files="$(find "$home" -type f -o -type l | sort)"
after_hashes="$(find "$home" -type f -print0 | sort -z | xargs -0 -r sha256sum)"
if [[ "$after_files" != "$before_files" || "$after_hashes" != "$before_hashes" ]]; then
    printf 'dry-run changed files in its disposable HOME\nBefore:\n%s\nAfter:\n%s\n' \
        "$before_files" "$after_files" >&2
    exit 1
fi

if HOME="$home" XDG_CACHE_HOME="$home/.cache" XDG_CONFIG_HOME="$home/.config" \
    USER=aroli CI=1 bash "$root/install.sh" --dry-run --yes --lang invalid config >/dev/null 2>&1; then
    echo 'installer accepted an invalid language' >&2
    exit 1
fi

printf 'installer dry-run isolation passed\n'
