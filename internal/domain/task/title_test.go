package task

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTitle_TrimsSpaces(t *testing.T) {
	const expected = "Buy milk"
	got, err := NewTitle("  Buy milk  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.String() != expected {
		t.Fatalf("expected %q, got %q", expected, got.String())
	}
}

func TestNewTitle_EmptyAfterTrim_IsInvalid(t *testing.T) {
	_, err := NewTitle("   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestNewTitle_MaxLen_IsInvalid(t *testing.T) {
	long := strings.Repeat("a", MaxLength+1)
	_, err := NewTitle(long)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrWrongLength) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestNewTitle_MaxLen_IsValid(t *testing.T) {
	exact := strings.Repeat("a", MaxLength)
	_, err := NewTitle(exact)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
