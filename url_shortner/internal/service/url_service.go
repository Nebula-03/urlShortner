package service

import (
	"context"
	"crypto/rand"
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

	if alias == "" {
		var err error

		for {
			alias, err = generateAlias()
			if err != nil {
				return nil, err
			}

			_, err = s.repository.GetURLByAlias(ctx, alias)

			if errors.Is(err, pgx.ErrNoRows) {
				break
			}

			if err != nil {
				return nil, err
			}
		}
	} else {
		_, err := s.repository.GetURLByAlias(ctx, alias)

		if err == nil {
			return nil, ErrAliasAlreadyExists
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	url := &model.URL{
		Alias:       alias,
		OriginalURL: originalURL,
	}

	err := s.repository.CreateURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func generateAlias() (string, error) {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, 6)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	alias := make([]byte, 6)

	for i := range bytes {
		alias[i] = characters[int(bytes[i])%len(characters)]
	}

	return string(alias), nil
}

func (s *URLService) GetURLByAlias(
	ctx context.Context,
	alias string,
) (*model.URL, error) {

	return s.repository.GetURLByAlias(ctx, alias)
}
