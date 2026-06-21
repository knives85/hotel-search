// Command hotel-search starts the HTTP server for the app module.
//
// This is the Go reimplementation of the `app` module of Hotel Search.
// At this stage it only wires configuration and the HTTP routes; the routes
// return 501 Not Implemented until the corresponding adapters are ported.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/knives85/hotel-search/internal/adapter/web"
	"github.com/knives85/hotel-search/internal/config"
)

func main() {
	cfg := config.Load()

	srv := web.NewServer(cfg.ContextPath)
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (context path %q)", cfg.Addr, cfg.ContextPath)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
