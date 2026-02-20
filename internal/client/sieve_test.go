package client

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"

	"github.com/cboone/fm/internal/jmap/sieve"
)

func sieveTestClient(doFunc func(*jmap.Request) (*jmap.Response, error)) *Client {
	return &Client{
		jmap: &jmap.Client{
			Session: &jmap.Session{
				RawCapabilities: map[jmap.URI]json.RawMessage{
					sieve.URI: json.RawMessage("{}"),
				},
			},
		},
		accountID: "acct-1",
		doFunc:    doFunc,
	}
}

func TestListSieveScripts(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{
						List: []*sieve.SieveScript{
							{ID: "S1", Name: "Block spam", IsActive: true},
							{ID: "S2", Name: "Custom filter", IsActive: false},
						},
					},
				},
			},
		}, nil
	})

	result, err := c.ListSieveScripts()
	if err != nil {
		t.Fatalf("ListSieveScripts() error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("ListSieveScripts() total = %d, want 2", result.Total)
	}
	if result.Scripts[0].ID != "S1" {
		t.Errorf("ListSieveScripts() scripts[0].ID = %q, want %q", result.Scripts[0].ID, "S1")
	}
	if !result.Scripts[0].IsActive {
		t.Error("ListSieveScripts() scripts[0].IsActive = false, want true")
	}
	if result.Scripts[1].Name != "Custom filter" {
		t.Errorf("ListSieveScripts() scripts[1].Name = %q, want %q", result.Scripts[1].Name, "Custom filter")
	}
}

func TestListSieveScripts_Empty(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{List: nil},
				},
			},
		}, nil
	})

	result, err := c.ListSieveScripts()
	if err != nil {
		t.Fatalf("ListSieveScripts() error: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("ListSieveScripts() total = %d, want 0", result.Total)
	}
}

func TestListSieveScripts_MethodError(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "error",
					Args: &jmap.MethodError{Type: "serverFail"},
				},
			},
		}, nil
	})

	_, err := c.ListSieveScripts()
	if err == nil {
		t.Fatal("ListSieveScripts() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "serverFail") {
		t.Errorf("ListSieveScripts() error = %q, want to contain %q", err.Error(), "serverFail")
	}
}

func TestListSieveScripts_NoCapability(t *testing.T) {
	c := &Client{
		jmap: &jmap.Client{
			Session: &jmap.Session{
				RawCapabilities: map[jmap.URI]json.RawMessage{},
			},
		},
		accountID: "acct-1",
	}

	_, err := c.ListSieveScripts()
	if err == nil {
		t.Fatal("ListSieveScripts() expected error for missing capability")
	}
	if !strings.Contains(err.Error(), "does not support sieve") {
		t.Errorf("ListSieveScripts() error = %q, want capability error", err.Error())
	}
}

func TestGetSieveScript(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{
						List: []*sieve.SieveScript{
							{ID: "S1", Name: "Block spam", BlobID: "B1", IsActive: true},
						},
					},
				},
			},
		}, nil
	})
	c.downloadFunc = func(accountID, blobID jmap.ID) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("require [\"fileinto\"];\nkeep;\n")), nil
	}

	result, err := c.GetSieveScript("S1")
	if err != nil {
		t.Fatalf("GetSieveScript() error: %v", err)
	}
	if result.ID != "S1" {
		t.Errorf("GetSieveScript() ID = %q, want %q", result.ID, "S1")
	}
	if result.Content != "require [\"fileinto\"];\nkeep;\n" {
		t.Errorf("GetSieveScript() content = %q, want sieve script", result.Content)
	}
}

func TestGetSieveScript_NotFound(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{
						NotFound: []jmap.ID{"S99"},
					},
				},
			},
		}, nil
	})

	_, err := c.GetSieveScript("S99")
	if err == nil {
		t.Fatal("GetSieveScript() expected error for not found")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("GetSieveScript() error = %q, want not found error", err.Error())
	}
}

func TestCreateSieveScript(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{
						Created: map[jmap.ID]*sieve.SieveScript{
							"create0": {ID: "S-new"},
						},
					},
				},
			},
		}, nil
	})
	c.uploadFunc = func(accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
		return &jmap.UploadResponse{ID: "B-new"}, nil
	}

	result, err := c.CreateSieveScript("test script", "keep;\n", false)
	if err != nil {
		t.Fatalf("CreateSieveScript() error: %v", err)
	}
	if result.ID != "S-new" {
		t.Errorf("CreateSieveScript() ID = %q, want %q", result.ID, "S-new")
	}
	if result.Name != "test script" {
		t.Errorf("CreateSieveScript() Name = %q, want %q", result.Name, "test script")
	}
	if result.IsActive {
		t.Error("CreateSieveScript() IsActive = true, want false")
	}
}

func TestCreateSieveScript_WithActivation(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{
						Created: map[jmap.ID]*sieve.SieveScript{
							"create0": {ID: "S-new"},
						},
					},
				},
			},
		}, nil
	})
	c.uploadFunc = func(accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
		return &jmap.UploadResponse{ID: "B-new"}, nil
	}

	result, err := c.CreateSieveScript("active script", "keep;\n", true)
	if err != nil {
		t.Fatalf("CreateSieveScript() error: %v", err)
	}
	if !result.IsActive {
		t.Error("CreateSieveScript() IsActive = false, want true")
	}
}

func TestCreateSieveScript_NotCreated(t *testing.T) {
	desc := "name already exists"
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{
						NotCreated: map[jmap.ID]*jmap.SetError{
							"create0": {Type: "invalidProperties", Description: &desc},
						},
					},
				},
			},
		}, nil
	})
	c.uploadFunc = func(accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
		return &jmap.UploadResponse{ID: "B-new"}, nil
	}

	_, err := c.CreateSieveScript("test", "keep;\n", false)
	if err == nil {
		t.Fatal("CreateSieveScript() expected error for NotCreated")
	}
	if !strings.Contains(err.Error(), "name already exists") {
		t.Errorf("CreateSieveScript() error = %q, want description", err.Error())
	}
}

func TestValidateSieveScript_Valid(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/validate",
					Args: &sieve.ValidateResponse{Error: nil},
				},
			},
		}, nil
	})
	c.uploadFunc = func(accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
		return &jmap.UploadResponse{ID: "B-val"}, nil
	}

	result, err := c.ValidateSieveScript("keep;\n")
	if err != nil {
		t.Fatalf("ValidateSieveScript() error: %v", err)
	}
	if !result.Valid {
		t.Error("ValidateSieveScript() valid = false, want true")
	}
}

func TestValidateSieveScript_Invalid(t *testing.T) {
	desc := "unsupported extension"
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/validate",
					Args: &sieve.ValidateResponse{
						Error: &jmap.SetError{Type: "invalidSieve", Description: &desc},
					},
				},
			},
		}, nil
	})
	c.uploadFunc = func(accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
		return &jmap.UploadResponse{ID: "B-val"}, nil
	}

	result, err := c.ValidateSieveScript("bad script")
	if err != nil {
		t.Fatalf("ValidateSieveScript() error: %v", err)
	}
	if result.Valid {
		t.Error("ValidateSieveScript() valid = true, want false")
	}
	if result.Error != "unsupported extension" {
		t.Errorf("ValidateSieveScript() error = %q, want %q", result.Error, "unsupported extension")
	}
}

func TestActivateSieveScript(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{},
				},
			},
		}, nil
	})

	result, err := c.ActivateSieveScript("S1")
	if err != nil {
		t.Fatalf("ActivateSieveScript() error: %v", err)
	}
	if result.ID != "S1" {
		t.Errorf("ActivateSieveScript() ID = %q, want %q", result.ID, "S1")
	}
	if !result.IsActive {
		t.Error("ActivateSieveScript() IsActive = false, want true")
	}
}

func TestDeactivateSieveScript(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{},
				},
			},
		}, nil
	})

	result, err := c.DeactivateSieveScript()
	if err != nil {
		t.Fatalf("DeactivateSieveScript() error: %v", err)
	}
	if result.IsActive {
		t.Error("DeactivateSieveScript() IsActive = true, want false")
	}
}

func TestDeleteSieveScript(t *testing.T) {
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{
						List: []*sieve.SieveScript{
							{ID: "S1", Name: "my-filter"},
						},
					},
				},
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{
						Destroyed: []jmap.ID{"S1"},
					},
				},
			},
		}, nil
	})

	result, err := c.DeleteSieveScript("S1")
	if err != nil {
		t.Fatalf("DeleteSieveScript() error: %v", err)
	}
	if result.ID != "S1" {
		t.Errorf("DeleteSieveScript() ID = %q, want %q", result.ID, "S1")
	}
	if result.Name != "my-filter" {
		t.Errorf("DeleteSieveScript() Name = %q, want %q", result.Name, "my-filter")
	}
}

func TestDeleteSieveScript_ActiveScript(t *testing.T) {
	desc := "script is currently active"
	c := sieveTestClient(func(req *jmap.Request) (*jmap.Response, error) {
		return &jmap.Response{
			Responses: []*jmap.Invocation{
				{
					Name: "SieveScript/get",
					Args: &sieve.GetResponse{
						List: []*sieve.SieveScript{
							{ID: "S1", Name: "my-filter"},
						},
					},
				},
				{
					Name: "SieveScript/set",
					Args: &sieve.SetResponse{
						NotDestroyed: map[jmap.ID]*jmap.SetError{
							"S1": {Type: "sieveIsActive", Description: &desc},
						},
					},
				},
			},
		}, nil
	})

	_, err := c.DeleteSieveScript("S1")
	if err == nil {
		t.Fatal("DeleteSieveScript() expected error for active script")
	}
	if !strings.Contains(err.Error(), "deactivate it first") {
		t.Errorf("DeleteSieveScript() error = %q, want deactivation hint", err.Error())
	}
}
