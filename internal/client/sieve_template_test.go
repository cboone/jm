package client

import "testing"

func TestGenerateSieveScript(t *testing.T) {
	tests := []struct {
		name    string
		opts    SieveTemplateOptions
		want    string
		wantErr string
	}{
		{
			name: "from address with junk action",
			opts: SieveTemplateOptions{From: "spam@example.com", Action: "junk"},
			want: "require [\"fileinto\"];\n\nif address :is \"from\" \"spam@example.com\" {\n    fileinto \"Junk\";\n    stop;\n}\n",
		},
		{
			name: "from domain with junk action",
			opts: SieveTemplateOptions{FromDomain: "example.com", Action: "junk"},
			want: "require [\"fileinto\"];\n\nif address :domain :is \"from\" \"example.com\" {\n    fileinto \"Junk\";\n    stop;\n}\n",
		},
		{
			name: "from address with discard action",
			opts: SieveTemplateOptions{From: "spam@example.com", Action: "discard"},
			want: "if address :is \"from\" \"spam@example.com\" {\n    discard;\n    stop;\n}\n",
		},
		{
			name: "from address with keep action",
			opts: SieveTemplateOptions{From: "important@example.com", Action: "keep"},
			want: "if address :is \"from\" \"important@example.com\" {\n    keep;\n    stop;\n}\n",
		},
		{
			name: "from domain with fileinto action",
			opts: SieveTemplateOptions{FromDomain: "example.com", Action: "fileinto", FileInto: "Archive"},
			want: "require [\"fileinto\"];\n\nif address :domain :is \"from\" \"example.com\" {\n    fileinto \"Archive\";\n    stop;\n}\n",
		},
		{
			name:    "missing from and from-domain",
			opts:    SieveTemplateOptions{Action: "junk"},
			wantErr: "either --from or --from-domain is required",
		},
		{
			name:    "both from and from-domain",
			opts:    SieveTemplateOptions{From: "a@b.com", FromDomain: "b.com", Action: "junk"},
			wantErr: "--from and --from-domain are mutually exclusive",
		},
		{
			name:    "missing action",
			opts:    SieveTemplateOptions{From: "a@b.com"},
			wantErr: "--action is required",
		},
		{
			name:    "unsupported action",
			opts:    SieveTemplateOptions{From: "a@b.com", Action: "delete"},
			wantErr: "unsupported action \"delete\": use junk, discard, keep, or fileinto",
		},
		{
			name:    "fileinto action without mailbox",
			opts:    SieveTemplateOptions{From: "a@b.com", Action: "fileinto"},
			wantErr: "--fileinto is required when --action is fileinto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSieveScript(tt.opts)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("GenerateSieveScript() error = nil, want error containing %q", tt.wantErr)
					return
				}
				if got := err.Error(); got != tt.wantErr {
					t.Errorf("GenerateSieveScript() error = %q, want %q", got, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("GenerateSieveScript() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateSieveScript() =\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}
