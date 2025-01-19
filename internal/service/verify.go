package service

import (
	"fmt"
	"net/http"

	"github.com/pwmorreale/rapid/internal/config"
)

func (s *Context) VerifyHeaderValues(httpHeaders http.Header, expectedHeader *config.HeaderData) error {

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

func (s *Context) VerifyHeaders(httpResponse *http.Response, response *config.Response, request *config.Request) error {

	for i := range response.Headers {

		err := s.VerifyHeaderValues(httpResponse.Header, &response.Headers[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Context) VerifyCookies(_ *http.Response, _ *config.Response, _ *config.Request) error {
	return nil
}

func (s *Context) VerifyContent(_ *http.Response, _ *config.Response, _ *config.Request) error {
	return nil
}
