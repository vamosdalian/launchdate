package subscription

import (
	"strings"
	"testing"
)

func TestBuildSubscriptionStatusEmailHTML(t *testing.T) {
	html := buildSubscriptionStatusEmailHTML(
		"Your subscription status is active.",
		"You will stay on the LaunchDate subscription list for future mission updates.",
		"https://launch-date.com/api/v1/subscriptions/unsubscribe?token=abc123",
	)

	checks := []string{
		"LaunchDate",
		"Your subscription status is active.",
		"Unsubscribe",
		"token=abc123",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Fatalf("expected HTML to contain %q", check)
		}
	}
}
