// Package alert contains type definitions and functionality for handling alert
// events.
package alert

import (
	"encoding/json"
)

// NewEvent unmarshals data and returns an *Event.
func NewEvent(data []byte) (*Event, error) {
	event := &Event{}

	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// DeletePod deletes a running pod, allowing Kubernetes to reschedule it.
func DeletePod(instance *Instance) {

}
