package links

import (
	"context"
	"errors"
	"fmt"
	"time"

	"trimora/internal/shortcode"
	"trimora/internal/validate"
)

// ErrAliasUnavailable is returned when a user-supplied alias is taken.
var ErrAliasUnavailable = errors.New("alias unavailable")

// ErrExpired is returned when a resolved link has passed its expiry.
var ErrExpired = errors.New("link expired")

const (
	codeLength       = 7
	createMaxRetries = 5
)

type Service struct {
	repo *Repository
	now  func() time.Time
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Create stores a new link, generating a code when alias is empty.
// expiresIn of zero means the link never expires.
func (s *Service) Create(ctx context.Context, rawURL, alias string, expiresIn time.Duration) (Link, error) {
	target, err := validate.URL(rawURL)
	if err != nil {
		return Link{}, err
	}

	var expiresAt *time.Time
	if expiresIn > 0 {
		t := s.now().Add(expiresIn).UTC()
		expiresAt = &t
	}

	if alias != "" {
		code, err := validate.Alias(alias)
		if err != nil {
			return Link{}, err
		}
		link, err := s.repo.Create(ctx, code, target, expiresAt)
		if errors.Is(err, ErrAliasTaken) {
			return Link{}, ErrAliasUnavailable
		}
		if err != nil {
			return Link{}, err
		}
		return link, nil
	}

	if expiresAt == nil {
		if existing, err := s.repo.FindReusableByTargetURL(ctx, target); err == nil {
			return existing, nil
		} else if !errors.Is(err, ErrNotFound) {
			return Link{}, err
		}
	}

	for i := 0; i < createMaxRetries; i++ {
		code, err := shortcode.Generate(codeLength)
		if err != nil {
			return Link{}, err
		}
		link, err := s.repo.Create(ctx, code, target, expiresAt)
		if errors.Is(err, ErrAliasTaken) {
			continue
		}
		if err != nil {
			return Link{}, err
		}
		return link, nil
	}
	return Link{}, fmt.Errorf("could not allocate short code after %d attempts", createMaxRetries)
}

// Resolve returns the target URL for a short code, or ErrExpired if past expiry.
func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return "", err
	}
	if link.ExpiresAt != nil && !s.now().Before(*link.ExpiresAt) {
		return "", ErrExpired
	}
	return link.TargetURL, nil
}
