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

	"github.com/ewohltman/pod-bouncer/internal/pkg/alert"
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

const (
	errorInternalServerError    = "Internal server error"
	errorUnmarshalingAlertEvent = errorInternalServerError + " unmarshaling Event"
	errorReadingRequestBody     = errorInternalServerError + " reading request body"
	errorClosingRequestBody     = errorInternalServerError + " closing request body"
	errorWritingResponseBody    = errorInternalServerError + " writing response body"
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
		log.Info("Event received")

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Warn(errorReadingRequestBody)

			sendErrorResponse(log, w, errorReadingRequestBody)

			return
		}

		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				log.WithError(closeErr).Warn(errorClosingRequestBody)
			}
		}()

		event, err := alert.NewEvent(reqBody)
		if err != nil {
			log.WithError(err).Error(errorUnmarshalingAlertEvent)

			sendErrorResponse(log, w, errorUnmarshalingAlertEvent)

			return
		}

		log.WithField("event", event).Info("Parsed event")
	}
}

func rootHandler(log logging.Interface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			log.WithError(err).Warn(errorReadingRequestBody)
		}

		err = r.Body.Close()
		if err != nil {
			log.WithError(err).Warn(errorClosingRequestBody)
		}
	}
}

func sendErrorResponse(log logging.Interface, w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)

	_, err := w.Write([]byte(message))
	if err != nil {
		log.WithError(err).Error(errorWritingResponseBody)
	}
}
