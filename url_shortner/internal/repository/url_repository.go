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
	// Check whether the alias stored in PostgreSQL exactly matches
	// the alias received from the request.
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			alias,
			alias = $1 AS exact_match,
			length(alias) AS alias_length,
			length($1) AS input_length,
			encode(convert_to(alias, 'UTF8'), 'hex') AS alias_hex,
			encode(convert_to($1, 'UTF8'), 'hex') AS input_hex
		FROM urls
		`,
		alias,
	)

	if err != nil {
		return nil, fmt.Errorf("diagnostic alias comparison failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			existingAlias string
			exactMatch    bool
			aliasLength   int
			inputLength   int
			aliasHex      string
			inputHex      string
		)

		if err := rows.Scan(
			&existingAlias,
			&exactMatch,
			&aliasLength,
			&inputLength,
			&aliasHex,
			&inputHex,
		); err != nil {
			return nil, fmt.Errorf("diagnostic alias comparison scan failed: %w", err)
		}

		fmt.Printf(
			"DIAGNOSTIC: database alias=%q exact_match=%t alias_length=%d input_length=%d alias_hex=%s input_hex=%s\n",
			existingAlias,
			exactMatch,
			aliasLength,
			inputLength,
			aliasHex,
			inputHex,
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
