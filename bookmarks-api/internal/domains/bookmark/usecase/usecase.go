package usecase

import (
	"context"

	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
)

type useCase struct {
	repo bookmark.Repository
}

func (uc *useCase) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) (*models.Bookmark, error) {
	bookmark, err := uc.repo.CreateBookmark(ctx, bookmark)

	return bookmark, err
}

func (uc *useCase) GetBookmark(ctx context.Context, bookmarkId string) (*models.Bookmark, error) {
	bookmark, err := uc.repo.GetBookmark(ctx, bookmarkId)

	return bookmark, err
}

func (uc *useCase) GetBookmarks(ctx context.Context) ([]*models.Bookmark, error) {
	bookmarks, err := uc.repo.GetBookmarks(ctx)

	return bookmarks, err
}

func (uc *useCase) UpdateBookmark(ctx context.Context, bookmarkId string, bookmark *models.Bookmark) (*models.Bookmark, error) {
	bookmarks, err := uc.repo.UpdateBookmark(ctx, bookmarkId, bookmark)

	return bookmarks, err
}

func (uc *useCase) DeleteBookmark(ctx context.Context, bookmarkId string) error {
	err := uc.repo.DeleteBookmark(ctx, bookmarkId)

	return err
}

func NewUseCase(repo bookmark.Repository) bookmark.UseCase {
	return &useCase{repo}
}
