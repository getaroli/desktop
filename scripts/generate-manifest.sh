#!/usr/bin/env bash
# Generates the package contract shipped with a release. No clock, host, or
# package-manager state enters the result, so CI can verify the exact asset.
set -euo pipefail
root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
out="${1:-$root/manifest.json}"
cd "$root"

python3 - "$out" <<'PY'
import json, pathlib, sys
root = pathlib.Path("packages")
sets = {}
for path in sorted(root.glob("*.txt")):
    items = []
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#"): continue
        items.append(line.split("\t", 1)[0].split()[0])
    sets[path.stem] = sorted(set(items))
payload = {"schema": 1, "packages": sets}
pathlib.Path(sys.argv[1]).write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n")
PY
