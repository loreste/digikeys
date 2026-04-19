package validator

import (
	"strings"
	"testing"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=1,max=150"`
	Role  string `validate:"oneof=admin agent verifier"`
}

func TestValidateSuccess(t *testing.T) {
	s := testStruct{
		Name:  "Amadou",
		Email: "amadou@example.bf",
		Age:   30,
		Role:  "admin",
	}
	if err := Validate(s); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateRequiredField(t *testing.T) {
	s := testStruct{
		Email: "test@example.bf",
		Age:   25,
		Role:  "admin",
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for missing Name")
	}
	if !strings.Contains(err.Error(), "obligatoire") {
		t.Errorf("expected French 'obligatoire' in error, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "Name") {
		t.Errorf("expected field name 'Name' in error, got: %s", err.Error())
	}
}

func TestValidateEmailField(t *testing.T) {
	s := testStruct{
		Name:  "Amadou",
		Email: "not-an-email",
		Age:   25,
		Role:  "admin",
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for invalid email")
	}
	if !strings.Contains(err.Error(), "email valide") {
		t.Errorf("expected French 'email valide' in error, got: %s", err.Error())
	}
}

func TestValidateMinField(t *testing.T) {
	s := testStruct{
		Name:  "Amadou",
		Email: "test@example.bf",
		Age:   0,
		Role:  "admin",
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for Age below min")
	}
	if !strings.Contains(err.Error(), "au moins") {
		t.Errorf("expected French 'au moins' in error, got: %s", err.Error())
	}
}

func TestValidateOneOfField(t *testing.T) {
	s := testStruct{
		Name:  "Amadou",
		Email: "test@example.bf",
		Age:   25,
		Role:  "invalid_role",
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for invalid oneof")
	}
	if !strings.Contains(err.Error(), "l'une des valeurs") {
		t.Errorf("expected French 'l'une des valeurs' in error, got: %s", err.Error())
	}
}

func TestValidateMultipleErrors(t *testing.T) {
	s := testStruct{} // All fields invalid
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation errors")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	// At least Name (required), Email (required) should fail
	if len(ve.Messages) < 2 {
		t.Errorf("expected at least 2 validation messages, got %d: %v", len(ve.Messages), ve.Messages)
	}
}

func TestValidateMaxField(t *testing.T) {
	type maxStruct struct {
		Code string `validate:"max=5"`
	}
	s := maxStruct{Code: "toolongvalue"}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for max exceeded")
	}
	if !strings.Contains(err.Error(), "ne doit pas") {
		t.Errorf("expected French 'ne doit pas' in error, got: %s", err.Error())
	}
}

func TestValidationErrorString(t *testing.T) {
	ve := &ValidationError{Messages: []string{"erreur 1", "erreur 2"}}
	got := ve.Error()
	if got != "erreur 1; erreur 2" {
		t.Errorf("Error() = %q, want %q", got, "erreur 1; erreur 2")
	}
}

func TestValidateNilInput(t *testing.T) {
	// Passing a non-struct should return an error (not panic)
	err := Validate("not a struct")
	if err == nil {
		t.Fatal("expected error when validating a non-struct")
	}
}
