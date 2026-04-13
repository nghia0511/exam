package service

import (
	"errors"
	"math/rand"
	"net/url"
	"shortener/models"
	"shortener/repository"
	"time"
)

const baseURL = "http://localhost:8080"

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type ShortLinkService struct {
	repo *repository.ShortLinkRepo
}

func NewShortLinkService(repo *repository.ShortLinkRepo) *ShortLinkService {
	return &ShortLinkService{repo: repo}
}

func (s *ShortLinkService) Create(originalURL string) (*models.ShortLink, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	// dedup: same URL gets same short code
	if existing, ok := s.repo.FindByURL(originalURL); ok {
		return existing, nil
	}

	id := generateCode(7)
	// retry on the off-chance of collision (extremely rare but still)
	for i := 0; i < 5; i++ {
		if _, exists := s.repo.FindByID(id); !exists {
			break
		}
		id = generateCode(7)
	}
	link := &models.ShortLink{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    baseURL + "/shortlinks/" + id,
		CreatedAt:   time.Now(),
	}
	s.repo.Save(link)
	return link, nil
}

func (s *ShortLinkService) GetByID(id string) (*models.ShortLink, error) {
	link, ok := s.repo.FindByID(id)
	if !ok {
		return nil, errors.New("not found")
	}
	return link, nil
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("invalid url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url must start with http or https")
	}
	if u.Host == "" {
		return errors.New("url missing host")
	}
	return nil
}

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
