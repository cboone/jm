package client

import (
	"testing"

	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

func TestValidateTargetMailbox_AllowedMailboxes(t *testing.T) {
	allowed := []struct {
		name string
		mb   *mailbox.Mailbox
	}{
		{"archive role", &mailbox.Mailbox{Name: "Archive", Role: mailbox.RoleArchive}},
		{"inbox role", &mailbox.Mailbox{Name: "Inbox", Role: mailbox.RoleInbox}},
		{"junk role", &mailbox.Mailbox{Name: "Junk", Role: mailbox.RoleJunk}},
		{"sent role", &mailbox.Mailbox{Name: "Sent", Role: mailbox.RoleSent}},
		{"drafts role", &mailbox.Mailbox{Name: "Drafts", Role: mailbox.RoleDrafts}},
		{"custom folder", &mailbox.Mailbox{Name: "Receipts", Role: ""}},
		{"custom folder with spaces", &mailbox.Mailbox{Name: "My Important Emails", Role: ""}},
		{"folder named Trashcan", &mailbox.Mailbox{Name: "Trashcan", Role: ""}},
		{"folder named Deleted", &mailbox.Mailbox{Name: "Deleted", Role: ""}},
		{"folder named Trash Bin", &mailbox.Mailbox{Name: "Trash Bin", Role: ""}},
	}

	for _, tc := range allowed {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateTargetMailbox(tc.mb); err != nil {
				t.Errorf("expected %q to be allowed, got: %v", tc.mb.Name, err)
			}
		})
	}
}

func TestValidateTargetMailbox_BlockedMailboxes(t *testing.T) {
	blocked := []struct {
		name string
		mb   *mailbox.Mailbox
	}{
		{"trash role", &mailbox.Mailbox{Name: "Trash", Role: mailbox.RoleTrash}},
		{"trash role custom name", &mailbox.Mailbox{Name: "My Trash", Role: mailbox.RoleTrash}},
		{"trash name no role", &mailbox.Mailbox{Name: "Trash", Role: ""}},
		{"TRASH uppercase", &mailbox.Mailbox{Name: "TRASH", Role: ""}},
		{"Trash mixed case", &mailbox.Mailbox{Name: "tRaSh", Role: ""}},
		{"deleted items", &mailbox.Mailbox{Name: "Deleted Items", Role: ""}},
		{"DELETED ITEMS uppercase", &mailbox.Mailbox{Name: "DELETED ITEMS", Role: ""}},
		{"deleted messages", &mailbox.Mailbox{Name: "Deleted Messages", Role: ""}},
		{"DELETED MESSAGES uppercase", &mailbox.Mailbox{Name: "DELETED MESSAGES", Role: ""}},
	}

	for _, tc := range blocked {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTargetMailbox(tc.mb)
			if err == nil {
				t.Fatalf("expected %q to be blocked", tc.mb.Name)
			}
			if _, ok := err.(*ErrForbidden); !ok {
				t.Errorf("expected *ErrForbidden, got %T: %v", err, err)
			}
		})
	}
}

func TestErrForbidden_Error(t *testing.T) {
	err := &ErrForbidden{
		Operation: "move",
		Reason:    "mailbox is trash",
	}
	msg := err.Error()
	if msg != "forbidden operation: move: mailbox is trash" {
		t.Errorf("unexpected error message: %s", msg)
	}
}

func TestErrForbidden_ImplementsError(t *testing.T) {
	var err error = &ErrForbidden{Operation: "test", Reason: "test"}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}
