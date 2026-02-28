package task

import (
	"errors"
	"testing"
)

func TestNewDescription_TrimsSpaces(t *testing.T) {
	const expected = "Buy milk description"
	got, err := NewDescription("  Buy milk description  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.String() != expected {
		t.Fatalf("expected %q, got %q", expected, got.String())
	}
}

func TestNewDescription_EmptyAfterTrim_IsInvalid(t *testing.T) {
	_, err := NewDescription("   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrEmptyDescription) {
		t.Fatalf("expected ErrEmptyDescription, got %v", err)
	}
}

func TestNewDescription_MaxLen_IsValid(t *testing.T) {
	_, err := NewDescription("Buy milk description")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
