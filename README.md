# RAPID
RAPID is a REST apI diagnostic tool you can use to check REST servers for both conformance to your API spec, as well as other types of load/performance testing.

RAPID takes one or more *testcase* configurations and executes them either serially, or concurrently.  RAPID can be used as a CLI command, or installed as a server instance that exports a set of APIs for managing *testcase* scenarios.

RAPID also collects and maintains metrics about each test case and these can be queried from the server instance. 

## TestCases
A testcase is the basic unit that encapsulates a testing scenario.

You first define *operations*, then define one or more *tests* referencing those operations.

For example, you could define a set of operations that includes DELETE, GET, and PUT HTTP methods.  Then define a test that orders them as:

PUT, GET, DELETE