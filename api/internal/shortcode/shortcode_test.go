package shortcode

import (
	"strings"
	"testing"
)

func TestGenerateLength(t *testing.T) {
	for _, n := range []int{1, 7, 16, 64} {
		code, err := Generate(n)
		if err != nil {
			t.Fatalf("unexpected error for length %d: %v", n, err)
		}
		if len(code) != n {
			t.Fatalf("expected length %d, got %d", n, len(code))
		}
	}
}

func TestGenerateAlphabet(t *testing.T) {
	code, err := Generate(128)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range code {
		if !strings.ContainsRune(alphabet, r) {
			t.Fatalf("character %q not in alphabet", r)
		}
	}
}

func TestGenerateRejectsNonPositive(t *testing.T) {
	if _, err := Generate(0); err == nil {
		t.Fatal("expected error for length 0")
	}
	if _, err := Generate(-3); err == nil {
		t.Fatal("expected error for negative length")
	}
}

func TestGenerateUnique(t *testing.T) {
	seen := make(map[string]struct{}, 200)
	for i := 0; i < 200; i++ {
		code, err := Generate(10)
		if err != nil {
			t.Fatal(err)
		}
		if _, dup := seen[code]; dup {
			t.Fatalf("duplicate code generated: %s", code)
		}
		seen[code] = struct{}{}
	}
}
