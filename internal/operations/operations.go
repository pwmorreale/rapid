//
//  Copyright © 2023 Peter W. Morreale. All RIghts Reserved.
//

// Package operations defines http/s operations
package operations

// Definition defines a URL request.
type Definition struct {
	Scheme       string              `yaml:"scheme"`        // http/s
	Token        string              `yaml:"bearer_token"`  // Bearer token
	UserName     string              `yaml:"user_name"`     // simple auth user
	UserPassword string              `yaml:"user_password"` // simple auth password.
	Host         string              `yaml:"host"`          // host or host:port
	Path         string              `yaml:"path"`          // path (relative paths may omit leading slash)
	Fragment     string              `yaml:"fragment"`      // fragment for references, without '#'
	Values       map[string][]string `yaml:query"`          // Query name/value pairs
}
