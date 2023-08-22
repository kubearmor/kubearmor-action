// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package urlfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// URLVisitor downloads the contents of a URL, and if successful, returns
// an info object representing the downloaded object.
type URLVisitor struct {
	URL              *url.URL
	HTTPAttemptCount int
}

// readHTTPWithRetries tries to http.Get the v.URL retries times before giving up.
func readHTTPWithRetries(get httpget, duration time.Duration, u string, attempts int) (io.ReadCloser, error) {
	var err error
	if attempts <= 0 {
		return nil, fmt.Errorf("http attempts must be greater than 0, was %d", attempts)
	}
	for i := 0; i < attempts; i++ {
		var (
			statusCode int
			status     string
			body       io.ReadCloser
		)
		if i > 0 {
			time.Sleep(duration)
		}

		// Try to get the URL
		statusCode, status, body, err = get(u)

		// Retry Errors
		if err != nil {
			continue
		}

		if statusCode == http.StatusOK {
			return body, nil
		}
		err := body.Close()
		if err != nil {
			return nil, err
		}
		// Error - Set the error condition from the StatusCode
		err = fmt.Errorf("unable to read URL %q, server reported %s, status code=%d", u, status, statusCode)

		if statusCode >= 500 && statusCode < 600 {
			// Retry 500's
			continue
		} else {
			// Don't retry other StatusCodes
			break
		}
	}
	return nil, err
}

// httpget Defines function to retrieve a url and return the results.  Exists for unit test stubbing.
type httpget func(url string) (int, string, io.ReadCloser, error)

// httpgetImpl Implements a function to retrieve a url and return the results.
func httpgetImpl(url string) (int, string, io.ReadCloser, error) {
	// TODO: G107 (CWE-88): Potential HTTP request made with variable url (Confidence: MEDIUM, Severity: MEDIUM)
	resp, err := http.Get(url) // #nosec
	if err != nil {
		return 0, "", nil, err
	}
	return resp.StatusCode, resp.Status, resp.Body, nil
}

// ReadJSONFromURL reads the json file from the given url and returns it as a []byte array.
func ReadJSONFromURL(url string) ([]byte, error) {
	body, err := readHTTPWithRetries(httpgetImpl, time.Second*5, url, 3)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
