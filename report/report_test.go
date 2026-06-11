//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

package report

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pwmorreale/rapid/config"
	"github.com/stretchr/testify/assert"
)

func makeScenario() *config.Scenario {
	sc := &config.Scenario{
		Name:    "test-scenario",
		Version: "1.0",
		Sequence: config.Sequence{
			Iterations: 5,
			Requests: []config.Request{
				{
					Name:   "get-users",
					Method: "get",
					Responses: []*config.Response{
						{Name: "success", StatusCode: 200},
						{Name: "not-found", StatusCode: 404},
					},
				},
			},
		},
	}

	start := time.Now().Add(-50 * time.Millisecond)
	sc.Sequence.Requests[0].Stats.Success(start)
	sc.Sequence.Requests[0].Stats.Success(start)
	sc.Sequence.Requests[0].Stats.Error(start)

	sc.Sequence.Requests[0].Responses[0].Stats.Success(start)
	sc.Sequence.Requests[0].Responses[0].Stats.Success(start)
	sc.Sequence.Requests[0].Responses[1].Stats.Error(start)

	return sc
}

func TestBuildSummary(t *testing.T) {

	sc := makeScenario()
	s := BuildSummary(sc)

	assert.Equal(t, "test-scenario", s.Name)
	assert.Equal(t, "1.0", s.Version)
	assert.Equal(t, 5, s.Iterations)
	assert.Len(t, s.Requests, 1)

	req := s.Requests[0]
	assert.Equal(t, "get-users", req.Name)
	assert.Equal(t, "get", req.Method)
	assert.Equal(t, int64(2), req.Count)
	assert.Equal(t, int64(1), req.Errors)
	assert.Len(t, req.Responses, 2)

	assert.Equal(t, "success", req.Responses[0].Name)
	assert.Equal(t, int64(2), req.Responses[0].Count)
	assert.Equal(t, "not-found", req.Responses[1].Name)
	assert.Equal(t, int64(1), req.Responses[1].Errors)
}

func TestWriteJSON(t *testing.T) {

	sc := makeScenario()
	path := filepath.Join(t.TempDir(), "report.json")

	err := WriteJSON(path, sc)
	assert.Nil(t, err)

	data, err := os.ReadFile(path)
	assert.Nil(t, err)

	var s Summary
	err = json.Unmarshal(data, &s)
	assert.Nil(t, err)

	assert.Equal(t, "test-scenario", s.Name)
	assert.Len(t, s.Requests, 1)
	assert.Equal(t, int64(2), s.Requests[0].Count)
}

func TestWriteJUnit(t *testing.T) {

	sc := makeScenario()
	path := filepath.Join(t.TempDir(), "report.xml")

	err := WriteJUnit(path, sc)
	assert.Nil(t, err)

	data, err := os.ReadFile(path)
	assert.Nil(t, err)

	var suites JUnitTestSuites
	err = xml.Unmarshal(data, &suites)
	assert.Nil(t, err)

	assert.Len(t, suites.Suites, 1)
	suite := suites.Suites[0]
	assert.Equal(t, "get-users", suite.Name)
	assert.Equal(t, 2, suite.Tests)
	assert.Equal(t, 1, suite.Failures)

	// First case: success response, no failure.
	assert.Len(t, suite.Cases, 2)
	assert.Nil(t, suite.Cases[0].Failure)
	// Second case: not-found response with errors.
	assert.NotNil(t, suite.Cases[1].Failure)
	assert.Contains(t, suite.Cases[1].Failure.Message, "1 errors")
}

func TestWriteJUnitRequestLevelErrors(t *testing.T) {

	sc := &config.Scenario{
		Name:    "connection-fail",
		Version: "1.0",
		Sequence: config.Sequence{
			Iterations: 1,
			Requests: []config.Request{
				{
					Name:   "broken",
					Method: "get",
					Responses: []*config.Response{
						{Name: "ok", StatusCode: 200},
					},
				},
			},
		},
	}

	// Simulate request-level errors (no response matched).
	start := time.Now().Add(-10 * time.Millisecond)
	sc.Sequence.Requests[0].Stats.Error(start)
	sc.Sequence.Requests[0].Stats.Error(start)
	sc.Sequence.Requests[0].Stats.Error(start)

	path := filepath.Join(t.TempDir(), "report.xml")
	err := WriteJUnit(path, sc)
	assert.Nil(t, err)

	data, err := os.ReadFile(path)
	assert.Nil(t, err)

	var suites JUnitTestSuites
	err = xml.Unmarshal(data, &suites)
	assert.Nil(t, err)

	suite := suites.Suites[0]
	// 1 configured response + 1 request-error case
	assert.Equal(t, 2, suite.Tests)
	assert.Equal(t, 1, suite.Failures)

	// Last case should be the request-level error.
	last := suite.Cases[len(suite.Cases)-1]
	assert.NotNil(t, last.Failure)
	assert.Equal(t, "ExecutionError", last.Failure.Type)
	assert.Contains(t, last.Failure.Message, "3 request-level errors")
}
