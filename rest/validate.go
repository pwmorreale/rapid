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
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pwmorreale/rapid/config"
)

// Used for adding unknown response structs to a request.
// Using a global lock here prevents 'vet' warnings for atomic noCopy.
// While less than perfect, the only contention is ffrom the same request. (eg: thundering herd resolutions)
var unknownResponseMutex sync.Mutex

func cookieExists(expected string, all []string) bool {

	for i := range all {
		if expected == all[i] {
			return true
		}
	}

	return false
}

func readBody(httpResponse *http.Response, maxSize int64) ([]byte, error) {

	if maxSize == 0 {
		maxSize = int64(config.DefaultContentLimit)
	}

	r := io.LimitReader(httpResponse.Body, maxSize)
	return io.ReadAll(r)
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

func (r *Context) verifyExpectedContentType(contentBytes []byte, httpResponse *http.Response, response *config.Response) error {
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

	// Verify contents vs. Content-Type
	mediaType := mimetype.Detect(contentBytes)
	if !mediaType.Is(contentType) {
		return fmt.Errorf("mismatched content/types:  Content_Type: %s, detected content as: %s", contentType, mediaType)
	}

	return nil
}

func (r *Context) extractContent(contentBytes []byte, response *config.Response) error {

	var v string
	var err error

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
		default:
			return fmt.Errorf("unknown extract type: %q (must be json, xml, or text)", e.Type)
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

func (r *Context) verifyContains(contentBytes []byte, response *config.Response) error {

	for i := range response.Content.Contains {
		ok, err := regexp.Match(response.Content.Contains[i], contentBytes)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("content sequence not found: %s", response.Content.Contains[i])
		}
	}

	return nil
}

func (r *Context) verifyContentLength(httpLength int64, contentLength int64) error {

	// httpLength is -1 when the server doesn't send Content-Length (chunked, etc.)
	// Only validate when the server explicitly declared a length.
	if httpLength < 0 {
		return nil
	}

	switch {
	case httpLength == 0 && contentLength != 0:
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpLength, contentLength)
	case httpLength > 0 && contentLength == 0:
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpLength, contentLength)
	}

	return nil
}

func (r *Context) verifyContent(contentBytes []byte, httpResponse *http.Response, response *config.Response) error {

	nrBytes := len(contentBytes)

	// Check for no expected content...
	if !response.Content.Expected {
		if nrBytes > 0 {
			return fmt.Errorf("no content expected yet read: %d response bytes", nrBytes)
		}
		return nil
	}

	err := r.verifyContentLength(httpResponse.ContentLength, int64(nrBytes))
	if err != nil {
		return err
	}

	err = r.verifyExpectedContentType(contentBytes, httpResponse, response)
	if err != nil {
		return err
	}

	return r.verifyContains(contentBytes, response)
}

func lookupResponses(statusCode int, r []*config.Response) []*config.Response {

	var matches []*config.Response
	for i := range r {
		if statusCode == r[i].StatusCode {
			matches = append(matches, r[i])
		}
	}

	return matches
}

func (r *Context) findOrCreateUnknown(httpResponse *http.Response, request *config.Request) *config.Response {

	unknownResponseMutex.Lock()
	defer unknownResponseMutex.Unlock()

	matches := lookupResponses(httpResponse.StatusCode, request.UnknownResponses)
	if len(matches) > 0 {
		return matches[0]
	}

	resp := new(config.Response)
	request.UnknownResponses = append(request.UnknownResponses, resp)
	resp.Name = config.DefaultResponseName
	resp.StatusCode = httpResponse.StatusCode

	return resp
}

func (r *Context) verifyResponse(contentBytes []byte, httpResponse *http.Response, response *config.Response) error {

	err := r.verifyHeaders(httpResponse, response)
	if err != nil {
		return err
	}

	err = r.verifyCookies(httpResponse, response)
	if err != nil {
		return err
	}

	return r.verifyContent(contentBytes, httpResponse, response)
}

func (r *Context) validateResponse(httpResponse *http.Response, request *config.Request) (*config.Response, error) {

	// Determine max content size from configured responses.
	var maxSize int64
	for _, resp := range request.Responses {
		if int64(resp.Content.MaxSize) > maxSize {
			maxSize = int64(resp.Content.MaxSize)
		}
	}

	contentBytes, err := readBody(httpResponse, maxSize)
	if err != nil {
		return nil, err
	}

	matches := lookupResponses(httpResponse.StatusCode, request.Responses)

	// No configured response for this status code.
	if len(matches) == 0 {
		resp := r.findOrCreateUnknown(httpResponse, request)
		return resp, r.verifyResponse(contentBytes, httpResponse, resp)
	}

	// Try each matching response; succeed on the first that passes.
	var lastErr error
	for _, resp := range matches {
		err := r.verifyResponse(contentBytes, httpResponse, resp)
		if err == nil {
			return resp, r.extractContent(contentBytes, resp)
		}
		lastErr = err
	}

	return matches[0], lastErr
}
