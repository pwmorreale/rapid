//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package verify performs a check on the scenario configuration.
package verify

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/logger"
)

// CheckCookies verifies cookie syntax.
func CheckCookies(request *config.Request, response *config.Response, cookies []config.CookieData) {

	for i := range cookies {
		cookie := cookies[i]

		logger.Info(request, response, "parsing cookie value: %s", cookie.Value)
		cookies, err := http.ParseCookie(cookie.Value)
		if err != nil {
			logger.Error(request, response, "parsing cookie: %s", err)
			continue
		}

		logger.Debug(request, response, "contains %d cookies", len(cookies))

		for n := range cookies {
			err := cookies[n].Valid()
			if err != nil {
				logger.Error(request, response, "Invalid cookie: %s", err)
			} else {
				logger.Debug(request, response, "cookie: %s is valid", cookies[n].String())
			}
		}
	}
}

// CheckHeaders verifies headers
func CheckHeaders(request *config.Request, response *config.Response, headers []config.HeaderData) {

	for i := range headers {
		if headers[i].Value != "" || headers[i].Name == "" {
			logger.Error(request, response, "missing header name, but have value: %s", headers[i].Value)
		}
	}
}

// CheckRequestContent verifies content and content type.
func CheckRequestContent(request *config.Request) {

	logger.Info(request, nil, "checking content and content_type")

	if request.Content == "" && request.ContentType == "" {
		return
	}

	if request.Content != "" && request.ContentType == "" {
		logger.Error(request, nil, "mismatched content/type, content_type is blank, but have content")
	}

	if request.Content == "" && request.ContentType != "" {
		logger.Error(request, nil, "mismatched content/type, content is blank, but have content_type")
	}

	if !strings.Contains(request.ContentType, "/") {
		logger.Warn(request, nil, "content_type: %s not in form of type/subtype", request.ContentType)
	}

	mime := mimetype.Lookup(request.ContentType)
	if mime == nil {
		logger.Error(request, nil, "invalid content_type: %s not a recognized mime type", request.ContentType)
		return
	}

	mediaType := mimetype.Detect([]byte(request.Content))
	if !mediaType.Is(request.ContentType) {
		logger.Error(request, nil, "mismatched content/types:  Content_Type: %s, detected content as: %s", request.ContentType, mediaType)
	}

}

// CheckResponseContent checks the response content.
func CheckResponseContent(request *config.Request, response *config.Response) {

	if response.Content.Expected && response.Content.MediaType == "" {
		logger.Error(request, response, "response content expected, but no content_type specified")
	}

	if !response.Content.Expected && response.Content.MediaType != "" {
		logger.Warn(request, response, "response content_type specified, but no content.expected is false")
	}

	if response.Content.MediaType != "" {
		mime := mimetype.Lookup(request.ContentType)
		if mime == nil {
			logger.Error(request, response, "invalid content_type: %s not a recognized mime type", response.Content.MediaType)
		}
	}

	for i := range response.Content.Extract {

		if response.Content.Extract[i].Type == "" {
			logger.Error(request, response, "extract type must be defined")
		}

		if response.Content.Extract[i].Path == "" {
			logger.Error(request, response, "extract path must be defined")
		}

		if response.Content.Extract[i].Name == "" {
			logger.Error(request, response, "extract data_name must be defined")
		}
	}

}

// CheckResponse verifies a response
func CheckResponse(request *config.Request, response *config.Response) {

	if http.StatusText(response.StatusCode) == "" {
		logger.Error(request, response, "invalid status code: %d", response.StatusCode)
	}

	CheckCookies(request, response, response.Cookies)

	CheckHeaders(request, response, response.Headers)

	CheckResponseContent(request, response)

}

// CheckURL verifies the URL
func CheckURL(request *config.Request) {

	u, err := url.ParseRequestURI(request.URL)
	if err != nil {
		logger.Error(request, nil, "URL error: %v", err)
	}

	if u.Scheme == "https" {
		// Validate TLS config....
		logger.Info(request, nil, "validating TLS config")
	}
}

// CheckThunderingHerd checks the herd configuration
func CheckThunderingHerd(request *config.Request) {

	if request.ThunderingHerd.Max == 0 {
		logger.Warn(request, nil, "missing thundering_herd.maximum_requests, default is 1")
	}

	if request.ThunderingHerd.Size == 0 {
		logger.Warn(request, nil, "missing thundering_herd.active_size, default is 1")
	}
}

// CheckRequest verifies a request
func CheckRequest(request *config.Request) {

	if request.Name == "" {
		logger.Error(request, nil, "missing request name")
	}

	CheckURL(request)

	CheckCookies(request, nil, request.Cookies)

	CheckRequestContent(request)

	CheckThunderingHerd(request)

	if len(request.Responses) == 0 {
		logger.Error(request, nil, "no responses defined")
	}

	CheckHeaders(request, nil, request.ExtraHeaders)
}

// CheckReplacements verifies the replacement data.
func CheckReplacements(r []config.ReplaceData) {

	for i := range r {
		if r[i].Regex != "" && r[i].Value == "" {
			logger.Error(nil, nil, "missing value for keyword: %s", r[i].Regex)
		}

		if r[i].Regex == "" && r[i].Value != "" {
			logger.Error(nil, nil, "missing keyword for value: %s", r[i].Value)
		}
	}
}

// Check verifies a scenario configuration.
func Check(scenarioFile string) error {

	c := config.New()
	sc, err := c.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	if sc.Name == "" {
		logger.Error(nil, nil, "missing scenario name")
	}

	if sc.Version == "" {
		logger.Warn(nil, nil, "missing scenario version")
	}

	CheckReplacements(sc.Replacements)

	if len(sc.Sequence.Requests) == 0 {
		logger.Error(nil, nil, "no requests defined")
	}

	for i := range sc.Sequence.Requests {

		request := &sc.Sequence.Requests[i]

		logger.Info(request, nil, "request check started")
		CheckRequest(request)

		for n := range request.Responses {
			CheckResponse(request, request.Responses[n])
		}
		logger.Info(request, nil, "request check complete")
	}

	return nil
}
