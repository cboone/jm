---
applyTo: "**/*.go"
---

- **Test coverage**: This project has unit tests in `*_test.go` files alongside source files. Check for existing test files before suggesting missing tests. In particular, `internal/output/text_test.go` has comprehensive tests for all `TextFormatter.Format` paths, and `internal/client/email_test.go` tests all conversion helpers and patch structures.
- **MoveEmails replaces mailboxIds intentionally**: `MoveEmails` and `MarkAsSpam` in `internal/client/email.go` replace all `mailboxIds` with only the target mailbox. This is intentional move semantics, not a bug. Do not suggest preserving other mailbox associations.
- **Core client methods require a live JMAP server**: Functions like `ListEmails`, `ReadEmail`, `SearchEmails`, `MoveEmails` in the client package cannot be unit-tested without a real JMAP server. Their helper functions are tested instead. Do not flag missing unit tests for these methods.
