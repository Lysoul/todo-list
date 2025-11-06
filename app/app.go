package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lysoul/gocommon/ginserver"
	"github.com/Lysoul/gocommon/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

//nolint:gochecknoglobals // we need this for versioning
var Version = "unknown"

type Config struct {
	HTTP ginserver.Config
	// Postgres postgres.Config // will use it later

	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"20s"`
}

func Start() error {
	log := monitoring.Logger()

	var config Config
	envconfig.MustProcess("", &config)

	router, httpStart := ginserver.InitGin(config.HTTP, log)
	basePath := config.HTTP.Prefix
	router.GET(basePath+"/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": Version,
		})
	})

	apiGroup := router.Group(basePath)

	apiGroup.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	_, httpStop := httpStart()
	monitoring.ServeTelemetry(3030)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	httpStop(ctx)

	return nil
}
