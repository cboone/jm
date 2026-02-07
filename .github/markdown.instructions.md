---
applyTo: "**/*.md"
---

- **PLAN.md documents actual behavior**: PLAN.md is a living design document that reflects the current implementation. Do not flag discrepancies that have already been updated.
- **Exit codes**: All errors exit with code 1. There are no distinct exit codes for different error categories.
- **Error format**: Runtime command errors use structured JSON to stderr. Cobra argument/flag validation errors use Cobra's default plain-text format. This is intentional.
