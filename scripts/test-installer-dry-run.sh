#!/usr/bin/env bash
# Exercise config planning against a disposable HOME and prove dry-run writes
# neither user files nor generated identity/state files.
set -euo pipefail

root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
home="$(mktemp -d "${TMPDIR:-/tmp}/aroli-installer-home.XXXXXXXX")"
trap 'rm -rf -- "$home"' EXIT

mkdir -p "$home/.config/hypr"
mkdir -p "$home/.cache/aroli-desktop"
legacy="$home/.local/share/diegoMalagrida-dotfiles"
mkdir -p "$legacy" "$home/.config"
cp "$root/README.md" "$legacy/README.md"
ln -s "$legacy/README.md" "$home/.config/legacy-example"
printf 'user shell content\n' > "$home/.zshrc"
printf 'user hypr content\n' > "$home/.config/hypr/user.lua"
printf '{"previous_ref":"fixture-ref","current_ref":"current-ref"}\n' \
    > "$home/.cache/aroli-desktop/state.json"
before_files="$(find "$home" -type f -o -type l | sort)"
before_hashes="$(find "$home" -type f -print0 | sort -z | xargs -0 -r sha256sum)"

HOME="$home" XDG_CACHE_HOME="$home/.cache" XDG_CONFIG_HOME="$home/.config" \
    USER=aroli CI=1 bash "$root/install.sh" --dry-run --yes --lang en config >/dev/null

rollback_output="$(HOME="$home" XDG_CACHE_HOME="$home/.cache" \
    XDG_CONFIG_HOME="$home/.config" AROLI_DESKTOP_REPO="$root" \
    bash "$root/home/.local/bin/aroli-backend" --dry-run rollback)"
grep -q 'would: git checkout fixture-ref' <<<"$rollback_output" || {
    echo 'rollback dry-run did not report the saved rollback target' >&2
    exit 1
}

update_output="$(HOME="$home" XDG_CACHE_HOME="$home/.cache" \
    XDG_CONFIG_HOME="$home/.config" AROLI_DESKTOP_REPO="$root" \
    bash "$root/home/.local/bin/aroli-backend" --dry-run \
        --channel=unstable-dev update)"
grep -q 'would: git fetch origin unstable-dev' <<<"$update_output" || {
    echo 'update dry-run did not report the development-channel plan' >&2
    exit 1
}

migration_output="$(HOME="$home" bash "$root/scripts/migrate-from-legacy-rice.sh" \
    --source "$legacy" --target "$home/target-checkout")"
grep -q 'Would migrate:' <<<"$migration_output" || {
    echo 'migration inspection did not report the legacy symlink' >&2
    exit 1
}
[[ "$(readlink "$home/.config/legacy-example")" == "$legacy/README.md" ]] || {
    echo 'migration inspection changed the legacy symlink' >&2
    exit 1
}
[[ ! -e "$home/target-checkout" ]] || {
    echo 'migration inspection created the target checkout' >&2
    exit 1
}

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
