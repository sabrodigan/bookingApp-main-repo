package forms

import (
	"testing"
)

func TestNewErrors(t *testing.T) {
	e := NewErrors()
	if e == nil {
		t.Error("expected non-nil Errors struct")
	}
}

func TestAddAndGetError(t *testing.T) {
	e := NewErrors()
	field := "email"
	msg := "invalid email"

	// Test Get on a field with no errors
	if e.Get(field) != "" {
		t.Errorf("expected empty string for field with no errors, got %s", e.Get(field))
	}

	// Test Add and Get
	e.Add(field, msg)
	if e.Get(field) != msg {
		t.Errorf("expected %s, got %s", msg, e.Get(field))
	}
}

func TestHas(t *testing.T) {
	e := NewErrors()
	field := "password"
	if e.Has(field) {
		t.Error("expected Has to be false for unset field")
	}
	e.Add(field, "required")
	if !e.Has(field) {
		t.Error("expected Has to be true after Add")
	}
}
