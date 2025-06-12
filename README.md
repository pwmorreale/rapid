# REST API Diagnostic (RAPID) tool

## *RAPID is under construction and is incomplete at this time.*



Rapid is a tool used to both verify conformance of your REST server against your API spec, as well as measure load and/or performance behavior.  Rapid also makes it possible for you to test policies such as circuit-breaking, rate-limiting, load-balancing, etc. 

Rapid works entirely through a YAML configuration file.   The configuration is called a *scenario* and consists of a sequence of one of more http/https requests along with their possible responses.  Sequences can contain iteration counts or execute in a loop for a specific period of time.  You can also easily configure multiple concurrent requests to load your infrastructure.

You define content, headers, and cookies for both requests and the responses.  Rapid compares the actual response data with the expected response configuration and informs you of any discrepancies.  

## Install
To install:

```bash
% go install github.com/pwmorreale/rapid@latest
```

## Usage
To execute a scenario, use the ***run*** command:

```bash
% rapid -s ./scenario.yaml
```
You can also check a scenario configuration to find common typos/etc by using the ***verify*** command:

```bash
% rapid verify -s ./scenario.yaml
```

## Features

Rapid provides several features.

### *Find&Replace*
Find&replace allows you to predefine a set of regex terms and their associated replacement strings.  When a regex matches in a header value, cookie, or the URL, the replacement term is inserted in its place.  This allows you to define a term once, and have it referenced throughout the entire configuration.

### *Data Extraction*
Rapid also allows you to extract data from response payloads for use in future requests.  You can search through JSON, XML, or text responses and have the data saved to the **find&replace** module.

This allows for example, extraction of a security token from an authorization response body for use in **Authorization** headers in future requests.  Another possibility would be to use returned response data to modify the URL of a future request in the sequence.

### *Thundering Herd* 
Rapid allows you to create *thundering herd* configurations that allow you to specify a number of concurrent requests for a specific duration of time, or a maximum number of requests.  For example, you could configure Rapid to execute 1000 requests concurrently for 5 minutes, or 20 concurrent requests until 500 requests have completed.  This can be useful to test circuit breaking, rate limiting, and other infrastructure behaviors.


## Configuration
A scenario is the basic unit that describes a test case for RAPID.  A scenario is wholly contained within a single YAML file.  Scenarios consist of a *sequence* of one or more *requests* and their expected *responses*.

A complete list of scenario fields can be found in the docs/template file.  A listing and discussion of the fields follows:
```yaml
name:
version:
comment:
find_replace:
  - match:
    replace:
tls_configuration:
  client_cert_path:
  client_key_path:
  ca_cert_path:
  insecure_skip_verify:
sequence:
  iterations:
  iteration_time_limit:
  abort_on_error:
  ignore_duplicate_errors:
  requests:
    - name:
      once_only:
      thundering_herd:
        maximum_requests:
        concurrent_requests:
        time_limit:
        delay:
      method:
      url:
      extra_headers:
        - name:
          value:
      cookies:
        - value:
      content:
      content_type:
      responses:
        - status_code:
          name:
          headers:
            - name:
              value:
          cookies:
            - value:
          content:
            expected:
            content_type:
            max_content:
            contains:
              - ""
            extract:
              - type:
                path:
                match:
```
## Scenario fields

```yaml
name:
version:
comment:
```

| Field | Notes| Type|
|-------|---|---|
|name | Optional name for this  scenario | string |
|version | Optional version for the scenario, or whatever you choose | string |
| comment | Optional comment | string |

The *name*, *version* and *comment* fields are optional and if present will appear in the report.  Use these fields to identify the test run, server instance, etc.

## Find&Replace fields

```yaml
find_replace:
  match:
  replace:
```

 |Field | Notes| Type|
|-------|---|---|
|match | Regex2 to match | string |
|replace | Replacement value for a successful match | string |

FInd&Replace allows you to define fields for replacement during execution of the scenario.  Headers (names and values), cookies, and URLs  are passed through this module for expansion prior to being referenced.  

Use https://golang.org/s/re2syntax for the *match* regular expression.  **Note you must take care to avoid any collisions between the match string and data within the field being modified**  Rapid will indescriminately replace all successful matches within the field.

