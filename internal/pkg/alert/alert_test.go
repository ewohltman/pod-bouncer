package alert

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNewEvent(t *testing.T) {
	testEvent, err := ioutil.ReadFile("testdata/event.json")
	if err != nil {
		t.Fatalf("Error reading testdata file: %s", err)
	}

	expected := &Event{}

	err = json.Unmarshal(testEvent, expected)
	if err != nil {
		t.Fatalf("Error unmarshaling test event: %s", err)
	}

	actual, err := NewEvent(testEvent)
	if err != nil {
		t.Fatalf("Error creating new test event: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"Unexpected result. Got: %+v, Expected: %+v",
			actual,
			expected,
		)
	}
}
