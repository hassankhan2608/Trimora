package validate

import (
	"errors"
	"strings"
	"testing"
)

func TestURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want error
	}{
		{"empty", "", ErrURLEmpty},
		{"whitespace", "   ", ErrURLEmpty},
		{"no scheme", "example.com", ErrURLScheme},
		{"ftp", "ftp://example.com", ErrURLScheme},
		{"no host", "http://", ErrURLInvalid},
		{"too long", "http://" + strings.Repeat("a", MaxURLLength), ErrURLTooLong},
		{"valid http", "http://example.com/path?q=1", nil},
		{"valid https", "https://example.com", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := URL(tc.in)
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}

func TestAlias(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want error
	}{
		{"too short", "ab", ErrAliasLength},
		{"too long", strings.Repeat("a", MaxAliasLength+1), ErrAliasLength},
		{"bad char", "hello!", ErrAliasFormat},
		{"reserved", "api", ErrAliasReserved},
		{"reserved upper", "API", ErrAliasReserved},
		{"valid", "my-link_1", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Alias(tc.in)
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}
