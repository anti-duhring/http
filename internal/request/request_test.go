package request_test

import (
	"httpfromtcp/internal/request"
	"strings"
	"testing"
)

func TestRequestFromReaderRequestLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		wantMethod  string
		wantTarget  string
		wantVersion string
	}{
		{
			name:        "Good GET Request line",
			input:       "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantErr:     false,
			wantMethod:  "GET",
			wantTarget:  "/",
			wantVersion: "1.1",
		},
		{
			name:        "Good GET Request line with path",
			input:       "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantErr:     false,
			wantMethod:  "GET",
			wantTarget:  "/coffee",
			wantVersion: "1.1",
		},
		{
			name:    "Invalid number of parts in request line",
			input:   "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := request.RequestFromReader(strings.NewReader(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r == nil {
				t.Fatalf("expected request, got nil")
			}
			if got := r.RequestLine.Method; got != tc.wantMethod {
				t.Errorf("Method = %q, want %q", got, tc.wantMethod)
			}
			if got := r.RequestLine.RequestTarget; got != tc.wantTarget {
				t.Errorf("RequestTarget = %q, want %q", got, tc.wantTarget)
			}
			if got := r.RequestLine.HttpVersion; got != tc.wantVersion {
				t.Errorf("HttpVersion = %q, want %q", got, tc.wantVersion)
			}
		})
	}
}
