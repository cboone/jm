package client

import (
	"testing"

	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

func TestValidateTargetMailbox_AllowsArchive(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Archive", Role: mailbox.RoleArchive}
	if err := ValidateTargetMailbox(mb); err != nil {
		t.Errorf("expected archive to be allowed, got: %v", err)
	}
}

func TestValidateTargetMailbox_AllowsInbox(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Inbox", Role: mailbox.RoleInbox}
	if err := ValidateTargetMailbox(mb); err != nil {
		t.Errorf("expected inbox to be allowed, got: %v", err)
	}
}

func TestValidateTargetMailbox_AllowsCustomFolder(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Receipts", Role: ""}
	if err := ValidateTargetMailbox(mb); err != nil {
		t.Errorf("expected custom folder to be allowed, got: %v", err)
	}
}

func TestValidateTargetMailbox_BlocksTrashRole(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Trash", Role: mailbox.RoleTrash}
	err := ValidateTargetMailbox(mb)
	if err == nil {
		t.Fatal("expected trash role to be blocked")
	}
	if _, ok := err.(*ErrForbidden); !ok {
		t.Errorf("expected ErrForbidden, got: %T", err)
	}
}

func TestValidateTargetMailbox_BlocksTrashName(t *testing.T) {
	// Trash name without the role set (edge case).
	mb := &mailbox.Mailbox{Name: "Trash", Role: ""}
	err := ValidateTargetMailbox(mb)
	if err == nil {
		t.Fatal("expected 'Trash' name to be blocked")
	}
}

func TestValidateTargetMailbox_BlocksDeletedItems(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Deleted Items", Role: ""}
	err := ValidateTargetMailbox(mb)
	if err == nil {
		t.Fatal("expected 'Deleted Items' to be blocked")
	}
}

func TestValidateTargetMailbox_BlocksDeletedMessages(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "Deleted Messages", Role: ""}
	err := ValidateTargetMailbox(mb)
	if err == nil {
		t.Fatal("expected 'Deleted Messages' to be blocked")
	}
}

func TestValidateTargetMailbox_CaseInsensitive(t *testing.T) {
	mb := &mailbox.Mailbox{Name: "TRASH", Role: ""}
	err := ValidateTargetMailbox(mb)
	if err == nil {
		t.Fatal("expected case-insensitive 'TRASH' to be blocked")
	}
}
