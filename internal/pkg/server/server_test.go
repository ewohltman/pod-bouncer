package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ewohltman/pod-bouncer/internal/pkg/logging"
)

const (
	testPort = "8080"
	testURL  = "http://localhost:" + testPort

	httpClientTimeout = time.Second
	contextTimeout    = time.Second
)

func TestNew(t *testing.T) {
	log := logging.New()
	log.Out = ioutil.Discard

	testServer := New(log, testPort)

	go func() {
		err := testServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Test server error: %s", err)
		}
	}()

	client := &http.Client{Timeout: httpClientTimeout}

	err := testAlertEndpoint(client)
	if err != nil {
		t.Errorf("Error testing alert endpoint: %s", err)
	}

	err = testRootEndpoint(client)
	if err != nil {
		t.Errorf("Error testing root endpoint: %s", err)
	}

	ctx, cancelContext := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelContext()

	err = testServer.Shutdown(ctx)
	if err != nil {
		t.Errorf("Error closing test server: %s", err)
	}
}

func testAlertEndpoint(client *http.Client) error {
	resp, err := doRequest(client, http.MethodPost, alertEndpoint, nil)
	if err != nil {
		return err
	}

	err = drainCloseResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func testRootEndpoint(client *http.Client) error {
	resp, err := doRequest(client, http.MethodGet, rootEndpoint, nil)
	if err != nil {
		return err
	}

	err = drainCloseResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func doRequest(client *http.Client, method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, testURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error creating test request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing test request: %w", err)
	}

	return resp, nil
}

func drainCloseResponse(resp *http.Response) (err error) {
	defer func() {
		err = closeResponse(resp, err)
	}()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		err = fmt.Errorf("error draining test response body: %w", err)
	}

	return
}

func closeResponse(resp *http.Response, err error) error {
	closeErr := resp.Body.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("error closing test response body: %w", closeErr)

		if err != nil {
			return fmt.Errorf("%s: %w", closeErr, err)
		}

		return closeErr
	}

	return err
}
