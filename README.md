# REST API Diagnostic (RAPID) tool
**rapid** can be used to both verify conformance against your API spec, as well as measure load and/or performance behavior against your REST server.

To use RAPID, you create YAML configurations called a *scenario*, and feed that configuration into **rapid**.

```bash
% rapid run -s [path to scenario]
```

A RAPID *scenario* consists of one of more *requests* along with their expected *responses*.  Sequences can contain iteration counts or execute in a loop for a specific period of time. 

You define content, headers, and cookies for both requests and the expected responses.  **rapid* compares the actual response data with the expected response configuration and informs you of any discrepencies.

## Scenario
A scenario is the basic unit that describes a test case for RAPID.  Scenarios consist of one or more *sequences* of *requests* and their expected *responses*.. 

Simply put, a request is a definition of a REST API instance.  You can define the parameters of the URL and query parameters, and any additional headers for the request.

## Features

another test.
