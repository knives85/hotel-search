// Command hotel-search starts the HTTP server for the app module.
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

	oslib "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/knives85/hotel-search/internal/adapter/opensearch"
	"github.com/knives85/hotel-search/internal/adapter/web"
	"github.com/knives85/hotel-search/internal/config"
)

// The OpenSearch index that holds hotel documents. Hardcoded for now; promote
// to config when more than one environment uses a non-default index name.
const hotelsIndex = "hotels"

func main() {
	cfg := config.Load()

	osClient, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: oslib.Config{
			Addresses: []string{cfg.OpenSearchEndpoint},
			// TODO: configure the AWS SigV4 signer when targeting Amazon OpenSearch
			// — use github.com/opensearch-project/opensearch-go/v4/signer/awsv2.
		},
	})
	if err != nil {
		log.Fatalf("opensearch client init failed: %v", err)
	}
	searchRepo := opensearch.NewRepository(osClient, hotelsIndex)

	srv := web.NewServer(cfg.ContextPath, web.Deps{Search: searchRepo})
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (context path %q, opensearch %s)",
			cfg.Addr, cfg.ContextPath, cfg.OpenSearchEndpoint)
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
