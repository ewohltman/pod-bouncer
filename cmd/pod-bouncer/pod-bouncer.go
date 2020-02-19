package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
	"github.com/ewohltman/pod-bouncer/internal/pkg/server"
)

const (
	port = "8080"

	contextTimeout = 5 * time.Second
)

func main() {
	log := logging.New()

	log.Info("pod-bouncer starting up")

	httpServer := server.New(log, port)
	sigTerm := make(chan os.Signal, 1)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.WithError(err).Error("HTTP server error")
				sigTerm <- syscall.SIGTERM
			}
		}
	}()

	signal.Notify(sigTerm, syscall.SIGTERM)

	<-sigTerm // Block until the SIGTERM OS signal

	ctx, cancelFunc := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelFunc()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Error("Error shutting down HTTP server gracefully")
	}
}
