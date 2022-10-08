package bookmark

import (
	"context"

	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
)

type UseCase interface {
	CreateBookmark(ctx context.Context, bookmark *models.Bookmark) (*models.Bookmark, error)
	GetBookmark(ctx context.Context, bookmarkId string) (*models.Bookmark, error)
	GetBookmarks(ctx context.Context) ([]*models.Bookmark, error)
	UpdateBookmark(ctx context.Context, bookmarkId string, bookmark *models.Bookmark) (*models.Bookmark, error)
	DeleteBookmark(ctx context.Context, bookmarkId string) error
}
