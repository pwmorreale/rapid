package service

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"

	"github.com/pwmorreale/rapid/internal/config"
)

func cookieExists(expected string, all []string) bool {

	for i := range all {
		if expected == all[i] {
			return true
		}
	}

	return false
}

func verifyNoContent(httpResponse *http.Response) error {

	if httpResponse.ContentLength > 0 {
		return fmt.Errorf("response ContentLength: %d", httpResponse.ContentLength)
	}

	buf := make([]byte, 10)

	n, err := httpResponse.Body.Read(buf)

	if err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("response Body contained data")
	}
	return nil
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

// VerifyHeaders verifies headers in the response.
func (s *Context) VerifyHeaders(httpResponse *http.Response, response *config.Response) error {

	for i := range response.Headers {

		err := verifyHeaderValues(httpResponse.Header, &response.Headers[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyCookies verifies any cookies returned in the response.
func (s *Context) VerifyCookies(httpResponse *http.Response, response *config.Response) error {

	var AllResponseCookies []string

	r := httpResponse.Cookies()
	for i := range r {
		AllResponseCookies = append(AllResponseCookies, r[i].String())
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

// VerifyCOntentAndExtract verifies the content, if any and extracts saved values, if any.
func (s *Context) VerifyContentAndExtract(httpResponse *http.Response, response *config.Response) error {

	// no content-type configured means that none is expected, prove it.
	if response.Content.MediaType == "" {
		return verifyNoContent(httpResponse)
	}

	contentType, _, err := mime.ParseMediaType(httpResponse.Header.Get("Content-Type"))
	expectedContentType, _, err := mime.ParseMediaType(response.Content.MediaType)
	if contentType != expectedContentType {
		return fmt.Errorf("content-type: %s != %s", contentType, response.Content.MediaType)
	}
	maxSize := response.Content.MaxSize
	if maxSize == 0 {
		maxSize = config.DefaultContentLimit
	}

	lr := io.LimitReader(httpResponse.Body, int64(maxSize))
	contentBytes, err := io.ReadAll(lr)
	if err != nil {
		return err
	}

	// Verify contents vs. Content-Type
	mediaType := http.DetectContentType(contentBytes)
	if mediaType != contentType {
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
			v, err = s.datum.ExtractJSON(e.Path, rb)
			break
		case "xml":
			v, err = s.datum.ExtractXML(e.Path, rb)
			break
		case "text":
			v, err = s.datum.ExtractRegex(e.Path, rb)
			break
		}

		if err != nil {
			return err
		}

		err = s.datum.AddReplacement(e.Name, v)
		if err != nil {
			return err
		}
	}

	return nil
}
