package links

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrAliasTaken is returned when a code already exists.
var ErrAliasTaken = errors.New("alias already taken")

// ErrNotFound is returned when no link matches a code.
var ErrNotFound = errors.New("link not found")

type Link struct {
	Code      string
	TargetURL string
	CreatedAt time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, code, targetURL string) (Link, error) {
	const q = `
		INSERT INTO links (code, target_url)
		VALUES ($1, $2)
		RETURNING code, target_url, created_at
	`
	var link Link
	err := r.db.QueryRowContext(ctx, q, code, targetURL).
		Scan(&link.Code, &link.TargetURL, &link.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Link{}, ErrAliasTaken
		}
		return Link{}, err
	}
	return link, nil
}

func (r *Repository) FindByCode(ctx context.Context, code string) (Link, error) {
	const q = `SELECT code, target_url, created_at FROM links WHERE code = $1`
	var link Link
	err := r.db.QueryRowContext(ctx, q, code).
		Scan(&link.Code, &link.TargetURL, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Link{}, ErrNotFound
	}
	if err != nil {
		return Link{}, err
	}
	return link, nil
}

// FindByTargetURL returns the oldest link for a target URL, if any.
// Used to deduplicate when a user shortens the same URL without a custom alias.
func (r *Repository) FindByTargetURL(ctx context.Context, targetURL string) (Link, error) {
	const q = `
		SELECT code, target_url, created_at
		FROM links
		WHERE target_url = $1
		ORDER BY created_at ASC
		LIMIT 1
	`
	var link Link
	err := r.db.QueryRowContext(ctx, q, targetURL).
		Scan(&link.Code, &link.TargetURL, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Link{}, ErrNotFound
	}
	if err != nil {
		return Link{}, err
	}
	return link, nil
}
