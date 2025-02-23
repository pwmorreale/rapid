//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package sanity performs a check on the scenario configuration.
package sanity

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
)

var log *logger.Context

func checkRequestCookies(request *config.Request) {

	for i := range request.Cookies {
		log.Info(request, nil, "Parsing cookie value: %s", request.Cookies[i].Value)
		cookies, err := http.ParseCookie(request.Cookies[i].Value)
		if err != nil {
			log.Error(request, nil, "parsing cookie: %s", err)
			continue
		}

		log.Debug(request, nil, "contains %d cookies", len(cookies))

		for n := range cookies {
			err := cookies[n].Valid()
			if err != nil {
				log.Error(request, nil, "Invalid cookie: %s", err)
			}
			log.Debug(request, nil, "cookie: %s is valid", cookies[n].String())
		}
	}
}

func checkRequestContent(request *config.Request) {

	if !strings.Contains(request.ContentType, "/") {
		log.Warn(request, nil, "ContentType: %s not in form of type/subtype", request.ContentType)
	}

	mime := mimetype.Lookup(request.ContentType)
	if mime == nil {
		log.Error(request, nil, "invalid ContentType: %s not a recognized mime type", request.ContentType)
		return
	}

	mediaType := mimetype.Detect([]byte(request.Content))
	if !mediaType.Is(request.ContentType) {
		log.Error(request, nil, "mismatched content/types:  Content_Type: %s, detected content as: %s", request.ContentType, mediaType)
	}

}

func checkResponse(_ *config.Request, _ *config.Response) {

}

func checkURL(request *config.Request) {

	u, err := url.ParseRequestURI(request.URL)
	if err != nil {
		log.Error(request, nil, "URL error: %v", err)
	}

	if u.Scheme == "https" {
		// Validate TLS config....
		log.Info(request, nil, "validating TLS config")
	}
}

func checkRequest(request *config.Request) {

	if request.Name == "" {
		log.Error(request, nil, "Missing request name")
	}

	if request.TimeLimit == 0 {
		log.Warn(request, nil, "Missing request time limit, default is infinity")
	}

	checkURL(request)

	checkRequestCookies(request)

	if request.Content == "" && request.ContentType != "" {
		log.Error(request, nil, "ContentType defined, but no content")
	}

	if request.Content != "" && request.ContentType == "" {
		log.Error(request, nil, "Content defined, but no ContentType specified")
	}

	checkRequestContent(request)

	if len(request.Responses) == 0 {
		log.Warn(request, nil, "No responses defined")
	}
}

// Check verifies a scenario configuration.
func Check(scenarioFile string, opts *logger.Options) int {

	log = logger.New(opts)

	c := config.New()
	sc, err := c.ParseFile(scenarioFile)
	if err != nil {
		log.Error(nil, nil, "Parse config: %s", err.Error())
		return 1
	}

	if sc.Name == "" {
		log.Error(nil, nil, "Missing scenario name")
	}

	if sc.Version == "" {
		log.Warn(nil, nil, "Missing scenario version")
	}

	if len(sc.Sequence.Requests) == 0 {
		log.Error(nil, nil, "No requests defined")
	}

	for i := range sc.Sequence.Requests {

		request := &sc.Sequence.Requests[i]

		log.Info(request, nil, "Request check started")
		checkRequest(request)

		for n := range request.Responses {
			checkResponse(request, &request.Responses[n])
		}
		log.Info(request, nil, "Request check complete")
	}

	return 0
}
