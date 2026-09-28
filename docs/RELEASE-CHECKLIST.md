# Beta release checklist

Use this checklist for each `v0.1.0-beta.N` pre-release.

## Before tagging

- [ ] Update `VERSION` and add the matching `CHANGELOG.md` entry.
- [ ] Run CI, `scripts/smoke-test.sh`, and `scripts/test-installer-dry-run.sh`.
- [ ] Complete the disposable Arch VM checklist in [Support](SUPPORT.md); record
  any unavailable Omarchy session check explicitly.
- [ ] Verify the release commit is a fast-forward of `main` and that no user
  changes or secrets are in the release tree.
- [ ] Create and push the signed `v0.1.0-beta.N` tag only after review.

## After CI publishes assets

- [ ] Confirm the release workflow passed tag/version, CLI, manifest, and
  artifact checks.
- [ ] Confirm the source archive, `aroli-linux-amd64`,
  `aroli-linux-arm64`, and `SHA256SUMS.txt` exist.
- [ ] Recompute SHA-256 for downloaded assets and compare with the published
  checksum file.
- [ ] Confirm GitHub marks the release as a pre-release and does not promote
  it as a stable release.
- [ ] Link the sanitized VM smoke record from the release PR.
