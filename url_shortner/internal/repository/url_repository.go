package repository

import (
	"context"
	"fmt"

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

	// Temporary diagnostic check:
	// Find out exactly which database/schema the deployed app is using.
	var databaseName string
	var schemaName string

	err := r.db.QueryRow(
		ctx,
		"SELECT current_database(), current_schema()",
	).Scan(&databaseName, &schemaName)

	if err != nil {
		return nil, fmt.Errorf("diagnostic database check failed: %w", err)
	}

	fmt.Printf(
		"DIAGNOSTIC: database=%s schema=%s alias=%q\n",
		databaseName,
		schemaName,
		alias,
	)

	// Temporary diagnostic check:
	// Show the aliases visible to the running application.
	rows, err := r.db.Query(
		ctx,
		"SELECT alias FROM urls ORDER BY id",
	)

	if err != nil {
		return nil, fmt.Errorf("diagnostic alias query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existingAlias string

		if err := rows.Scan(&existingAlias); err != nil {
			return nil, fmt.Errorf("diagnostic alias scan failed: %w", err)
		}

		fmt.Printf(
			"DIAGNOSTIC: database alias=%q\n",
			existingAlias,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("diagnostic rows error: %w", err)
	}

	query := `
		SELECT id, alias, original_url, created_at
		FROM urls
		WHERE alias = $1
	`

	url := &model.URL{}

	err = r.db.QueryRow(
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
