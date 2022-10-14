package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/arturyumaev/bookmarks/bookmarks-api/metrics"
	"github.com/arturyumaev/bookmarks/bookmarks-api/middleware"
	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"

	// bookmark
	"github.com/arturyumaev/bookmarks/bookmarks-api/internal/domains/bookmark"
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
	defer app.BoltDB.Close()

	return app.HttpServer.Shutdown(ctx)
}

func NewApplication(config *models.Config) *Application {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	metrics.InitMetrics()

	boltDB, err := initBoltDB(config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(middleware.PrometheusMiddleware)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
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
	if err != nil {
		return nil, err
	}

	err = initBoltDBBucket(db, bookmark.BookmarkBucketName, !config.IsProduction())
	if err != nil {
		return nil, err
	}

	err = initBoltDBBucket(db, bookmark.ColorBucketName, !config.IsProduction())
	if err != nil {
		return nil, err
	}

	return db, err
}

func initBoltDBBucket(db *bolt.DB, name string, isDevelopment bool) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if isDevelopment {
		log.Println("successfully created bucket", name)
	}

	return nil
}
