//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package reporter_test

import (
	"os"
	"testing"

	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/test-go/testify/assert"
)

func TestNew(t *testing.T) {
	sc := &mocks.FakeScenario{}

	rpt := reporter.New(sc)
	assert.NotNil(t, rpt)
}

func TestGenerate(t *testing.T) {
	sc := &mocks.FakeScenario{}

	rpt := reporter.New(sc)
	assert.NotNil(t, rpt)

	err := rpt.Generate(os.Stdout)
	assert.Nil(t, err)
}
