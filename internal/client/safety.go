package client

import (
	"fmt"
	"strings"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

// ErrForbidden is returned when a safety guardrail blocks an operation.
type ErrForbidden struct {
	Operation string
	Reason    string
}

func (e *ErrForbidden) Error() string {
	return fmt.Sprintf("forbidden operation: %s: %s", e.Operation, e.Reason)
}

// trashNames are mailbox names that indicate a trash/deleted items destination.
var trashNames = []string{
	"trash",
	"deleted items",
	"deleted messages",
}

// ValidateTargetMailbox checks that a mailbox is not a trash destination.
func ValidateTargetMailbox(mb *mailbox.Mailbox) error {
	if mb.Role == mailbox.RoleTrash {
		return &ErrForbidden{
			Operation: "move",
			Reason:    fmt.Sprintf("mailbox %q has role 'trash'; deletion is not permitted", mb.Name),
		}
	}
	lower := strings.ToLower(mb.Name)
	for _, name := range trashNames {
		if lower == name {
			return &ErrForbidden{
				Operation: "move",
				Reason:    fmt.Sprintf("mailbox %q appears to be a trash folder; deletion is not permitted", mb.Name),
			}
		}
	}
	return nil
}

// ValidateSetForDraft checks that an Email/set request is a valid draft creation.
// It enforces that:
//   - Destroy is empty (no deletions)
//   - Update is empty (no modifications)
//   - Create has exactly one entry
//   - That entry targets only the given drafts mailbox
//   - That entry has the $draft keyword set
func ValidateSetForDraft(set *email.Set, draftsMailboxID jmap.ID) error {
	if len(set.Destroy) > 0 {
		return &ErrForbidden{
			Operation: "draft",
			Reason:    "Email/set destroy is not allowed in draft creation",
		}
	}
	if len(set.Update) > 0 {
		return &ErrForbidden{
			Operation: "draft",
			Reason:    "Email/set update is not allowed in draft creation",
		}
	}
	if len(set.Create) != 1 {
		return &ErrForbidden{
			Operation: "draft",
			Reason:    fmt.Sprintf("draft creation must have exactly 1 create entry, got %d", len(set.Create)),
		}
	}

	for _, e := range set.Create {
		if len(e.MailboxIDs) != 1 || !e.MailboxIDs[draftsMailboxID] {
			return &ErrForbidden{
				Operation: "draft",
				Reason:    "draft must target only the Drafts mailbox",
			}
		}
		if !e.Keywords["$draft"] {
			return &ErrForbidden{
				Operation: "draft",
				Reason:    "draft must have $draft keyword",
			}
		}
	}
	return nil
}
