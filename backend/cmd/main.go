package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/app"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
)

func listenHTTPServer(r *gin.Engine, port int) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		log.Printf("Starting HTTP server on port %d\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", err)
		}
	}()

	return server
}

func main() {
	cfg := config.MustLoadConfig()

	r := gin.Default()

	close := app.MustConfigureApp(r, cfg)
	defer close()

	server := listenHTTPServer(r, cfg.RestConfig.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown", err)
	}

	log.Println("Server stopped")
}
