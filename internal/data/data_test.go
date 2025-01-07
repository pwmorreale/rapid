//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package data_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/data"
	"github.com/test-go/testify/assert"
)

func TestAdd(t *testing.T) {

	d := data.New()

	err := d.Add("goo", "kkk")
	assert.Nil(t, err)
	assert.Equal(t, len(d.All), 1)
}

func TestReplace(t *testing.T) {

	d := data.New()

	err := d.Add("goo", "value")
	assert.Nil(t, err)
	assert.Equal(t, len(d.All), 1)

	s := d.Replace("$goo")
	assert.Equal(t, "value", s)

	s = d.Replace("moo$goo")
	assert.Equal(t, "moovalue", s)

	s = d.Replace("$gooboo")
	assert.Equal(t, "valueboo", s)

	s = d.Replace(`\$goo`)
	assert.Equal(t, `\value`, s)

	s = d.Replace(`\$foo`)
	assert.Equal(t, `\$foo`, s)

	s = d.Replace(`$foo`)
	assert.Equal(t, `$foo`, s)
}

func TestMultiReplace(t *testing.T) {

	d := data.New()

	err := d.Add("hello", "hi")
	assert.Nil(t, err)

	err = d.Add("George", "Hank")
	assert.Nil(t, err)

	err = d.Add("Steve", "Fred")
	assert.Nil(t, err)

	assert.Equal(t, len(d.All), 3)

	before := "$hello $George, this is $Steve paying you $50"
	after := "hi Hank, this is Fred paying you $50"

	s := d.Replace(before)
	assert.Equal(t, after, s)

}
