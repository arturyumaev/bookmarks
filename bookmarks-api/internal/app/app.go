package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
	"github.com/gin-gonic/gin"

	// bookmark
	bookmarkHTTP "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/delivery/http"
	bookmarkRepo "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/repository/boltdb"
	bookmarkUC "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/usecase"
)

type IApplication interface {
	Run() error
}

type application struct {
	httpServer *http.Server
	config     *models.Config
}

func (app *application) Run() error {
	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return app.httpServer.Shutdown(ctx)
}

func NewApplication(config *models.Config) IApplication {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// bookmark
	bookmarkRepo := bookmarkRepo.NewRepository()
	bookmarkUseCase := bookmarkUC.NewUseCase(bookmarkRepo)
	bookmarkHTTP.RegisterHTTPEndpoints(r, bookmarkUseCase)

	httpServer := &http.Server{
		Addr:           ":" + config.Server.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &application{
		httpServer,
		config,
	}
}
