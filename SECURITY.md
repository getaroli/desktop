# Security policy

## Supported versions

Security fixes are made against the current beta release and the current
development branch. Older tags are not maintained; upgrade to the latest beta
release before reporting a vulnerability that may already be fixed.

## Reporting a vulnerability

Please report vulnerabilities privately through
[GitHub Security Advisories](https://github.com/getaroli/desktop/security/advisories/new).
Include the affected release or commit, the smallest reproducible case, and
the impact. Do not open a public issue for an unpatched vulnerability. The
maintainer will acknowledge the report and coordinate a fix and disclosure
timeline through the advisory.

Do not include passwords, access tokens, private keys, browser profiles, or
unredacted diagnostic bundles in a report.

## Update integrity

The CLI updater accepts only beta `v0.1.0-beta.N` release tags, downloads assets
from the fixed `github.com/getaroli/desktop` release URL, verifies the binary
against its SHA-256 entry in `SHA256SUMS.txt`, writes a temporary file beside
the installed binary with mode `0755`, then atomically renames it into place.
If a download or checksum check fails, the installed binary is left untouched.

The checksum detects corruption and mismatched assets; because the checksum
file and binary come from the same GitHub release, this check alone does not
protect against a compromised release account or release pipeline. Do not
describe it as a cryptographic signature or independent publisher
authentication.
