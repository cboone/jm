package cmd

import "testing"

func TestParseSort_Default(t *testing.T) {
	field, asc := parseSort("receivedAt desc")
	if field != "receivedAt" {
		t.Errorf("expected field=receivedAt, got %s", field)
	}
	if asc {
		t.Error("expected ascending=false for desc")
	}
}

func TestParseSort_Ascending(t *testing.T) {
	field, asc := parseSort("sentAt asc")
	if field != "sentAt" {
		t.Errorf("expected field=sentAt, got %s", field)
	}
	if !asc {
		t.Error("expected ascending=true for asc")
	}
}

func TestParseSort_AscendingCaseInsensitive(t *testing.T) {
	field, asc := parseSort("subject ASC")
	if field != "subject" {
		t.Errorf("expected field=subject, got %s", field)
	}
	if !asc {
		t.Error("expected ascending=true for ASC")
	}
}

func TestParseSort_FieldOnly(t *testing.T) {
	field, asc := parseSort("from")
	if field != "from" {
		t.Errorf("expected field=from, got %s", field)
	}
	if asc {
		t.Error("expected ascending=false when direction omitted")
	}
}

func TestParseSort_Empty(t *testing.T) {
	field, asc := parseSort("")
	if field != "receivedAt" {
		t.Errorf("expected default field=receivedAt, got %s", field)
	}
	if asc {
		t.Error("expected ascending=false for empty input")
	}
}

func TestParseSort_ExtraWhitespace(t *testing.T) {
	field, asc := parseSort("  receivedAt   asc  ")
	if field != "receivedAt" {
		t.Errorf("expected field=receivedAt, got %s", field)
	}
	if !asc {
		t.Error("expected ascending=true")
	}
}

func TestParseSort_UnknownDirection(t *testing.T) {
	field, asc := parseSort("receivedAt up")
	if field != "receivedAt" {
		t.Errorf("expected field=receivedAt, got %s", field)
	}
	if asc {
		t.Error("expected ascending=false for unknown direction")
	}
}
