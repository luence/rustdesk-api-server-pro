# Release Process

Use the original tag-based workflows for releases:

- `.github/workflows/build-release.yml`
- `.github/workflows/ghcr-docker.yml`

Do not create one-off release workflows unless absolutely necessary.

## Standard Release Steps

1. Ensure `main` is clean and all intended changes are merged.
2. Confirm compliance files are present:
   - `LICENSE`
   - `NOTICE`
   - `THIRD_PARTY_NOTICES.md`
   - `DISCLAIMER.md`
   - `SECURITY.md`
   - `PRIVACY.md`
   - `COMPLIANCE.md`
   - `CONTRIBUTING.md`
   - `CODE_OF_CONDUCT.md`
   - `SUPPORT.md`
3. Confirm no secrets or runtime data are tracked:
   - `server.db`
   - `*.sqlite`
   - production `server.yaml`
   - OAuth secrets
   - SMTP passwords
   - private keys
   - logs
   - `record_uploads/`
4. Push a semantic tag matching `v*.*.*`.

Example:

```bash
git checkout main
git pull origin main
git tag v1.4.8
git push origin v1.4.8
```

## Expected Workflows

Tag push should trigger:

- `build-release.yml`: builds Linux, Windows and macOS packages for amd64 and arm64, then uploads zip assets to GitHub Release.
- `ghcr-docker.yml`: builds and pushes GHCR image tags, including `latest`, version tag and `sha-*`.

## Do Not Include in Release Assets

- production databases;
- production config files;
- OAuth client secrets;
- SMTP credentials;
- user recordings;
- private logs;
- private keys or certificates.

## If a Release Fails

Prefer fixing the original workflow instead of adding a new one-off workflow.

If a tag or release is partially created:

1. Delete the failed GitHub Release if present.
2. Delete the failed tag both remotely and locally.
3. Fix the workflow or source issue.
4. Re-create the tag.

Example:

```bash
git tag -d v1.4.8 || true
git push origin :refs/tags/v1.4.8 || true
git tag v1.4.8
git push origin v1.4.8
```
