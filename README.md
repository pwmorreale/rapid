# REST API Diagnostic (RAPID) tool

## *RAPID is under construction and is incomplete at this time.*



Rapid is a tool used to both verify conformance of your REST server against your API spec, as well as measure load and/or performance behavior.  Rapid also makes it possible for you to test policies such as circuit-breaking, rate-limiting, load-balancing, etc. 

Rapid works entirely through a YAML configuration file.   The configuration is called a *scenario* and consists of a sequence of one of more http/https requests along with their possible responses.  Sequences can contain iteration counts or execute in a loop for a specific period of time.  You can also easily configure multiple concurrent requests to load your infrastructure.

You define content, headers, and cookies for both requests and the responses.  Rapid compares the actual response data with the expected response configuration and informs you of any discrepancies.  

In addition, Rapid allows you to dynamically extract data from previous responses and insert that data into future requests.  This allows you to create dynamic paths through your service infrastructure.

## Install
To install:

```bash
% go install github.com/pwmorreale/rapid@latest
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

## Features

Rapid provides several features:

### Iterations
You can define an iteration count and iteration time limit.  Each iteration will loop through all configured in order and each iteration must complete within the specified time limit.

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
      method:
      url:
      content:
      content_type:
      thundering_herd:
        maximum_requests:
        concurrent_requests:
        time_limit:
        delay:
      extra_headers:
        - name:
          value:
      cookies:
        - value:
      responses:
        - name:
          status_code:
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
### Scenario fields
The *name*, *version* and *comment* fields are optional and if present will appear in the report.  Use these fields to identify the test run, server instance, etc.

```yaml
name:
version:
comment:
```

| Field | Notes| Default |Type|
|-------|---|---|--|
|name | Optional name for this  scenario || string |
|version | Optional version for the scenario, or whatever you choose | |string |
| comment | Optional comment || string |

### Find&Replace Configuration
Find&Replace allows you to define fields for replacement during execution of the scenario.  Headers (names and values), cookies, and URLs  are passed through this module for expansion prior to being referenced.  

Use https://golang.org/s/re2syntax for the *match* regular expression.  **Note you must take care to avoid any collisions between the match string and data within the field being modified**  Rapid will indescriminately replace all successful matches within the field.

Also see the *extract* configuration below.

```yaml
find_replace:
  match:
  replace:
```

 |Field | Notes| Default| Type|
|-------|---|---|--|
|match | Regex2 to match | | string |
|replace | Replacement value for a successful match || string |

### TLS Configuration
This section allows you to specify the certificates used for TLS connections.  This configuration is used for all requests in the configuration. 
If this configuration is omitted, then TLS on the client will not be enabled.

If you specify a CA certificate, it will be used to verify the identify of the server certificate during the TLS handshake and the system certificates will be ignored.

The *insecure_skip_verify* field tells the client to ignore attempts to verify the server certificate.  This can be useful when you do not have a CA certificate from the server.

```yaml
tls_configuration:
  client_cert_path:
  client_key_path:
  ca_cert_path:
  insecure_skip_verify:
```

 |Field | Notes| Default| Type|
|-------|---|---|--|
|client_cert_path | Path to client certificate file in PEM format. | |string |
|client_key_path |Path to client key certificate file in PEM format|| string |
|ca_cert_path | Path to CA certificate file in PEM format.  | | string |
|insecure_skip_verify| If set to *true*, then the client will not attempt to verify the server certificate. |false| boolean |

### Sequence Configuration
The *sequence* section defines iterations of the *requests*. 
```yaml
sequence:
  iterations:
  iteration_time_limit:
  abort_on_error:
  ignore_duplicate_errors:
  requests:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|iterations | The number of times to iterate through the array of requests.  Specify at least '1' to execute the requests. |0| integer |
|iteration_time_limit | The maximum amount of time to allow an iteration to complete. Specify an integer with a modifier of *s* (seconds), *m* (minutes), or *h* (hours) | 0 | string |
| abort_on_error | currently unimplemented | false| boolean |
| ignore_duplicate_errors | currently unimplemented | false | boolean |
| requests| The array of requests, see next section| | array |


### Request Configuration
The *requests* array defines the *requests*.
```yaml
  requests:
    - name:
      once_only:
      method:
      url:
      content:
      content_type:
      thundering_herd:
        maximum_requests:
        concurrent_requests:
        time_limit:
        delay:
      extra_headers:
        - name:
          value:
      cookies:
        - value:
      responses:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | The name for this request.  Used in logging and reports|| string |
|once_only | Execute this request exactly one time, regardless of the *iterations* count, and/or the subsequent *thundering_herd* configuration. |false| boolean |
|method  | The [HTTP method](http://www.w3schools.com/TAgs/ref_httpmethods.asp).  This will be converted to upper case.|| string |
|url | The complete URL of the request.  Include and query parameters, fragments, etc.  The URL is only passed to the **Find&Replace** module for modification.  No other modifications are made. || string |
|content |Content for this request.  If present, the contents of this parameter becomes the body of the request. This parameter is passed through the **Find&Replace** module for modification. Note that Rapid will automatically set the *Content-Length* header. || string |
|content_type | The [mime type](https://developer.mozilla.org/en-US/docs/Glossary/MIME_type) of the *content* parameter.  If defined, a *Content-Type* header will be added to the request with this value. || string |
|thundering_herd | See next section||  |

#### Thundering Herd Configuration
The *thundering_herd* configuration allows you to control concurrent execution of this request within the current iteration.

If this configuration is omitted, then exactly one instance of the request will be executed for each iteration.  Also be aware that regardless of this configuration for any of the requests in the array, the iteration time limit will be enforced.

Note that Rapid always creates a separate client for each request.  This means that, by design, rapid will use a separate system socket for each transaction with the implied server.  Consequently, the actual number of concurrent requests that can be created can be limited by the platform used to execute Rapid.  The intent here is to mimic real world application where a HTTP server is handling multiple concurrent requests from multiple separate clients.  

```yaml
thundering_herd:
  maximum_requests:
  concurrent_requests:
  time_limit:
  delay:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|maximum_requests |The total number of requests to execute.  Note that *maximum_requests* is ignored if *time_limit* is configured |1| integer |
|concurrent_requests |The number of requests to execute concurrently.  Rapid will maintain this number of active requests until one of the limits is reached. |1| integer |
|time_limit | The time limit for the herd. Specify an integer with a modifier of *s* (seconds), *m* (minutes), *h* (hours).  If this parameter is set, then *maximum_requests* is ignored.|0| string |
|delay | Optional delay between starting each of the *concurrent_requests*. Specify an integer with a modifier of *ms* (milliseconds), *s* (minutes) or  *m* (minutes). |0| integer |

#### Extra Headers
This section allows you to specify an array of additional headers to the request.  Header names and values are passed through the **Find&Replace** module for expansion prior to being added to the request.

```yaml
extra_headers:
  - name:
    value:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | The header name|| string |
|value |The header content|| string |

#### Cookies Configuration
The *cookies* section allows you to define an array of cookies sent with the request.  An appropriate *Cookie* header will be constructed and sent with the request.

```yaml
cookies:
  - value:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|value | The cookie string|| string |

### Response Configuration
The *responses* section of *request* configuration is an array of the possible responses to this request. Here you configure the expected data associated with a particular response.  Rapid will compare the actual response against what is expected and report any discrepencies.

```yaml
responses:
  - name:
    status_code:
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

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | The name for this response.  Used in logs and reports|| string |
|status_code | The HTTP status code for this response.|0| integer |
|headers | See next section|| array |

#### Headers
These are headers expected to be set by the server in the reponse.  Rapid will perform exact, case sensitive matches and report any discrepencies.

```yaml
headers:
  - name:
    value:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|name | The case-sensitive name for this header.|| string |
|value | The expected case-sensitive contents of this header.|| string |

#### Cookies
The array of Cookies expected to be returned from the server.  Rapid will parse the cookie(s) and compare both the cookie and its attributes against those returned by the server.

```yaml
cookies:
  - value:
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|value | The expected cookie|| string |

#### Content
This section allows you to verify content returned from the server.   
```yaml
content:
  expected:
  content_type:
  max_content:
  contains:
    - ""
```

| Field | Notes| Default| Type|
|-------|---|---|---|
|expected | Whether the body of the response should contain content. This can be useful to check for data leaks from the server. |false| boolean |
|content_type | The expected [mime type](https://developer.mozilla.org/en-US/docs/Glossary/MIME_type) of the content.  If defined, the value is verified against the *Content-Type* header || string |
|max_content | A maximum length of the expected content.  The actual content length (if it can be detected) will be compared to this value if defined |false| integer |
|contains | An array of [regular expressions](https://golang.org/s/re2syntax) used to verify the content. || string |

#### Extract
Use this configuration to extract data from the response body and insert it into the **Find&Replace** module for future use in headers/URLs/bodies/etc. 

Rapid supports extraction from three different types of data:  text, JSON, or XML.  Note that the type you specify in the cnfiguration below is independent of the MIME type present in the *Content-Type* header.  Thus you can mix and match extraction types from the same response body.

```yaml
  extract:
    - type:
      path:
      match:
```
| Field | Notes| Default| Type|
|-------|---|---|---|
|type | Specify 'text', 'json', or 'xml'. || string |
|path | The path to search.  If *type* is 'text', then specify a [regular expression](https://golang.org/s/re2syntax). If *type* is 'json', then specify a [GJSON](https://github.com/tidwall/gjson) path.  If *type* is 'xml', then specify a suitable [XPATH](https://github.com/antchfx/xmlquery).   || string |
|match | A **Find&Replace** *match* [regular expression](https://golang.org/s/re2syntax) for future modifications in headers/URLs/bodies/etc.  The extracted data will be used as the replacement text on matches to this expression.|| string |

