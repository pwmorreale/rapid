# REST API Diagnostic (RAPID) tool
Rapid is a tool used to both verify conformance of your REST server against your API spec, as well as measure load and/or performance behavior.

You create YAML configurations called a *scenario*, and feed that configuration into Rapid.  A *scenario* consists of a sequence of one of more http/https *requests* along with their possible *responses*.  Sequences can contain iteration counts or execute in a loop for a specific period of time. 

You define content, headers, and cookies for both requests and the responses.  Rapid compares the actual response data with the expected response configuration and informs you of any discrepancies.  

Rapid makes it easy to test circuit breakers, rate limiting and other behaviors by allowing you to configure a full range of responses.

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

### *find&replace*
Find&replace allows you to predefine a set of regex terms and their associated replacement strings.  When a regex matches in a header value, cookie, or the URL, the replacement term is inserted in its place.  This allows you to define a term once, and have it referenced throughout the entire configuration.

### Data Extraction
Rapid also allows you to extract data from response payloads for use in future requests.  You can search through JSON, XML, or text responses and have the data saved to the find&replace module.

This allows for example, extraction of a security token from an authorization request response for use in **Authorization** headers in future requests.  Another possibility would be to use returned response data to modify the URL of a future request in the configuration.

### Thundering Herd 
Rapid allows you to create *thundering herd* configurations that allow you to specify a number of concurrent requests for a specific duration of time, or a maximum number of requests.  For example, you could configure Rapid to execute 1000 requests concurrently for 5 minutes, or 20 concurrent requests until 500 requests have completed.

## Configuration
A scenario is the basic unit that describes a test case for RAPID.  A scenario is wholly contained within a single YAML file.  

Scenarios consist of a *sequence* of one or more *requests* and their expected *responses*.

A request is a definition of a REST API instance.  You can define the parameters of the URL and query parameters, and any additional headers, cookies and/or payload for the request and define the possible responses.



