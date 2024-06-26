# REST API Diagnostic (RAPID) tool
RAPID is a REST API diagnostic tool you can use to both verify conformance against your API spec, as well as measure load and/or performance behavior against your REST server.  It can be used as a CLI command, or installed as a server instance that exports a set of, well, REST APIs for managing testcases.

To use RAPID, you create YAML configurations that describe a *scenario*, and feed that configuration into RAPID.  A RAPID *scenario* consists of *operation* definitions and *sequences* of those operations.  Sequences can contain iteration counts or execute for a specific period of time. 


RAPID collects and maintains metrics about each scenario and these can be queried from the server instance.  In the CLI instance, these metrics are spewed to stdout when the scenario has completed.

## Scenario
A scenario is the basic unit that describes a test case for RAPID.  Scenarios consist of *operations* and *sequences*. 

Simply put, an operation is a definition of a REST API instance.  You can define the parameters of the URL and, as needed, query parameters, and headers for the operation.  Operations are required to have a unique name, and these names are used to create sequences.

Sequences, as the name implies, is an ordered set of defined operations.  

For example, you could define a set of operations that includes DELETE, GET, and PUT HTTP methods.  Then define a scenario that orders them as:

PUT, GET, DELETE
