package repository

import (
	"context"

	"url_shortner/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) CreateURL(ctx context.Context, url *model.URL) error {
	query := `
		INSERT INTO urls (alias, original_url)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		url.Alias,
		url.OriginalURL,
	).Scan(
		&url.ID,
		&url.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *URLRepository) GetURLByAlias(ctx context.Context, alias string) (*model.URL, error) {
	query := `
		SELECT id, alias, original_url, created_at
		FROM urls
		WHERE alias = $1
	`

	url := &model.URL{}

	err := r.db.QueryRow(
		ctx,
		query,
		alias,
	).Scan(
		&url.ID,
		&url.Alias,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return url, nil
}
