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

func readContent(httpResponse *http.Response, response *config.Response) (int, []byte, error) {

	// Default if not specified
	maxSize := int64(4096)
	if response.Content.MaxSize != 0 {
		maxSize = int64(response.Content.MaxSize)
	}

	size := httpResponse.ContentLength
	if size < 0 || size > maxSize {
		size = maxSize
	}

	buf := make([]byte, size)
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

	// N.B.  We may not know the exact response body length, so only check what we can.
	switch {
	case (httpLength == 0 && contentLength != 0):
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpLength, contentLength)
	case (httpLength > 0 && contentLength == 0):
		return fmt.Errorf("mismatched Content-Length header (%d) and actual content (at least %d bytes)", httpLength, contentLength)
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

	err = r.verifyContentLength(httpResponse.ContentLength, int64(nrBytes))
	if err != nil {
		return err
	}

	err = r.verifyExpectedContentType(contentBytes, httpResponse, response)
	if err != nil {
		return err
	}

	err = r.verifyContains(contentBytes, response)
	if err != nil {
		return err
	}

	return r.extractContent(contentBytes, response)
}

func lookupResponse(statusCode int, r []*config.Response) *config.Response {

	for i := range r {
		if statusCode == r[i].StatusCode {
			return r[i]
		}
	}

	return nil
}

func (r *Context) findResponse(httpResponse *http.Response, request *config.Request) *config.Response {

	resp := lookupResponse(httpResponse.StatusCode, request.Responses)
	if resp != nil {
		return resp
	}

	// Not found on the configured array, so find/add to the unkowns array.

	unknownResponseMutex.Lock()
	defer unknownResponseMutex.Unlock()

	resp = lookupResponse(httpResponse.StatusCode, request.UnknownResponses)
	if resp != nil {
		return resp
	}

	resp = new(config.Response)

	request.UnknownResponses = append(request.UnknownResponses, resp)

	resp.Name = config.DefaultResponseName
	resp.StatusCode = httpResponse.StatusCode

	return resp
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

	return r.verifyContentAndExtract(httpResponse, response)
}

func (r *Context) validateResponse(httpResponse *http.Response, request *config.Request) (*config.Response, error) {

	response := r.findResponse(httpResponse, request)
	err := r.verifyResponse(httpResponse, response)
	return response, err
}
