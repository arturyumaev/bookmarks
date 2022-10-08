package http

import (
	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc bookmark.UseCase) {
	h := NewHandler(uc)

	userEndpoints := router.Group("/bookmarks")

	{
		userEndpoints.POST("/", h.CreateBookmark)
		userEndpoints.GET("/", h.GetBookmarks)
		userEndpoints.GET("/:id", h.GetBookmark)
		userEndpoints.PUT("/:id", h.UpdateBookmark)
		userEndpoints.DELETE("/:id", h.DeleteBookmark)
	}
}
