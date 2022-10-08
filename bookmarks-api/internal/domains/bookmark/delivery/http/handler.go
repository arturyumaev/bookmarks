package http

import (
	"net/http"

	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/gin-gonic/gin"
)

type handler struct {
	uc bookmark.UseCase
}

func (h *handler) CreateBookmark(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *handler) GetBookmark(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *handler) GetBookmarks(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *handler) UpdateBookmark(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *handler) DeleteBookmark(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func NewHandler(useCase bookmark.UseCase) bookmark.Handler {
	return &handler{useCase}
}
