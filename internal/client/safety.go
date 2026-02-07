package client

import (
	"fmt"
	"strings"

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
