package alertmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

func TestNewHandler(t *testing.T) {
	_, err := testHandler()
	if err != nil {
		t.Fatalf("Error creating test Handler: %s", err)
	}
}

func TestHandler_DeletePod(t *testing.T) {
	handler, err := testHandler()
	if err != nil {
		t.Fatalf("Error creating test Handler: %s", err)
	}

	_, testEvent, err := testEvent()
	if err != nil {
		t.Fatalf("Error creating test Event: %s", err)
	}

	for _, alert := range testEvent.Alerts {
		err = handler.DeletePod(alert)
		if err != nil {
			t.Errorf("Error deleting pod: %s", err)
		}
	}
}

func TestNewEvent(t *testing.T) {
	testEventData, testEvent, err := testEvent()
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

func testHandler() (*Handler, error) {
	log := logging.New()
	log.Out = ioutil.Discard

	_, testEvent, err := testEvent()
	if err != nil {
		return nil, fmt.Errorf("error creating test event: %w", err)
	}

	handler := NewHandler(log, testClientset(testEvent))

	return handler, nil
}

func testEvent() (eventData []byte, event *Event, err error) {
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

func testClientset(event *Event) *fake.Clientset {
	pods := make([]runtime.Object, len(event.Alerts))

	for i, alert := range event.Alerts {
		pod := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      alert.Labels.Pod,
				Namespace: alert.Labels.Namespace,
			},
			Spec:   corev1.PodSpec{},
			Status: corev1.PodStatus{},
		}

		pods[i] = pod
	}

	return fake.NewSimpleClientset(pods...)
}