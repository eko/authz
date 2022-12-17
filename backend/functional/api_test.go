//nolint:unused
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
)

var (
	baseURL = "http://localhost:8080"
)

type apiFeature struct {
	httpClient *http.Client
	req        *http.Request
	resp       *http.Response
	token      string
}

func (a *apiFeature) reset(*godog.Scenario) error {
	a.req = nil
	a.resp = nil
	return nil
}

func (a *apiFeature) iSendRequestTo(method, endpoint string) error {
	return a.httpCall(method, endpoint, nil, nil)
}

func (a *apiFeature) iSendRequestToWithPayload(method, endpoint string, body *godog.DocString) error {
	return a.httpCall(method, endpoint, body, nil)
}

func (a *apiFeature) httpCall(method, endpoint string, body *godog.DocString, writer *multipart.Writer) error {
	var reader io.Reader

	if body != nil {
		reader = strings.NewReader(body.Content)
	}

	url := baseURL + endpoint

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return fmt.Errorf("unable to prepare http request: %w", err)
	}

	if a.token != "" {
		req.Header.Add("Authorization", "Bearer "+a.token)
	}

	if writer != nil {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to query http request: %w", err)
	}

	a.req = req
	a.resp = resp

	return nil
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if a.resp == nil {
		return fmt.Errorf("http response is nil")
	}

	if code != a.resp.StatusCode {
		if a.resp.StatusCode >= 400 {
			bodyBytes, err := io.ReadAll(a.resp.Body)
			if err != nil {
				return fmt.Errorf("unable to read request body: %v", err)
			}

			return fmt.Errorf("expected response code to be: %d, but actual is: %d, response message: %s", code, a.resp.StatusCode, string(bodyBytes))
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.StatusCode)
	}
	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(body *godog.DocString) (err error) {
	if a.resp == nil {
		return fmt.Errorf("http response is nil")
	}

	var expected, actual interface{}

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	bodyBytes, err := io.ReadAll(a.resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body: %v", err)
	}
	// re-encode actual response too
	if err = json.Unmarshal(bodyBytes, &actual); err != nil {
		return
	}

	// the matching may be adapted per different requirementa.
	if !reflect.DeepEqual(expected, actual) {
		expectedBytes, err := json.MarshalIndent(expected, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal expected to JSON: %v", expected)
		}
		actualBytes, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal actuel JSON: %v", actual)
		}

		return fmt.Errorf("expected JSON does not match.\n-> expected:\n%v\n-> actual:\n%v", string(expectedBytes), string(actualBytes))
	}
	return nil
}
