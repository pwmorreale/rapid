//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package reporter_test

import (
	"os"
	"testing"

	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rpt := reporter.New()
	assert.NotNil(t, rpt)
}

func TestGenerate(t *testing.T) {
	rpt := reporter.New()
	assert.NotNil(t, rpt)

	err := rpt.Generate(nil, os.Stdout)
	assert.Nil(t, err)
}
