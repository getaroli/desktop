#!/usr/bin/env bash
# verify-delivery.sh: checkable part of the delivery contract (docs/UPDATES.md).
# A user-facing file must be delivered by a path a user runs. Fails on:
#   1. installer phase/function parity (ALL_PHASES/ON_DEMAND_PHASES vs phase_<name>())
#   2. orphan helper scripts (scripts/*.sh referenced by nothing)
#   3. orphan package lists (packages/*.txt consumed by nothing)
#   4. tracked backup/whitespace residue
# Read-only: changes no files. Run from anywhere inside the checkout.
set -euo pipefail

root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

fail=0
deny() { printf 'verify-delivery: %s\n' "$1" >&2; fail=1; }

# 1. Every listed phase has an implementation, every implementation is listed.
# (phase_is_done is a status helper, not a phase implementation.)
phases="$(sed -n 's/^ALL_PHASES=(\(.*\))/\1/p; s/^ON_DEMAND_PHASES=(\(.*\))/\1/p' install.sh | tr '\n' ' ')"
for p in $phases; do
  grep -qE "^phase_${p}\(\)" install.sh || deny "phase '$p' is listed but phase_${p}() is missing in install.sh"
done
while IFS= read -r fn; do
  name="${fn#phase_}"
  [ "$name" = "is_done" ] && continue
  case " $phases " in
    *" $name "*) ;;
    *) deny "phase_${name}() exists in install.sh but is in no phase list" ;;
  esac
done < <(grep -oE '^phase_[a-z_]+\(\)' install.sh | sed 's/^phase_//; s/()$//' | sort -u)

# 2-3. Every helper script and package list is referenced somewhere that runs.
refs=(install.sh diagnose Makefile cmd/aroli docs README.md OMARCHY.md LLMS.md CONTRIBUTING.md .github/workflows scripts)
for s in scripts/*.sh; do
  base="$(basename "$s")"
  [ "$base" = "verify-delivery.sh" ] && continue
  if ! grep -rl --exclude-dir=.git -- "$base" "${refs[@]}" 2>/dev/null | grep -qv "^scripts/$base$"; then
    deny "scripts/$base is referenced by nothing (installer, Makefile, CI, CLI, docs)"
  fi
done
for p in packages/*.txt; do
  base="$(basename "$p" .txt)"
  if ! grep -rl --exclude-dir=.git -- "$base" "${refs[@]}" 2>/dev/null | grep -qv "^packages/$base.txt$"; then
    deny "packages/$base.txt is consumed by nothing (installer, CLI, docs)"
  fi
done

# 4. No tracked residue.
if git ls-files | grep -qE '(\.bak$|\.orig$|\.rej$|~$)'; then
  deny "tracked backup files found (see git ls-files | grep -E '(\\.bak\$|\\.orig\$|\\.rej\$|~\$)')"
fi
git diff --check || deny "whitespace errors (see above)"

if [ "$fail" -ne 0 ]; then exit 1; fi
printf 'verify-delivery: ok\n'
