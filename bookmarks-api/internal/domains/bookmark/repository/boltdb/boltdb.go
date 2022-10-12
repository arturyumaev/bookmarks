package boltdb

import (
	"context"
	"encoding/json"

	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/arturyumaev/bookmarks/bookmarks-api/models"

	bolt "go.etcd.io/bbolt"
)

type repository struct {
	boltdb *bolt.DB
}

func (repo *repository) CreateBookmark(ctx context.Context, bm *models.Bookmark) (*models.Bookmark, error) {
	buf, err := repo.marshall(bm)
	if err != nil {
		return nil, err
	}

	err = repo.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bookmark.BookmarkBucketName))
		err := b.Put([]byte(bm.Id), buf)
		return err
	})

	return bm, nil
}

func (repo *repository) GetBookmark(ctx context.Context, bmId string) (*models.Bookmark, error) {
	bm := &models.Bookmark{}

	err := repo.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bookmark.BookmarkBucketName))
		bytes := b.Get([]byte(bmId))

		var err error
		bm, err = repo.unmarshall(bytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return bm, nil
}

func (repo *repository) GetBookmarks(ctx context.Context) ([]*models.Bookmark, error) {
	var bms = []*models.Bookmark{}

	if err := repo.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bookmark.BookmarkBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			bm, err := repo.unmarshall(v)
			if err != nil {
				return err
			}
			bms = append(bms, bm)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return bms, nil
}

func (repo *repository) UpdateBookmark(ctx context.Context, bmId string, bm *models.Bookmark) (*models.Bookmark, error) {
	return nil, nil
}

func (repo *repository) DeleteBookmark(ctx context.Context, bmId string) error {
	return nil
}

func (repo *repository) marshall(bm *models.Bookmark) ([]byte, error) {
	buf, err := json.Marshal(bm)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (repo *repository) unmarshall(bytes []byte) (*models.Bookmark, error) {
	bm := &models.Bookmark{}
	err := json.Unmarshal(bytes, bm)
	if err != nil {
		return nil, err
	}
	return bm, nil
}

func NewRepository(boltdb *bolt.DB) bookmark.Repository {
	return &repository{boltdb}
}
