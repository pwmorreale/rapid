//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package sequences_test

import (
	"fmt"
	"testing"

	"github.com/pwmorreale/rapid/internal/scenario"
	"github.com/pwmorreale/rapid/internal/sequences"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	sc := &mocks.FakeScenario{}
	rpt := &mocks.FakeReport{}

	seq := sequences.New(sc, rpt)
	assert.NotNil(t, seq)
}

func TestRun(t *testing.T) {
	sc := &mocks.FakeScenario{}
	rpt := &mocks.FakeReport{}

	seq := sequences.New(sc, rpt)
	assert.NotNil(t, seq)

	err := seq.Run()
	assert.Nil(t, err)
}

func TestUnmarshal(t *testing.T) {
	sc := scenario.New()

	err := sc.ParseFile("../../test/configs/test_scenario.yaml")
	fmt.Print(err)
	assert.Nil(t, err)
	assert.NotNil(t, sc.Viper())

	rpt := &mocks.FakeReport{}

	seq := sequences.New(sc, rpt)
	assert.NotNil(t, seq)

	s, err := seq.UnmarshalKey("sequence")
	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.Greater(t, len(s.Reqs), 0)

}
