package output

import "io"

// Formatter renders structured data to an output stream.
type Formatter interface {
	// Format writes v to the output stream.
	Format(w io.Writer, v any) error
	// FormatError writes a structured error to the output stream.
	FormatError(w io.Writer, code string, message string, hint string) error
}

// New returns a Formatter for the given format name ("json" or "text").
func New(format string) Formatter {
	if format == "text" {
		return &TextFormatter{}
	}
	return &JSONFormatter{}
}
