//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package report generates structured output from scenario results.
package report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/pwmorreale/rapid/config"
)

// RequestResult holds results for a single request.
type RequestResult struct {
	Name     string           `json:"name" xml:"name,attr"`
	Method   string           `json:"method" xml:"method,attr"`
	Count    int64            `json:"count" xml:"count,attr"`
	Errors   int64            `json:"errors" xml:"errors,attr"`
	MinTime  string           `json:"min_time" xml:"min-time,attr"`
	MaxTime  string           `json:"max_time" xml:"max-time,attr"`
	AvgTime  string           `json:"avg_time" xml:"avg-time,attr"`
	Responses []ResponseResult `json:"responses" xml:"response"`
}

// ResponseResult holds results for a single response.
type ResponseResult struct {
	Name       string `json:"name" xml:"name,attr"`
	StatusCode int    `json:"status_code" xml:"status-code,attr"`
	Count      int64  `json:"count" xml:"count,attr"`
	Errors     int64  `json:"errors" xml:"errors,attr"`
	MinTime    string `json:"min_time" xml:"min-time,attr"`
	MaxTime    string `json:"max_time" xml:"max-time,attr"`
	AvgTime    string `json:"avg_time" xml:"avg-time,attr"`
}

// Summary holds the full scenario report.
type Summary struct {
	Name       string          `json:"name" xml:"name,attr"`
	Version    string          `json:"version" xml:"version,attr"`
	Timestamp  string          `json:"timestamp" xml:"timestamp,attr"`
	Iterations int             `json:"iterations" xml:"iterations,attr"`
	Requests   []RequestResult `json:"requests" xml:"request"`
}

func avgDuration(total time.Duration, count int64) string {
	if count == 0 {
		return "0s"
	}
	return (total / time.Duration(count)).String()
}

// BuildSummary creates a Summary from a completed scenario.
func BuildSummary(sc *config.Scenario) *Summary {
	s := &Summary{
		Name:       sc.Name,
		Version:    sc.Version,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Iterations: sc.Sequence.Iterations,
	}

	for i := range sc.Sequence.Requests {
		req := &sc.Sequence.Requests[i]
		rr := RequestResult{
			Name:    req.Name,
			Method:  req.Method,
			Count:   req.Stats.GetCount(),
			Errors:  req.Stats.GetErrors(),
			MinTime: req.Stats.GetMinDuration().String(),
			MaxTime: req.Stats.GetMaxDuration().String(),
			AvgTime: avgDuration(req.Stats.GetDuration(), req.Stats.GetCount()),
		}

		for j := range req.Responses {
			resp := req.Responses[j]
			rr.Responses = append(rr.Responses, ResponseResult{
				Name:       resp.Name,
				StatusCode: resp.StatusCode,
				Count:      resp.Stats.GetCount(),
				Errors:     resp.Stats.GetErrors(),
				MinTime:    resp.Stats.GetMinDuration().String(),
				MaxTime:    resp.Stats.GetMaxDuration().String(),
				AvgTime:    avgDuration(resp.Stats.GetDuration(), resp.Stats.GetCount()),
			})
		}

		for j := range req.UnknownResponses {
			resp := req.UnknownResponses[j]
			rr.Responses = append(rr.Responses, ResponseResult{
				Name:       resp.Name,
				StatusCode: resp.StatusCode,
				Count:      resp.Stats.GetCount(),
				Errors:     resp.Stats.GetErrors(),
				MinTime:    resp.Stats.GetMinDuration().String(),
				MaxTime:    resp.Stats.GetMaxDuration().String(),
				AvgTime:    avgDuration(resp.Stats.GetDuration(), resp.Stats.GetCount()),
			})
		}

		s.Requests = append(s.Requests, rr)
	}

	return s
}

// WriteJSON writes the summary as JSON to the given file.
func WriteJSON(path string, sc *config.Scenario) error {
	s := BuildSummary(sc)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// JUnit XML types

// JUnitTestSuites is the top-level element.
type JUnitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents one request.
type JUnitTestSuite struct {
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Time     string          `xml:"time,attr"`
	Cases    []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents one response within a request.
type JUnitTestCase struct {
	Name    string        `xml:"name,attr"`
	Time    string        `xml:"time,attr"`
	Failure *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitFailure describes a failed test case.
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

// WriteJUnit writes the summary as JUnit XML to the given file.
func WriteJUnit(path string, sc *config.Scenario) error {
	s := BuildSummary(sc)

	suites := JUnitTestSuites{}

	for _, req := range s.Requests {
		suite := JUnitTestSuite{
			Name:     req.Name,
			Tests:    0,
			Failures: 0,
			Time:     req.AvgTime,
		}

		for _, resp := range req.Responses {
			suite.Tests++
			tc := JUnitTestCase{
				Name: fmt.Sprintf("%s/%s [%d]", req.Name, resp.Name, resp.StatusCode),
				Time: resp.AvgTime,
			}
			if resp.Errors > 0 {
				suite.Failures++
				tc.Failure = &JUnitFailure{
					Message: fmt.Sprintf("%d errors in %d executions", resp.Errors, resp.Count+resp.Errors),
					Type:    "ValidationError",
				}
			}
			suite.Cases = append(suite.Cases, tc)
		}

		// Request-level errors that didn't match any configured response.
		var responseErrors int64
		for _, resp := range req.Responses {
			responseErrors += resp.Errors
		}
		unmatched := req.Errors - responseErrors
		if unmatched > 0 {
			suite.Tests++
			suite.Failures++
			suite.Cases = append(suite.Cases, JUnitTestCase{
				Name: fmt.Sprintf("%s (request errors)", req.Name),
				Time: req.AvgTime,
				Failure: &JUnitFailure{
					Message: fmt.Sprintf("%d request-level errors (connection failures or unconfigured responses)", unmatched),
					Type:    "ExecutionError",
				},
			})
		}

		suites.Suites = append(suites.Suites, suite)
	}

	data, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return err
	}

	output := append([]byte(xml.Header), data...)
	return os.WriteFile(path, output, 0644)
}
