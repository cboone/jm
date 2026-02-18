# Add draft creation command

## Context

Issue #36: `fm` has no way to create email drafts. This is needed for triage workflows where Claude Code composes a draft (e.g., an unsubscribe reply) for the user to review and send from Fastmail.

Draft creation uses `Email/set` with `Create` -- a separate JMAP operation from `EmailSubmission/set` (sending). No new permissions or scopes are needed.

This plan follows the design in `docs/plans/done/2026-02-11-email-drafting-option-c.md`.

---

## Implementation

### 1. Add `DraftResult` type -- `internal/types/types.go`

```go
type DraftResult struct {
    ID        string           `json:"id"`
    Mode      string           `json:"mode"`
    Mailbox   *DestinationInfo `json:"mailbox"`
    From      []Address        `json:"from,omitempty"`
    To        []Address        `json:"to"`
    CC        []Address        `json:"cc,omitempty"`
    Subject   string           `json:"subject"`
    InReplyTo string           `json:"in_reply_to,omitempty"`
}
```

### 2. Add `ValidateSetForDraft` -- `internal/client/safety.go`

Validates the entire `Email/set` request before execution:
- `Destroy` must be empty
- `Update` must be empty
- `Create` must have exactly one entry
- That entry's `MailboxIDs` must contain only the Drafts mailbox ID
- That entry's `Keywords` must include `$draft: true`

### 3. Add safety tests -- `internal/client/safety_test.go`

Test cases for `ValidateSetForDraft`:
- Valid draft passes
- Reject non-empty `Destroy`
- Reject non-empty `Update`
- Reject empty or multiple `Create` entries
- Reject wrong mailbox target
- Reject missing `$draft` keyword

### 4. Create `internal/client/draft.go`

Types and methods:
- `DraftMode` type (new, reply, reply-all, forward)
- `DraftOptions` struct (mode, to, cc, bcc, subject, body, html flag, original email ID)
- `CreateDraft(opts DraftOptions) (types.DraftResult, error)` -- main method
- `fetchOriginalForReply(emailID string) (*email.Email, error)` -- fetches original with threading properties

`CreateDraft` flow:
1. Resolve Drafts mailbox via `GetMailboxByRole(mailbox.RoleDrafts)`
2. Derive `From` from `Session().Username` (omit if not a valid email)
3. For reply/reply-all/forward: fetch original email, compute recipients/subject/threading headers per the rules in the existing plan
4. Construct `email.Email` with `MailboxIDs`, `Keywords` (`$draft`, `$seen`), body, recipients
5. Call `ValidateSetForDraft` on the constructed `email.Set`
6. Execute via `c.Do(req)`, handle `Created`/`NotCreated` in response
7. Return `DraftResult`

Recipient computation rules (from existing plan):
- **Reply:** `To` = original `ReplyTo` (or `From`); user `--to` appended, deduped
- **Reply-all:** same `To` base; `CC` = original `To+CC` minus self minus final `To`; user `--cc` appended
- **Forward:** recipients from user flags only
- **New:** recipients from user flags only

Threading header rules:
- `InReplyTo` = original `messageId` list
- `References` = original `references` + original `messageId`, order-preserving deduped
- Omit both if original has no `messageId`

Subject prefixing:
- Reply/reply-all: prepend `Re: ` if not already present
- Forward: prepend `Fwd: ` if not already present

### 5. Create tests -- `internal/client/draft_test.go`

Using `Client.doFunc` mocking (same pattern as `email_test.go`):
- New draft creation (basic fields, HTML mode)
- Reply recipient derivation (ReplyTo fallback to From)
- Reply-all To/CC composition, self-exclusion, deduplication
- Forward subject prefixing and quoted body
- Threading headers from messageId/references
- From behavior (valid username sets From, invalid omits)
- Safety guard integration (wrong mailbox, missing $draft)
- Error cases (Drafts mailbox missing, original email missing, NotCreated, method errors)
- Address deduplication across user-supplied and auto-derived recipients

### 6. Create `cmd/draft.go`

Command: `fm draft [flags]`

Flags:
| Flag | Description |
|------|-------------|
| `--to` | Recipient addresses (required for new/forward) |
| `--cc` | CC addresses |
| `--bcc` | BCC addresses |
| `--subject` | Subject line (required for new) |
| `--body` | Message body (mutually exclusive with `--body-stdin`) |
| `--body-stdin` | Read body from stdin |
| `--reply-to` | Email ID to reply to |
| `--reply-all` | Email ID to reply-all to |
| `--forward` | Email ID to forward |
| `--html` | Treat body as HTML |

Address flags parsed with RFC 5322 (`net/mail`).

Validation:
- Exactly one of: no mode flag (new), `--reply-to`, `--reply-all`, `--forward`
- Mode flags are mutually exclusive
- New mode requires `--to` and `--subject`
- Forward mode requires `--to`
- Exactly one of `--body` or `--body-stdin`

Command flow:
1. Validate flags and determine mode
2. Read body from flag or stdin
3. Parse address flags
4. Build `DraftOptions`
5. Call `newClient()` then `c.CreateDraft(opts)`
6. Format result via `formatter().Format()`

### 7. Add text formatter -- `internal/output/text.go`

Add `types.DraftResult` case to the `Format` switch and a `formatDraftResult` method:

```
Draft created: M-new-draft-id
Mode: reply
To: Alice <alice@example.com>
Subject: Re: Meeting
Mailbox: Drafts (mb-drafts-id)
In-Reply-To: <CAExample1234@example.com>
```

### 8. Update documentation

**`README.md`**: Add draft command to the command tables, add usage examples, update safety section to clarify that draft creation (not sending) is allowed.

**`docs/CLI-REFERENCE.md`**: Add full `### draft` section with flags, validation rules, output schemas (JSON and text), and examples for all four modes. Add `DraftResult` to the Output Schemas section.

**`tests/help.md`**: Add `fm draft --help` scrut test block. Update root help to include `draft` in the command list.

**`docs/plans/done/PLAN.md`**: Update the "Disallowed Operations" wording to clarify that `Email/set create` is allowed only for validated drafts (Drafts mailbox + $draft keyword), while `EmailSubmission/set` (sending) remains structurally disallowed.

---

## Files summary

| File | Action |
|------|--------|
| `internal/types/types.go` | Modify: add `DraftResult` |
| `internal/client/safety.go` | Modify: add `ValidateSetForDraft` |
| `internal/client/safety_test.go` | Modify: add safety validation tests |
| `internal/client/draft.go` | Create: `CreateDraft`, `fetchOriginalForReply`, helpers |
| `internal/client/draft_test.go` | Create: comprehensive draft tests |
| `cmd/draft.go` | Create: draft command |
| `internal/output/text.go` | Modify: add `DraftResult` formatting |
| `README.md` | Modify: add draft command docs |
| `docs/CLI-REFERENCE.md` | Modify: add draft section and schema |
| `tests/help.md` | Modify: add draft help test |
| `docs/plans/done/PLAN.md` | Modify: update safety wording |

---

## Verification

1. `make build` compiles
2. `go test ./...` passes (all new and existing tests)
3. `make test-cli` passes (scrut snapshot tests including new draft help)
4. `go vet ./...` and `gofmt` clean
5. Manual live test: `fm draft --to test@example.com --subject "Test" --body "Hello"` -- verify draft appears in Fastmail Drafts
6. Manual reply test: `fm draft --reply-to <email-id> --body "Thanks"` -- verify recipients and threading
