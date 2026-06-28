# Legal and Open-Source Compliance Guide

This guide records project-level compliance expectations for maintainers, redistributors, and operators.

It is not legal advice. For commercial redistribution, hosted services, enterprise deployment, or product integration, consult qualified legal counsel.

## 1. License Model

This repository is licensed under AGPL-3.0-only.

Key practical meaning:

- You may use, study, modify, and redistribute the software under AGPL-3.0.
- If you modify the software and let users interact with it over a network, you must offer those users the corresponding source code of the modified version.
- Binary releases and container images should preserve license and notice files.
- Do not combine this project with code that cannot be distributed under AGPL-compatible terms unless it is clearly separate and legally compatible.

Core files:

- `LICENSE`
- `NOTICE`
- `THIRD_PARTY_NOTICES.md`
- `DISCLAIMER.md`
- `COMPLIANCE.md`

## 2. Third-Party Attribution

The frontend is based on Soybean Admin, which identifies itself as MIT-licensed in `soybean-admin/package.json`.

When redistributing built frontend assets or release packages, preserve attribution and notices. Release artifacts should include at least:

- `LICENSE`
- `NOTICE`
- `THIRD_PARTY_NOTICES.md`
- `DISCLAIMER.md`
- `SECURITY.md`
- `PRIVACY.md`
- `COMPLIANCE.md`

## 3. Trademark and Non-Affiliation Rules

Do not describe this project as official RustDesk, Huawei, HarmonyOS, OpenHarmony, GitHub, Google, Docker, GHCR, or OpenWrt software.

Allowed wording:

- independent compatibility implementation;
- compatible API workflow;
- deployment guide for OpenWrt Docker;
- OAuth provider integration example.

Avoid wording such as:

- official RustDesk server;
- Huawei-certified;
- HarmonyOS official SDK;
- endorsed by RustDesk;
- approved by OpenHarmony.

## 4. Release Artifact Rules

Release artifacts must not contain:

- production `server.yaml` files;
- SQLite databases such as `server.db`;
- OAuth client secrets;
- SMTP credentials;
- tokens or private keys;
- reverse proxy private certificates;
- uploaded recordings;
- logs containing personal data;
- private screenshots or issue attachments.

Release artifacts should contain license and notice files so that downloaded binaries remain compliant even outside the source repository.

## 5. Security and Privacy Rules

Operators may process sensitive data such as account records, device metadata, audit logs, OAuth identity claims, uploaded recordings, IP addresses and user-agent data.

Operators are responsible for:

- publishing their own privacy policy where required;
- securing `/app/data` and backups;
- using HTTPS for public access;
- rotating default credentials;
- protecting OAuth and SMTP secrets;
- complying with local laws and sector-specific requirements.

## 6. Dependency Review

Before major public releases:

```bash
cd backend
go list -m all
```

```bash
cd soybean-admin
pnpm install --frozen-lockfile
pnpm licenses list || true
```

Recommended extra checks:

- generate SBOM for Docker images;
- scan container images for known vulnerabilities;
- review dependency license changes;
- review GitHub Actions permissions.

## 7. Maintainer Checklist

Before merging release-related changes:

- [ ] License and notice files are present.
- [ ] Release workflow includes notice files in artifacts.
- [ ] No temporary release workflow is active.
- [ ] `.gitignore` protects runtime data and secrets.
- [ ] CI checks required compliance files.
- [ ] Documentation does not imply official third-party endorsement.
- [ ] OpenWrt scripts back up data before replacing containers.
- [ ] Docker commands keep host networking and port labels where documented.
