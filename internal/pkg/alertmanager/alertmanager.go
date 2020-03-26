// Package alertmanager contains type definitions and functionality for
// handling events from AlertManager.
package alertmanager

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

const (
	statusFiring     = "firing"
	severityCritical = "critical"

	errorDeletingPod = "error deleting pod"
)

// Handler is an interface for abstracting handling implementations.
type Handler interface {
	Handle(context.Context, []byte) error
}

// Annotations are extra metadata of an Alert.
type Annotations struct {
	Description string `json:"description,omitempty"`
}

// Labels are metadata of an Alert.
type Labels struct {
	AlertName  string `json:"alertname,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	Instance   string `json:"instance,omitempty"`
	Job        string `json:"job,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Pod        string `json:"pod,omitempty"`
	Prometheus string `json:"prometheus,omitempty"`
	Service    string `json:"service,omitempty"`
	Severity   string `json:"severity,omitempty"`
}

// Alert is an alert from AlertManager within an Event.
type Alert struct {
	Status       string       `json:"status,omitempty"`
	Labels       *Labels      `json:"labels,omitempty"`
	Annotations  *Annotations `json:"annotations,omitempty"`
	StartsAt     time.Time    `json:"startsAt,omitempty"`
	EndsAt       time.Time    `json:"endsAt,omitempty"`
	GeneratorURL string       `json:"generatorURL,omitempty"`
	Fingerprint  string       `json:"fingerprint,omitempty"`
}

// Event is the payload of an AlertManager webhook notification.
type Event struct {
	Receiver          string       `json:"receiver,omitempty"`
	Status            string       `json:"status,omitempty"`
	Alerts            []*Alert     `json:"alerts,omitempty"`
	GroupLabels       *Labels      `json:"groupLabels,omitempty"`
	CommonLabels      *Labels      `json:"commonLabels,omitempty"`
	CommonAnnotations *Annotations `json:"commonAnnotations,omitempty"`
	ExternalURL       string       `json:"externalURL,omitempty"`
	Version           string       `json:"version,omitempty"`
	GroupKey          string       `json:"groupKey,omitempty"`
}

// EventHandler provides methods to handle alert events.
type EventHandler struct {
	log           logging.Interface
	kubeClientset kubernetes.Interface
}

// NewEventHandler returns a new *EventHandler.
func NewEventHandler(log logging.Interface, kubeClientset kubernetes.Interface) *EventHandler {
	return &EventHandler{
		log:           log,
		kubeClientset: kubeClientset,
	}
}

// Handle handles the provided data, parsing it into an Event and then deleting
// pods in the alerts of the event.
func (handler *EventHandler) Handle(ctx context.Context, data []byte) error {
	handler.log.Debug("Handling AlertManager event")

	event, err := parseEvent(data)
	if err != nil {
		return err
	}

	handler.log.WithField("event", *event).Info("Parsed event")

	for _, alert := range event.Alerts {
		if alert.Status != statusFiring {
			continue
		}

		if alert.Labels.Severity != severityCritical {
			continue
		}

		err = handler.deletePod(ctx, alert)
		if err != nil {
			handler.log.WithError(err).Error(errorDeletingPod)
			continue
		}

		handler.log.WithFields(logrus.Fields{
			"namespace": alert.Labels.Namespace,
			"pod":       alert.Labels.Pod,
		}).Info("Pod deleted")
	}

	return nil
}

func (handler *EventHandler) deletePod(ctx context.Context, alert *Alert) error {
	return handler.kubeClientset.CoreV1().Pods(alert.Labels.Namespace).Delete(ctx, alert.Labels.Pod, metav1.DeleteOptions{})
}

func parseEvent(data []byte) (*Event, error) {
	event := &Event{}

	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
