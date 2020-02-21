package alert

import "time"

// Event is the payload of an AlertManager webhook notification.
type Event struct {
	Receiver          string       `json:"receiver,omitempty"`
	Status            string       `json:"status,omitempty"`
	Alerts            []*Instance  `json:"alerts,omitempty"`
	GroupLabels       *Labels      `json:"groupLabels,omitempty"`
	CommonLabels      *Labels      `json:"commonLabels,omitempty"`
	CommonAnnotations *Annotations `json:"commonAnnotations,omitempty"`
	ExternalURL       string       `json:"externalURL,omitempty"`
	Version           string       `json:"version,omitempty"`
	GroupKey          string       `json:"groupKey,omitempty"`
}

// Instance is an alert from AlertManager within an Event.
type Instance struct {
	Status       string       `json:"status,omitempty"`
	Labels       *Labels      `json:"labels,omitempty"`
	Annotations  *Annotations `json:"annotations,omitempty"`
	StartsAt     time.Time    `json:"startsAt,omitempty"`
	EndsAt       time.Time    `json:"endsAt,omitempty"`
	GeneratorURL string       `json:"generatorURL,omitempty"`
	Fingerprint  string       `json:"fingerprint,omitempty"`
}

// Labels are metadata of an alert Instance.
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

// Annotations are extra metadata of an alert Instance.
type Annotations struct {
	Description string `json:"description,omitempty"`
}
