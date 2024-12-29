//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package sequences_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/sequences"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rpt := &mocks.FakeReport{}

	seq := sequences.New(rpt)
	assert.NotNil(t, seq)
}

func TestRun(t *testing.T) {
	rpt := &mocks.FakeReport{}

	seq := sequences.New(rpt)
	assert.NotNil(t, seq)

	err := seq.Run(nil)
	assert.Nil(t, err)
}
