package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"url_shortner/internal/model"
	"url_shortner/internal/repository"
)

var ErrAliasAlreadyExists = errors.New("alias already exists")

type URLService struct {
	repository *repository.URLRepository
}

func NewURLService(repository *repository.URLRepository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) CreateURL(
	ctx context.Context,
	alias string,
	originalURL string,
) (*model.URL, error) {

	_, err := s.repository.GetURLByAlias(ctx, alias)

	if err == nil {
		return nil, ErrAliasAlreadyExists
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	url := &model.URL{
		Alias:       alias,
		OriginalURL: originalURL,
	}

	err = s.repository.CreateURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}
