//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package sequences_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/sequences"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/test-go/testify/assert"
)

func TestNew(t *testing.T) {
	sc := &mocks.FakeScenario{}
	rpt := &mocks.FakeReport{}

	seq := sequences.New(sc, rpt)
	assert.NotNil(t, seq)
}

func TestGenerate(t *testing.T) {
	sc := &mocks.FakeScenario{}
	rpt := &mocks.FakeReport{}

	seq := sequences.New(sc, rpt)
	assert.NotNil(t, seq)

	err := seq.Run()
	assert.Nil(t, err)
}
