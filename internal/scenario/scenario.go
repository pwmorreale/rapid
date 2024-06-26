//
//  Copyright © 2023 Peter W. Morreale. All RIghts Reserved.
//

// Package scenario defines a complete testing scenario.
package scenario

import (
	"github.com/pwmorreale/rapid/internal/operations"
	"github.com/pwmorreale/rapid/internal/sequences"
)

// Instance defines a scenario instance.
type Instance struct {
	Name      string                  `yaml:"name"`
	Id        string                  `yaml:"id"`
	Ops       []operations.Definition `yaml:"operations"`
	Sequences []sequences.Definition  `yaml:"sequences"`
}
