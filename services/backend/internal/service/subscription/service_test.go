package subscription

import "testing"

func TestNormalizeEmail(t *testing.T) {
	got, err := NormalizeEmail("  USER@example.COM ")
	if err != nil {
		t.Fatalf("expected email to be valid, got error: %v", err)
	}

	if got != "user@example.com" {
		t.Fatalf("expected normalized email to be user@example.com, got %q", got)
	}
}

func TestNormalizeEmailInvalid(t *testing.T) {
	_, err := NormalizeEmail("not-an-email")
	if err == nil {
		t.Fatal("expected invalid email error")
	}
}
