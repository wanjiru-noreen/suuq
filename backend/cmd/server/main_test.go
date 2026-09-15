package main

import (
	"strings"
	"testing"
)

func TestValidateSecret(t *testing.T) {
	for _, secret := range []string{"", "development-only-change-me", strings.Repeat(" ", 64), strings.Repeat("a", 31)} {
		if validateSecret(secret) == nil {
			t.Fatal("accepted missing or short secret")
		}
	}
	if err := validateSecret(strings.Repeat("a", 64)); err != nil {
		t.Fatal(err)
	}
}
