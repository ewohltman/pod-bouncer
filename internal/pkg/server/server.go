// Package server provides HTTP server functionality.
package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"

	stdLog "log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

const (
	metricsEndpoint = "/metrics"

	pprofIndexEndpoint   = "/debug/pprof/"
	pprofCmdlineEndpoint = "/debug/pprof/cmdline"
	pprofProfileEndpoint = "/debug/pprof/profile"
	pprofSymbolEndpoint  = "/debug/pprof/symbol"
	pprofTraceEndpoint   = "/debug/pprof/trace"

	alertEndpoint = "/alert"
	rootEndpoint  = "/"
)

// New returns a new pre-configured server instance.
func New(log logging.Interface, port string) *http.Server {
	mux := http.NewServeMux()

	mux.Handle(metricsEndpoint, promhttp.Handler())

	mux.HandleFunc(pprofIndexEndpoint, pprof.Index)
	mux.HandleFunc(pprofCmdlineEndpoint, pprof.Cmdline)
	mux.HandleFunc(pprofProfileEndpoint, pprof.Profile)
	mux.HandleFunc(pprofSymbolEndpoint, pprof.Symbol)
	mux.HandleFunc(pprofTraceEndpoint, pprof.Trace)

	mux.HandleFunc(alertEndpoint, alertHandler(log))
	mux.HandleFunc(rootEndpoint, rootHandler(log))

	errorLog := stdLog.New(log.WrappedLogger().WriterLevel(logrus.ErrorLevel), "", 0)

	return &http.Server{
		Addr:     "0.0.0.0:" + port,
		Handler:  mux,
		ErrorLog: errorLog,
	}
}

func alertHandler(log logging.Interface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Warn("Internal HTTP server error reading request body")
		}

		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				log.WithError(closeErr).Warn("Internal HTTP server error closing request body")
			}
		}()

		log.WithField("request", string(reqBody)).
			Infof("Request received on %s", alertEndpoint)
	}
}

func rootHandler(log logging.Interface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			log.WithError(err).Warn("Internal HTTP server error draining request body")
		}

		err = r.Body.Close()
		if err != nil {
			log.WithError(err).Warn("Internal HTTP server error closing request body")
		}
	}
}
