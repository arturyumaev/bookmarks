package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/arturyumaev/bookmarks/bookmarks-api/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
			logrus.Fatalf("failed to listen and serve: %+v", err)
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

	logrus.SetOutput(os.Stdout)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	httpServer := &http.Server{
		Addr:           ":" + config.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &application{
		httpServer,
		config,
	}
}
