package http

import (
	"log"
	"net/http"

	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
	"github.com/gin-gonic/gin"
)

type handler struct {
	uc bookmark.UseCase
}

func (h *handler) CreateBookmark(c *gin.Context) {
	bookmark := &models.Bookmark{}
	if err := c.BindJSON(bookmark); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	bm, err := h.uc.CreateBookmark(c.Request.Context(), bookmark)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, bm)
}

func (h *handler) GetBookmark(c *gin.Context) {
	id := c.Param("id")
	bm, err := h.uc.GetBookmark(c.Request.Context(), id)
	if err != nil {
		log.Println(err.Error())

		if err == bookmark.ErrNoSuchBookmark {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bm)
}

func (h *handler) GetBookmarks(c *gin.Context) {
	bms, err := h.uc.GetBookmarks(c.Request.Context())
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bms)
}

func (h *handler) UpdateBookmark(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	bm, err := h.uc.GetBookmark(ctx, id)
	if err != nil {
		log.Println(err.Error())

		if err == bookmark.ErrNoSuchBookmark {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	bm.Id = id
	nextBm, err := h.uc.UpdateBookmark(ctx, id, bm)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nextBm)
}

func (h *handler) DeleteBookmark(c *gin.Context) {
	id := c.Param("id")
	err := h.uc.DeleteBookmark(c.Request.Context(), id)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func NewHandler(useCase bookmark.UseCase) bookmark.Handler {
	return &handler{useCase}
}
