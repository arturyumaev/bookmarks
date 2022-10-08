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
	bolt "go.etcd.io/bbolt"

	// bookmark
	bookmarkHTTP "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/delivery/http"
	bookmarkRepo "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/repository/boltdb"
	bookmarkUC "github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark/usecase"
)

type Application struct {
	HttpServer *http.Server
	Config     *models.Config
	BoltDB     *bolt.DB
}

func (app *Application) Run() error {
	go func() {
		if err := app.HttpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return app.HttpServer.Shutdown(ctx)
}

func NewApplication(config *models.Config) *Application {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	boltDB, err := initBoltDB(config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// bookmark
	bookmarkRepo := bookmarkRepo.NewRepository(boltDB)
	bookmarkUseCase := bookmarkUC.NewUseCase(bookmarkRepo)
	bookmarkHTTP.RegisterHTTPEndpoints(r, bookmarkUseCase)

	httpServer := &http.Server{
		Addr:           ":" + config.Server.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Application{
		httpServer,
		config,
		boltDB,
	}
}

func initBoltDB(config *models.Config) (*bolt.DB, error) {
	db, err := bolt.Open(config.DB.BoltDB, 0600, nil)
	return db, err
}
