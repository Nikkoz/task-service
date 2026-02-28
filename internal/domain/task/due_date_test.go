package task

import (
	"errors"
	"testing"
	"time"
)

func TestNewDueDate_TimeInThePast_IsInvalid(t *testing.T) {
	date := time.Now().Add(-time.Hour)
	_, err := NewDueDate(date)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrDateInThePast) {
		t.Fatalf("expected ErrDateInThePast, got %v", err)
	}
}

func TestNewDueDate_TimeIsEqualNow_IsInvalid(t *testing.T) {
	date := time.Now()
	_, err := NewDueDate(date)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrDateInThePast) {
		t.Fatalf("expected ErrDateInThePast, got %v", err)
	}
}

func TestNewDueDate_IsValid(t *testing.T) {
	date := time.Now().Add(time.Hour)
	_, err := NewDueDate(date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
