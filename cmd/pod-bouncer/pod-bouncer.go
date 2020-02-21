package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
	"github.com/ewohltman/pod-bouncer/internal/pkg/server"
)

const (
	port = "8080"

	contextTimeout = 5 * time.Second
)

func newKubeClientset() (*kubernetes.Clientset, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(inClusterConfig)
}

func startServer(log *logging.Logger, kubeClientset kubernetes.Interface) (httpServer *server.Instance, sigTerm chan os.Signal) {
	httpServer = server.New(log, port, kubeClientset)

	sigTerm = make(chan os.Signal, 1)

	signal.Notify(sigTerm, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.WithError(err).Error("HTTP server error")
				sigTerm <- syscall.SIGTERM
			}
		}
	}()

	return
}

func main() {
	log := logging.New()

	log.Info("pod-bouncer starting up")

	kubeClientset, err := newKubeClientset()
	if err != nil {
		log.WithError(err).Error("Error creating new Kubernetes Clientset")
	}

	httpServer, sigTerm := startServer(log, kubeClientset)

	<-sigTerm // Block until the SIGTERM OS signal

	ctx, cancelFunc := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelFunc()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Error("Error shutting down HTTP server gracefully")
	}
}
