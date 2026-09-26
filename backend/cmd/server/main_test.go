package main

import (
	"strings"
	"testing"
)

func TestValidateSecret(t *testing.T) {
	for _, secret := range []string{
		"",
		"development-only-change-me",
		strings.Repeat(" ", 64),
		strings.Repeat("a", 31),
		"replace-with-a-random-secret-of-at-least-32-bytes",
	} {
		if validateSecret(secret) == nil {
			t.Fatalf("accepted invalid secret %q", secret)
		}
	}
	if err := validateSecret(strings.Repeat("a", 64)); err != nil {
		t.Fatal(err)
	}
}
