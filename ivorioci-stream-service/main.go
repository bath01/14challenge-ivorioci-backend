package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ivorioci-stream-service/config"
	"ivorioci-stream-service/handlers"
	"ivorioci-stream-service/routes"
	"ivorioci-stream-service/services"
)

func main() {
	cfg := config.Load()

	db, err := config.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[FATAL] %v", err)
	}
	defer db.Close()

	// Ensure video storage directory exists
	if err := os.MkdirAll(cfg.VideoStoragePath, 0o750); err != nil {
		log.Fatalf("[FATAL] Cannot create video storage path %q: %v", cfg.VideoStoragePath, err)
	}

	// Wire dependencies
	videoSvc := services.NewVideoService(db)
	categorySvc := services.NewCategoryService(db)

	videoH := handlers.NewVideoHandler(videoSvc, categorySvc)
	categoryH := handlers.NewCategoryHandler(categorySvc)
	streamH := handlers.NewStreamHandler(videoSvc, cfg.VideoStoragePath)

	router := routes.New(videoH, categoryH, streamH, cfg.JWTAccessSecret)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // disabled for streaming responses
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("[SERVER] ivorioci-stream-service running on :%s (env=%s)", cfg.Port, cfg.GoEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] %v", err)
		}
	}()

	<-stop
	log.Println("[SERVER] Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[SERVER] Shutdown error: %v", err)
	}
	log.Println("[SERVER] Stopped")
}
