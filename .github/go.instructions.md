---
applyTo: "**/*.go"
---

- **Test coverage**: This project has unit tests in `*_test.go` files alongside source files. Check for existing test files before suggesting missing tests. In particular, `internal/output/text_test.go` has comprehensive tests for all `TextFormatter.Format` paths, and `internal/client/email_test.go` tests all conversion helpers and patch structures.
- **MoveEmails replaces mailboxIds intentionally**: `MoveEmails` and `MarkAsSpam` in `internal/client/email.go` replace all `mailboxIds` with only the target mailbox. This is intentional move semantics, not a bug. Do not suggest preserving other mailbox associations.
- **Core client methods require a live JMAP server**: Functions like `ListEmails`, `ReadEmail`, `SearchEmails`, `MoveEmails` in the client package cannot be unit-tested without a real JMAP server. Their helper functions are tested instead. Do not flag missing unit tests for these methods.
- **Flag validation order**: Commands validate local flags (--limit, --before, --after) before calling `newClient()`. Do not suggest reordering validation before the client when it is already in that order.
- **Safety via structural omission**: `MoveEmails` and `MarkAsSpam` only populate the `Update` field of `Email/set`. The `Destroy` and `Create` fields are never set. This is the safety mechanism -- do not suggest adding `ValidateNoDestroy` or `ValidateNoSubmission` wrapper functions.
- **Search query is optional**: `jm search` accepts zero or one positional args. Filter-only search (no query, only flags) is a supported use case.
