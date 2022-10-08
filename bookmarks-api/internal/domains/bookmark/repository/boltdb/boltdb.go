package boltdb

import (
	"context"

	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/arturyumaev/bookmarks/bookmarks-api/models"

	bolt "go.etcd.io/bbolt"
)

type repository struct {
	boltdb *bolt.DB
}

func (repo *repository) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) (*models.Bookmark, error) {
	return nil, nil
}

func (repo *repository) GetBookmark(ctx context.Context, bookmarkId string) (*models.Bookmark, error) {
	return nil, nil
}

func (repo *repository) GetBookmarks(ctx context.Context) ([]*models.Bookmark, error) {
	return nil, nil
}

func (repo *repository) UpdateBookmark(ctx context.Context, bookmarkId string, bookmark *models.Bookmark) (*models.Bookmark, error) {
	return nil, nil
}

func (repo *repository) DeleteBookmark(ctx context.Context, bookmarkId string) error {
	return nil
}

func NewRepository(boltdb *bolt.DB) bookmark.Repository {
	return &repository{}
}
