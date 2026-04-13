package repository

import (
	"shortener/models"
	"sync"
)

type ShortLinkRepo struct {
	mu      sync.RWMutex
	byID    map[string]*models.ShortLink
	byURL   map[string]*models.ShortLink // for dedup
}

func NewShortLinkRepo() *ShortLinkRepo {
	return &ShortLinkRepo{
		byID:  make(map[string]*models.ShortLink),
		byURL: make(map[string]*models.ShortLink),
	}
}

func (r *ShortLinkRepo) Save(link *models.ShortLink) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[link.ID] = link
	r.byURL[link.OriginalURL] = link
}

func (r *ShortLinkRepo) FindByID(id string) (*models.ShortLink, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.byID[id]
	return link, ok
}

func (r *ShortLinkRepo) FindByURL(originalURL string) (*models.ShortLink, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.byURL[originalURL]
	return link, ok
}
