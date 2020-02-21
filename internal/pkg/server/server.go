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
	errorEmptyRequestBody       = errorInternalServerError + ": empty request body"
	errorReadingRequestBody     = errorInternalServerError + ": unable to read request body"
	errorClosingRequestBody     = errorInternalServerError + ": unable to close request body"
	errorWritingResponseBody    = errorInternalServerError + ": unable to write response body"
	errorUnmarshalingAlertEvent = errorInternalServerError + ": unable to unmarshal Event"
)

// Instance wraps an *http.Server for extending custom functionality.
type Instance struct {
	*http.Server
}

// New returns a new pre-configured server instance.
func New(log logging.Interface, port string) *Instance {
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

	return &Instance{
		Server: &http.Server{
			Addr:     "0.0.0.0:" + port,
			Handler:  mux,
			ErrorLog: errorLog,
		},
	}
}

func alertHandler(log logging.Interface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Event received")

		if r.Body == nil {
			log.Error(errorEmptyRequestBody)

			send400ErrorResponse(log, w, errorEmptyRequestBody)

			return
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error(errorReadingRequestBody)

			send500ErrorResponse(log, w, errorReadingRequestBody)

			return
		}

		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				log.WithError(closeErr).Error(errorClosingRequestBody)
			}
		}()

		event, err := alert.NewEvent(reqBody)
		if err != nil {
			log.WithError(err).Error(errorUnmarshalingAlertEvent)

			send500ErrorResponse(log, w, errorUnmarshalingAlertEvent)

			return
		}

		log.WithField("event", event).Info("Parsed event")
	}
}

func rootHandler(log logging.Interface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			return
		}

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

func send400ErrorResponse(log logging.Interface, w http.ResponseWriter, message string) {
	sendErrorResponse(log, w, message, http.StatusBadRequest)
}

func send500ErrorResponse(log logging.Interface, w http.ResponseWriter, message string) {
	sendErrorResponse(log, w, message, http.StatusInternalServerError)
}

func sendErrorResponse(log logging.Interface, w http.ResponseWriter, message string, respCode int) {
	w.WriteHeader(respCode)

	_, err := w.Write([]byte(message))
	if err != nil {
		log.WithError(err).Error(errorWritingResponseBody)
	}
}
