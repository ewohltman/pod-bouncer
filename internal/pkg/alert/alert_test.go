package alert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNewEvent(t *testing.T) {
	testEventData, testEvent, err := newTestEvent()
	if err != nil {
		t.Fatalf("Error creating new test event: %s", err)
	}

	event, err := NewEvent(testEventData)
	if err != nil {
		t.Fatalf("Error creating new test event: %s", err)
	}

	if !reflect.DeepEqual(event, testEvent) {
		t.Errorf(
			"Unexpected result. Got: %+v, Expected: %+v",
			event,
			testEvent,
		)
	}
}

func TestDeletePod(t *testing.T) {
	_, testEvent, err := newTestEvent()
	if err != nil {
		t.Fatalf("Error creating new test event: %s", err)
	}

	for _, alertInstance := range testEvent.Alerts {
		DeletePod(alertInstance)
	}
}

func newTestEvent() (eventData []byte, event *Event, err error) {
	eventData, err = ioutil.ReadFile("testdata/event.json")
	if err != nil {
		return nil, nil, fmt.Errorf("error reading testdata file: %w", err)
	}

	event = &Event{}

	err = json.Unmarshal(eventData, event)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling test event: %w", err)
	}

	return
}
