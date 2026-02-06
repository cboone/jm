package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRetryTransport_RetriesWithFreshRequestBody(t *testing.T) {
	var calls int
	var bodies []string

	rt := &retryTransport{base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		bodies = append(bodies, string(payload))

		if calls == 1 {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Header:     http.Header{"Retry-After": []string{"0"}},
				Body:       io.NopCloser(strings.NewReader("busy")),
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	})}

	req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString("payload"))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("roundtrip failed: %v", err)
	}
	defer resp.Body.Close()

	if calls != 2 {
		t.Fatalf("expected 2 transport calls, got %d", calls)
	}
	if len(bodies) != 2 {
		t.Fatalf("expected 2 captured bodies, got %d", len(bodies))
	}
	if bodies[0] != "payload" || bodies[1] != "payload" {
		t.Fatalf("expected both attempts to send full payload, got %q", bodies)
	}
}

type oneShotReader struct {
	data []byte
	off  int
}

func (r *oneShotReader) Read(p []byte) (int, error) {
	if r.off >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.off:])
	r.off += n
	return n, nil
}

func TestRetryTransport_FailsForNonRewindableBody(t *testing.T) {
	var calls int

	rt := &retryTransport{base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		_, _ = io.ReadAll(req.Body)
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header:     http.Header{"Retry-After": []string{"0"}},
			Body:       io.NopCloser(strings.NewReader("retry later")),
		}, nil
	})}

	req, err := http.NewRequest(http.MethodPost, "http://example.com", &oneShotReader{data: []byte("payload")})
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	_, err = rt.RoundTrip(req)
	if err == nil {
		t.Fatal("expected retry to fail for non-rewindable body")
	}
	if !strings.Contains(err.Error(), "non-rewindable body") {
		t.Fatalf("expected non-rewindable body error, got: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 transport call, got %d", calls)
	}
}
