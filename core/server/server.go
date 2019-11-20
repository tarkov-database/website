package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tarkov-database/website/route"

	"github.com/google/logger"
)

// Start starts the HTTP server
func Start() {
	mux := route.Load()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		fmt.Println()
		logger.Info("HTTP server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if cfg.TLS {
		logger.Infof("HTTPS server listen and serve on *:%v\n\n", cfg.Port)
		if err := srv.ListenAndServeTLS(cfg.Certificate, cfg.PrivateKey); err != http.ErrServerClosed {
			logger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		logger.Infof("HTTP server listen and serve on *:%v\n\n", cfg.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}

	<-idleConnsClosed
}
