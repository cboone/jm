---
applyTo: ".github/workflows/**"
---

- **scrut installation uses official install script**: The `curl | sh` pattern for installing scrut follows the project's official installation instructions. The script is fetched over HTTPS with TLS 1.2 enforcement from the official facebookincubator.github.io domain. Do not suggest replacing it with manual release artifact downloads.
