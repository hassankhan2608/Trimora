package links

import (
	"context"
	"errors"
	"fmt"

	"trimora/internal/shortcode"
	"trimora/internal/validate"
)

// ErrAliasUnavailable is returned when a user-supplied alias is taken.
var ErrAliasUnavailable = errors.New("alias unavailable")

const (
	codeLength       = 7
	createMaxRetries = 5
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create stores a new link, generating a code when alias is empty.
func (s *Service) Create(ctx context.Context, rawURL, alias string) (Link, error) {
	target, err := validate.URL(rawURL)
	if err != nil {
		return Link{}, err
	}

	if alias != "" {
		code, err := validate.Alias(alias)
		if err != nil {
			return Link{}, err
		}
		link, err := s.repo.Create(ctx, code, target)
		if errors.Is(err, ErrAliasTaken) {
			return Link{}, ErrAliasUnavailable
		}
		if err != nil {
			return Link{}, err
		}
		return link, nil
	}

	for i := 0; i < createMaxRetries; i++ {
		code, err := shortcode.Generate(codeLength)
		if err != nil {
			return Link{}, err
		}
		link, err := s.repo.Create(ctx, code, target)
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

// Resolve returns the target URL for a short code.
func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return "", err
	}
	return link.TargetURL, nil
}
