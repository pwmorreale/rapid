//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes REST calls
package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pwmorreale/rapid/config"
)

func cookieExists(expected string, all []string) bool {

	for i := range all {
		if expected == all[i] {
			return true
		}
	}

	return false
}

func readContent(httpResponse *http.Response, response *config.Response) (int, []byte, error) {

	maxSize := httpResponse.ContentLength
	if maxSize < 0 || maxSize > int64(response.Content.MaxSize) {
		maxSize = int64(response.Content.MaxSize)
	}

	buf := make([]byte, maxSize)
	n, err := httpResponse.Body.Read(buf)
	if err != nil && err != io.EOF {
		return 0, nil, err
	}

	return n, buf[:n], nil
}

func verifyHeaderValues(httpHeaders http.Header, expectedHeader *config.HeaderData) error {

	name := http.CanonicalHeaderKey(expectedHeader.Name)
	v := httpHeaders.Values(name)
	if len(v) == 0 {
		return fmt.Errorf("header: %s not found", name)
	}

	for n := range v {
		if v[n] == expectedHeader.Value {
			return nil
		}
	}

	return fmt.Errorf("header: %s, expected value (%s) not found", name, expectedHeader.Value)
}

func (r *Context) verifyHeaders(httpResponse *http.Response, response *config.Response) error {

	for i := range response.Headers {

		err := verifyHeaderValues(httpResponse.Header, &response.Headers[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Context) verifyCookies(httpResponse *http.Response, response *config.Response) error {

	var AllResponseCookies []string

	c := httpResponse.Cookies()
	for i := range c {
		AllResponseCookies = append(AllResponseCookies, c[i].String())
	}

	// Allow mutiple cookies within the config.
	for i := range response.Cookies {

		e, err := http.ParseSetCookie(response.Cookies[i].Value)
		if err != nil {
			return err
		}

		expectedCookie := e.String()
		ok := cookieExists(expectedCookie, AllResponseCookies)
		if !ok {
			return fmt.Errorf("cookie: %s not found in response", expectedCookie)
		}

	}
	return nil
}

func (r *Context) verifyContentAndExtract(httpResponse *http.Response, response *config.Response) error {

	nrBytes, contentBytes, err := readContent(httpResponse, response)
	if err != nil {
		return err
	}

	// Check for no expected content...
	if !response.Content.Expected {
		if nrBytes > 0 {
			return fmt.Errorf("no content expected yet read: %d response bytes", nrBytes)
		}
		return nil
	}

	contentType, _, err := mime.ParseMediaType(httpResponse.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	expectedContentType, _, err := mime.ParseMediaType(response.Content.MediaType)
	if err != nil {
		return err
	}

	if contentType != expectedContentType {
		return fmt.Errorf("content-type: %s != %s", contentType, response.Content.MediaType)
	}

	switch {
	case (httpResponse.ContentLength == 0 && nrBytes != 0):
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpResponse.ContentLength, nrBytes)
	case (httpResponse.ContentLength > 0 && nrBytes == 0):
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpResponse.ContentLength, nrBytes)
	}

	// Verify contents vs. Content-Type
	mediaType := mimetype.Detect(contentBytes)
	if !mediaType.Is(contentType) {
		return fmt.Errorf("mismatched content/types:  Content_Type: %s, detected content as: %s", contentType, mediaType)
	}

	for i := range response.Content.Contains {
		ok, err := regexp.Match(response.Content.Contains[i], contentBytes)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("content sequence not found: %s", response.Content.Contains[i])
		}
	}

	var v string

	for i := range response.Content.Extract {
		e := &response.Content.Extract[i]

		rb := bytes.NewReader(contentBytes)

		switch e.Type {
		case "json":
			v, err = r.datum.ExtractJSON(e.Path, rb)
		case "xml":
			v, err = r.datum.ExtractXML(e.Path, rb)
		case "text":
			v, err = r.datum.ExtractRegex(e.Path, rb)
		}

		if err != nil {
			return err
		}

		err = r.datum.AddReplacement(e.Name, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Context) findResponse(httpResponse *http.Response, request *config.Request) *config.Response {

	for i := range request.Responses {
		if httpResponse.StatusCode == request.Responses[i].StatusCode {
			return &request.Responses[i]
		}
	}
	return nil
}

func (r *Context) verifyResponse(httpResponse *http.Response, response *config.Response) error {

	err := r.verifyHeaders(httpResponse, response)
	if err != nil {
		return err
	}

	err = r.verifyCookies(httpResponse, response)
	if err != nil {
		return err
	}

	err = r.verifyContentAndExtract(httpResponse, response)
	if err != nil {
		return err
	}

	return nil
}

func (r *Context) validateResponse(httpResponse *http.Response, request *config.Request) error {

	response := r.findResponse(httpResponse, request)
	if response == nil {
		return fmt.Errorf("response not found for status code: %d for request: %s", response.StatusCode, request.Name)
	}

	err := r.verifyResponse(httpResponse, response)

	return err
}
