package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/api"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/image"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	ll2datasyncer "github.com/vamosdalian/launchdate-backend/internal/service/ll2_data_syncer"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Env == config.ENV_PRODUCTION {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	db, cleandb, err := db.NewMongoDB(cfg.MongodbURL, cfg.MongodbDatabase)
	if err != nil {
		logrus.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer cleandb()
	logrus.Infof("create mongodb database: %s", cfg.MongodbDatabase)

	hc := util.NewHTTPClient()
	coreservice := core.NewMainService(db)
	if err := coreservice.EnsurePageBackgroundIndexes(); err != nil {
		logrus.Fatalf("failed to ensure page background indexes: %v", err)
	}
	ll2server := ll2.NewLL2Service(cfg, db, hc)
	ll2syncer := ll2datasyncer.NewLL2DataSyncer(cfg, ll2server, coreservice)
	if err := ll2syncer.RestoreTasks(); err != nil {
		logrus.Errorf("failed to restore sync tasks: %v", err)
	}
	s3Client, err := util.CreateS3Client(cfg.ImageConf.AccessKey, cfg.ImageConf.SecretKey,
		cfg.ImageConf.Region, cfg.ImageConf.Endpoint)
	if err != nil {
		logrus.Fatalf("failed to create s3 client: %v", err)
	}
	imageService := image.NewImageService(s3Client, db, cfg.ImageConf.Bucket, cfg.ImageDomain)

	handler := api.NewHandler(logrus.StandardLogger(), cfg, db, ll2server, ll2syncer, coreservice, imageService)
	router := api.SetupRouter(handler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("server forced to shutdown: %v", err)
	}

	logrus.Info("server stopped")
}
