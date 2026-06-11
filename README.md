# REST API  Diagnostic (RAPID) tool
[![Go Report Card](https://goreportcard.com/badge/github.com/pwmorreale/rapid)](https://goreportcard.com/report/github.com/pwmorreale/rapid) [![Tests & Lint](https://github.com/pwmorreale/rapid/actions/workflows/makefile.yml/badge.svg)](https://github.com/pwmorreale/rapid/actions/workflows/makefile.yml) [![CodeQL Advanced](https://github.com/pwmorreale/rapid/actions/workflows/codeql.yml/badge.svg)](https://github.com/pwmorreale/rapid/actions/workflows/codeql.yml)

Rapid is a REST API testing tool used to both verify conformance of your REST server against your API spec, as well as measure load and/or performance behavior.  Rapid also makes it possible for you to test policies such as circuit-breaking, rate-limiting, load-balancing, etc.

Rapid works entirely through a YAML configuration file.   The configuration is called a *scenario* and consists of a sequence of one or more http/https requests along with their possible responses.  Sequences can contain iteration counts or execute in a loop for a specific period of time.  You can also easily configure multiple concurrent requests to load your infrastructure.

You define content, headers, and cookies for both requests and the responses.  Rapid compares the actual response data with the expected response configuration and informs you of any discrepancies.

In addition, Rapid allows you to dynamically extract data from previous responses and insert that data into future requests.  This allows you to create dynamic paths through your service infrastructure.

## Install
To install:

```bash
go install github.com/pwmorreale/rapid@latest
```

## Build
Rapid uses a Makefile for building.  The Makefile references three other tools: [staticcheck](https://github.com/dominikh/go-tools), [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) and [revive](https://github.com/mgechev/revive).

The executable will be located in the *target* directory.

The Makefile targets are:

```bash
$ make help
help                           Display this help screen
tests                          Run all tests/lints
generate                       Generate test mocks
lint                           Lint the files
test                           Run unit tests
race                           Run race detector
staticcheck                    Run staticcheck
build                          Build
clean                          Remove previous build
coverage                       Display test coverage
$
```

## Usage

Rapid has two commands.  To execute a scenario, use the ***run*** command:

```bash
% rapid run -s ./scenario.yaml
```
You can also check a scenario configuration to find common typos/etc by using the ***verify*** command:

```bash
% rapid verify -s ./scenario.yaml
```

The verify command will exit with a non-zero status if any errors are found.

There are also several options for controlling log messages.  See the help for the above commands.

## Quick Start

Here is a minimal scenario that performs a GET request and verifies a 200 response:

```yaml
name: quickstart
version: "1.0"
sequence:
  iterations: 1
  requests:
    - name: health-check
      method: get
      url: https://httpbin.org/get
      responses:
        - name: success
          status_code: 200
          content:
            expected: true
            content_type: application/json
            contains:
              - "httpbin.org"
```

Run it:
```bash
% rapid run -s quickstart.yaml
```

Output on success:
```
INF execution started request.name=health-check request.method=get
INF execution complete request.name=health-check request.method=get
INF count=1 errors=0 minTime=150.23ms maxTime=150.23ms avgTime=150.23ms request.name=health-check request.method=get
INF count=1 errors=0 minTime=150.23ms maxTime=150.23ms avgTime=150.23ms request.name=health-check request.method=get response.name=success response.status=200
```

Output when a validation fails (e.g., expected header is missing):
```
ERR header: X-Custom not found request.name=health-check request.method=get
INF count=1 errors=1 minTime=148.91ms maxTime=148.91ms avgTime=148.91ms request.name=health-check request.method=get
```

## Data Extraction Example

A common pattern is to authenticate first, extract a token from the response, and use it in subsequent requests.  Here's how that looks:

```yaml
name: auth-flow
version: "1.0"
find_replace:
  - match: AUTH_TOKEN
    replace: "placeholder"
sequence:
  iterations: 1
  requests:
    - name: login
      method: post
      url: https://api.example.com/auth/login
      content: '{"username": "test", "password": "secret"}'
      content_type: application/json
      responses:
        - name: login-success
          status_code: 200
          content:
            expected: true
            content_type: application/json
            extract:
              - type: json
                path: token
                match: AUTH_TOKEN

    - name: get-profile
      method: get
      url: https://api.example.com/users/me
      extra_headers:
        - name: Authorization
          value: "Bearer AUTH_TOKEN"
      responses:
        - name: profile-success
          status_code: 200
          content:
            expected: true
            content_type: application/json
            contains:
              - "username"
```

The flow:
1. The `login` request posts credentials and receives a JSON response containing a `token` field.
2. The `extract` configuration uses a GJSON path to pull the token value and registers it as the replacement for the regex `AUTH_TOKEN`.
3. The `get-profile` request's Authorization header contains `AUTH_TOKEN`, which Rapid replaces with the actual token before sending.

This pattern works for any chained data: session IDs, CSRF tokens, resource IDs returned from creation endpoints, etc.

## Features

### Iterations
You can define an iteration count and an optional iteration time limit.  Each iteration loops through all configured requests in order.  If a time limit is set, the iteration must complete within it or it is recorded as an error.

### Find&Replace
Find&Replace allows you to predefine a set of regex terms and their associated replacement strings.  When a regex matches in an extra header value, cookie value, request content, or the URL, the replacement term is inserted in its place.  This allows you to define a value once and reference it throughout the configuration.

Note: header *names* are not passed through Find&Replace, only header *values*.

### Data Extraction
Rapid allows you to extract data from response payloads for use in future requests.  You can search through JSON, XML, or text responses and save the extracted value as a new Find&Replace entry.

Extraction only occurs after all other response validations (headers, cookies, content checks) pass successfully.  This ensures you never extract data from an invalid response.

### Thundering Herd
Rapid allows you to create *thundering herd* configurations that specify a number of concurrent requests for a specific duration of time, or a maximum total request count.  For example, you could configure Rapid to execute 1000 requests concurrently for 5 minutes, or 20 concurrent requests until 500 requests have completed.  This can be useful to test circuit breaking, rate limiting, and other infrastructure behaviors.

### Multiple Response Matching
You can configure multiple responses with the same status code for a single request.  Rapid will try each matching response in order and succeed on the first one that fully validates.  This is useful when a server may return the same status code with different content depending on conditions (e.g., different backends behind a load balancer).

### Once Only
A request marked `once_only: true` will execute during the first iteration only.  On subsequent iterations it is skipped entirely (including its thundering herd configuration).  This is useful for setup requests like authentication that should not repeat.

### Prometheus Metrics
When configured, Rapid collects Prometheus metrics and pushes them to a [Prometheus PushGateway](https://prometheus.io/docs/instrumenting/pushing/) after the scenario completes.  The metrics follow the [RED](https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/) (Requests, Errors, Durations) paradigm with Prometheus counters for request and error counts, and a histogram for request durations.

To disable metrics gathering, omit the `prometheus_configuration` section entirely.

### Graceful Cancellation
Rapid handles SIGINT (Ctrl-C) gracefully, cancelling in-flight requests and stopping cleanly rather than terminating abruptly.  Statistics for completed requests are still printed.

## Statistics and Metrics

Rapid prints statistics at normal termination containing counts and timings for both requests and responses. A typical output:

```
INF count=10 errors=0 minTime=2.83ms maxTime=102.05ms avgTime=56.59ms request.name=get-users request.method=get
INF count=10 errors=0 minTime=2.83ms maxTime=102.05ms avgTime=56.59ms request.name=get-users request.method=get response.name=success response.status=200
```

The first line shows request-level totals.  The second line shows the breakdown per response.

Errors are counted when:
- A network or connection error occurs
- An expected header is missing or has the wrong value
- An expected cookie is not present in the response
- Response content fails a `contains` regex check
- Content-Type doesn't match the expected type
- The response status code doesn't match any configured response (tracked as an "unconfigured" response)

Note that Rapid checks for the *presence* of expected values.  Extra headers or cookies returned by the server that are not in your configuration are not flagged as errors.

## Contributing
Contributions, bug reports, suggestions are welcome.  Please note that pull requests must pass both [staticcheck](https://github.com/dominikh/go-tools), and [revive](https://github.com/mgechev/revive) linting. Please create an issue.

## Configuration Reference

A scenario is wholly contained within a single YAML file.  A complete template of all fields can be found in [docs/template.yaml](docs/template.yaml).

### Scenario Fields

| Field | Notes| Default |Type|
|-------|---|---|--|
|name | Optional name for this scenario || string |
|version | Optional version || string |
| comment | Optional comment || string |
| request_timeout | Timeout for individual HTTP requests. Specify a duration: *ms*, *s*, *m*, or *h*. | 30s | duration |

### Find&Replace

| Field | Notes| Default| Type|
|-------|---|---|--|
|match | [RE2 regular expression](https://golang.org/s/re2syntax) to match | | string |
|replace | Replacement value for a successful match || string |

Take care to avoid collisions between the match regex and data in the fields being modified.  Rapid replaces all successful matches within a field.

### TLS Configuration

Used for all requests.  Omit entirely to disable TLS client authentication.

| Field | Notes| Default| Type|
|-------|---|---|--|
|client_cert_path | Path to client certificate file in PEM format | |string |
|client_key_path | Path to client key file in PEM format || string |
|ca_cert_path | Path to CA certificate in PEM format. If set, used instead of system certificates. | | string |
|insecure_skip_verify| Skip server certificate verification |false| boolean |

### Prometheus Configuration

Omit this section entirely to disable metrics.

| Field | Notes| Default| Type|
|-------|---|---|--|
|job_name | Prometheus job name for pushed metrics | |string |
|push_gateway_url | URL to your Prometheus PushGateway | |string |
|tls_configuration | Same fields as TLS above, specific to the push gateway || |
|buckets | Histogram bucket configuration (see below) |||
|headers | Additional headers sent with the push request || array |

#### Histogram Buckets

Rapid uses Prometheus' [ExponentialBucketsRange](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#ExponentialBucketsRange) to generate buckets.

| Field | Notes| Default| Type|
|--|--|--|--|
|minimum_bucket_duration | Minimum bucket boundary. Duration: *us*, *ms*, *s*, *m*, *h*. | 1ns | duration |
|maximum_bucket_duration | Maximum bucket boundary. Must be non-zero. | 1m | duration |
|count | Number of buckets. Prometheus adds an +Inf bucket automatically. | 5 | integer |

### Sequence

| Field | Notes| Default| Type|
|-------|---|---|---|
|iterations | Number of times to loop through all requests. Must be at least 1 to execute. |0| integer |
|iteration_time_limit | Maximum time per iteration. Duration: *s*, *m*, *h*. Zero means no limit. | 0 | duration |
|abort_on_error | Stop execution immediately when any request encounters an error | false| boolean |
|ignore_duplicate_errors | During thundering herd execution, only log each unique error message once (errors are still counted in stats) | false | boolean |
|requests | The array of request definitions || array |

### Request

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | Name for this request, used in logging and metrics || string |
|once_only | Execute only on the first iteration |false| boolean |
|method | [HTTP method](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods). Converted to uppercase. || string |
|url | Complete URL including query parameters. Passed through Find&Replace. || string |
|content | Request body. Passed through Find&Replace. || string |
|content_type | MIME type for the content. Sets the Content-Type header. || string |
|thundering_herd | Concurrent execution configuration (see below) ||  |
|extra_headers | Additional headers (see below) || array |
|cookies | Cookies to send (see below) || array |
|retry | Retry configuration for transient failures (see below) || |
|responses | Expected responses (see below) || array |

#### Retry

Controls automatic retry of HTTP requests on connection errors or specific status codes.  Retries use exponential backoff.  Omit entirely to disable retries.

| Field | Notes| Default| Type|
|-------|---|---|---|
|max_attempts | Total attempts including the initial request. Must be >= 2 to enable retries. |1| integer |
|delay | Initial delay before the first retry. Duration: *ms*, *s*, *m*. |0| duration |
|max_delay | Maximum delay cap for exponential backoff |0| duration |
|status_codes | HTTP status codes that trigger a retry (e.g., 429, 503) || array of integers |

Only connection failures and responses with a matching status code are retried.  Validation errors (wrong headers, content mismatches) are never retried.

#### Thundering Herd

Controls concurrent execution of a request within an iteration.  Omit to execute exactly one request per iteration.

| Field | Notes| Default| Type|
|-------|---|---|---|
|maximum_requests | Total requests to execute. Ignored if *time_limit* is set. |1| integer |
|concurrent_requests | Number of concurrent in-flight requests |1| integer |
|time_limit | Duration limit for the herd. If set, *maximum_requests* is ignored. |0| duration |
|delay | Delay between launching each concurrent request |0| duration |

#### Extra Headers

Header values are passed through Find&Replace.  Header names are used as-is.

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | The header name|| string |
|value | The header value || string |

#### Request Cookies

Cookie values are passed through Find&Replace.

| Field | Notes| Default| Type|
|-------|---|---|---|
|value | The cookie string (e.g., `name=value; attr=x`)|| string |

### Response

Multiple responses may share the same status_code.  Rapid tries each match in order and succeeds on the first that validates fully.

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | Name for this response, used in logs and metrics || string |
|status_code | Expected HTTP status code |0| integer |
|headers | Expected response headers (see below) || array |
|cookies | Expected response cookies (see below) || array |
|content | Content validation (see below) || |

#### Response Headers

Rapid verifies that each configured header is *present* in the response with the expected value.  Extra headers returned by the server are not flagged.  Header names are matched using Go's canonical form (e.g., "content-type" matches "Content-Type").  Values are compared exactly.

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | Header name || string |
|value | Expected value (exact match) || string |

#### Response Cookies

Rapid verifies that each configured cookie is present in the response.  Extra cookies returned by the server are not flagged.

| Field | Notes| Default| Type|
|-------|---|---|---|
|value | Expected Set-Cookie value || string |

#### Content

| Field | Notes| Default| Type|
|-------|---|---|---|
|expected | Whether the response body should contain content |false| boolean |
|content_type | Expected MIME type, verified against the Content-Type header || string |
|max_content | Maximum bytes to read from the response body for validation |4096| integer |
|contains | Array of [RE2 regular expressions](https://golang.org/s/re2syntax) that must match the content || array |
|extract | Data extraction rules (see below) || array |

#### Extract

Extracts data from the response body and registers it as a new Find&Replace entry for use in subsequent requests.  Extraction only runs after all other validations pass.

| Field | Notes| Default| Type|
|-------|---|---|---|
|type | `text`, `json`, or `xml` || string |
|path | Search path: RE2 regex for text, [GJSON](https://github.com/tidwall/gjson) path for JSON, [XPATH](https://github.com/antchfx/xmlquery) for XML || string |
|match | RE2 regex to use as the Find&Replace match key. The extracted value becomes the replacement. || string |
