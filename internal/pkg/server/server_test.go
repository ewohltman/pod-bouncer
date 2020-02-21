package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

const (
	testPort = "8080"
	testURL  = "http://localhost:" + testPort
)

type testCase struct {
	req                  *http.Request
	expectedResponseCode int
}

func TestNew(t *testing.T) {
	log := logging.New()
	log.Out = ioutil.Discard

	testServer := New(log, testPort)
	if testServer == nil {
		t.Fatal("Unexpected nil *http.Server")
	}

	err := testAlertEndpoint(log)
	if err != nil {
		t.Errorf("Error testing alert endpoint: %s", err)
	}

	err = testRootEndpoint(log)
	if err != nil {
		t.Errorf("Error testing root endpoint: %s", err)
	}
}

func testAlertEndpoint(log logging.Interface) error {
	testEventData, err := ioutil.ReadFile("testdata/event.json")
	if err != nil {
		return fmt.Errorf("error reading testdata file: %w", err)
	}

	nilBodyReq, err := http.NewRequest(http.MethodPost, testURL+alertEndpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating nil body test request: %w", err)
	}

	invalidReq, err := http.NewRequest(http.MethodPost, testURL+alertEndpoint, bytes.NewReader([]byte{}))
	if err != nil {
		return fmt.Errorf("error creating invalid test request: %w", err)
	}

	validReq, err := http.NewRequest(http.MethodPost, testURL+alertEndpoint, bytes.NewReader(testEventData))
	if err != nil {
		return fmt.Errorf("error creating valid test request: %w", err)
	}

	testCases := []*testCase{
		{req: nilBodyReq, expectedResponseCode: http.StatusBadRequest},
		{req: invalidReq, expectedResponseCode: http.StatusInternalServerError},
		{req: validReq, expectedResponseCode: http.StatusOK},
	}

	return runTests(alertHandler(log), testCases)
}

func testRootEndpoint(log logging.Interface) error {
	nilBodyReq, err := http.NewRequest(http.MethodPost, testURL+rootEndpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating nil body test request: %w", err)
	}

	emptyBodyReq, err := http.NewRequest(http.MethodPost, testURL+rootEndpoint, bytes.NewReader([]byte{}))
	if err != nil {
		return fmt.Errorf("error creating empty body test request: %w", err)
	}

	testCases := []*testCase{
		{req: nilBodyReq, expectedResponseCode: http.StatusOK},
		{req: emptyBodyReq, expectedResponseCode: http.StatusOK},
	}

	return runTests(rootHandler(log), testCases)
}

func runTests(testFunc func(http.ResponseWriter, *http.Request), testCases []*testCase) error {
	for _, tc := range testCases {
		respRecorder := httptest.NewRecorder()

		testFunc(respRecorder, tc.req)

		resp := respRecorder.Result()

		if resp.StatusCode != tc.expectedResponseCode {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}

			err = resp.Body.Close()
			if err != nil {
				return fmt.Errorf("error closing response body: %w", err)
			}

			return fmt.Errorf(
				"unexpected response code (%d): %s",
				resp.StatusCode,
				string(body),
			)
		}
	}

	return nil
}
