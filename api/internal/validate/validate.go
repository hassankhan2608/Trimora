package validate

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	MinAliasLength = 3
	MaxAliasLength = 32
	MaxURLLength   = 2048
)

var (
	ErrURLEmpty      = errors.New("url is required")
	ErrURLTooLong    = errors.New("url is too long")
	ErrURLInvalid    = errors.New("url is invalid")
	ErrURLScheme     = errors.New("url must use http or https")
	ErrAliasLength   = errors.New("alias must be between 3 and 32 characters")
	ErrAliasFormat   = errors.New("alias may only contain letters, numbers, hyphens, and underscores")
	ErrAliasReserved = errors.New("alias is reserved")
	ErrExpiryInvalid = errors.New("expiry must be one of 1h, 1d, 7d, 30d")
)

var (
	aliasRe  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	reserved = map[string]struct{}{
		"api": {}, "healthz": {}, "health": {}, "static": {},
		"admin": {}, "assets": {}, "favicon.ico": {}, "robots.txt": {},
	}
	expiryOptions = map[string]time.Duration{
		"1h":  time.Hour,
		"1d":  24 * time.Hour,
		"7d":  7 * 24 * time.Hour,
		"30d": 30 * 24 * time.Hour,
	}
)

// URL validates and normalizes a target URL.
func URL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrURLEmpty
	}
	if len(raw) > MaxURLLength {
		return "", ErrURLTooLong
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", ErrURLInvalid
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrURLScheme
	}
	if u.Host == "" {
		return "", ErrURLInvalid
	}
	return u.String(), nil
}

// Alias validates a user-supplied short alias.
func Alias(raw string) (string, error) {
	alias := strings.TrimSpace(raw)
	if len(alias) < MinAliasLength || len(alias) > MaxAliasLength {
		return "", ErrAliasLength
	}
	if !aliasRe.MatchString(alias) {
		return "", ErrAliasFormat
	}
	if _, ok := reserved[strings.ToLower(alias)]; ok {
		return "", ErrAliasReserved
	}
	return alias, nil
}

// ExpiresIn parses an expiry option. Empty string means no expiry.
func ExpiresIn(raw string) (time.Duration, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	d, ok := expiryOptions[strings.ToLower(raw)]
	if !ok {
		return 0, ErrExpiryInvalid
	}
	return d, nil
}
