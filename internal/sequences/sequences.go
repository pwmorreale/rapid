//
//  Copyright © 2023 Peter W. Morreale. All Rights Reserved.
//

// Package sequences defines a sequence of RAPID operations
package sequences

import "time"

// Definition defines a sequence
type Definition struct {
	Name        string    `yaml:"scheme"name`           // Sequence Name.`
	Id          string    `yaml:id"`                    // Sequence Identifier.
	RepeatCount int32     `yaml:"repeat_count"`         // repeat count
	Duration    time.Time `yaml:"repeat_time_duration"` // repeat time limit
}
