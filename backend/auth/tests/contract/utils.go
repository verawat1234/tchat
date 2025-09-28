package contract

import (
	"io"
	"strings"
)

// Utility functions for contract testing

// StringPtr returns a pointer to string value
func StringPtr(s string) *string {
	return &s
}

// NewStringReader creates a new string reader (like pact-go utils)
func NewStringReader(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}