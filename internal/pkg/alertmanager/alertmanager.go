// Package alertmanager contains type definitions and functionality for
// handling events from AlertManager.
package alertmanager

import (
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

// Handler provides methods to handle alert events.
type Handler struct {
	log           logging.Interface
	kubeClientset kubernetes.Interface
}

// NewHandler returns a new *Handler.
func NewHandler(log logging.Interface, kubeClientset kubernetes.Interface) *Handler {
	return &Handler{
		log:           log,
		kubeClientset: kubeClientset,
	}
}

// DeletePod deletes a running pod, allowing Kubernetes to reschedule it.
func (handler *Handler) DeletePod(alert *Alert) error {
	return handler.kubeClientset.CoreV1().Pods(alert.Labels.Namespace).Delete(alert.Labels.Pod, &metav1.DeleteOptions{})
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

// NewEvent unmarshals data and returns an *Event.
func NewEvent(data []byte) (*Event, error) {
	event := &Event{}

	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
