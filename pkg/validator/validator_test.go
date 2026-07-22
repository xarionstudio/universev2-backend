package validator

import (
	"testing"
)

type TestUser struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestValidate(t *testing.T) {
	u := TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   15,
	}

	errs := Validate(u)
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(errs))
	}

	expectedFields := map[string]bool{
		"name":  true,
		"email": true,
		"age":   true,
	}

	for _, err := range errs {
		if !expectedFields[err.Field] {
			t.Errorf("unexpected field error: %s", err.Field)
		}
	}
}
