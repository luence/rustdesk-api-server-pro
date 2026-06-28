# Third-Party Notices

This document summarizes third-party components and compliance notes for this repository. It is intended as a practical compliance checklist and does not replace the upstream license texts.

## Project License

The repository is licensed under **AGPL-3.0-only**. See [`LICENSE`](./LICENSE).

AGPL-3.0 is a strong copyleft license. If you modify and run the program for users over a network, you must offer those users the corresponding source code of the modified version.

## Frontend Base

| Component | Source | License | Notes |
|---|---|---|---|
| Soybean Admin | <https://github.com/soybeanjs/soybean-admin> | MIT | The frontend package metadata identifies Soybean Admin as MIT-licensed. Preserve attribution and license notices when redistributing built frontend assets. |

## Backend Go Dependencies

Backend dependencies are declared in [`backend/go.mod`](./backend/go.mod). They include packages such as Iris, Cobra, SQLite, MySQL driver, gopsutil, yaml.v3 and other Go ecosystem modules.

Recommended release practice:

```bash
go version
cd backend
go list -m all
```

For formal releases, generate a dependency license inventory with a tool such as `go-licenses`, `licensee`, `scancode-toolkit`, or an equivalent CI process. Keep the generated report with release artifacts when distributing binaries or container images.

## Frontend npm Dependencies

Frontend dependencies are declared in [`soybean-admin/package.json`](./soybean-admin/package.json) and locked by the package lockfile used by the project.

Recommended release practice:

```bash
cd soybean-admin
pnpm licenses list
pnpm licenses list --json
```

If the package manager does not provide license output in the current environment, generate an SBOM or dependency license report with a dedicated tool before public binary/container release.

## Container Base Images

The Docker build currently uses:

- `golang:alpine`
- `node:20-alpine`
- `alpine:3.20.3`

Container images include their own upstream packages and licenses. For production releases, generate an image SBOM, for example with Syft or an equivalent scanner.

## Trademarks and Compatibility References

Names such as RustDesk, Huawei, HarmonyOS, OpenHarmony, GitHub, Google, OIDC, Docker, GHCR and OpenWrt belong to their respective owners. They are referenced only for compatibility, interoperability, deployment, or descriptive purposes.

This project is not affiliated with, endorsed by, sponsored by, or officially approved by RustDesk, Huawei, OpenHarmony, GitHub, Google, Docker, OpenWrt, or any other third-party vendor unless expressly stated by that vendor.

## Redistribution Checklist

Before publishing a binary, Docker image, release ZIP, or derived product:

- Keep `LICENSE`, `NOTICE`, and this file in source distributions.
- Keep AGPL-3.0 source availability obligations visible for network deployments.
- Preserve MIT attribution for Soybean Admin and other permissively licensed dependencies.
- Do not imply official RustDesk, Huawei, HarmonyOS or OpenHarmony endorsement.
- Do not include private keys, OAuth secrets, tokens, production passwords, user databases, or uploaded recordings in release artifacts.
- Publish corresponding source for modified AGPL deployments.
