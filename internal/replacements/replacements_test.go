//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

package replacements_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/replacements"
	"github.com/stretchr/testify/assert"
)

func TestAddReplacement(t *testing.T) {

	d := replacements.New()

	err := d.AddReplacement("goo", "kkk")
	assert.Nil(t, err)
	assert.Equal(t, d.Len(), 1)
}

func TestAddReplacementBad(t *testing.T) {

	d := replacements.New()

	err := d.AddReplacement(`\((?!['"]`, "kkk")
	assert.Equal(t, "error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err.Error())
	assert.Equal(t, d.Len(), 0)
}

func TestReplace(t *testing.T) {

	d := replacements.New()

	err := d.AddReplacement("goo", "value")
	assert.Nil(t, err)
	assert.Equal(t, d.Len(), 1)

	s := d.Replace("goo")
	assert.Equal(t, "value", s)

	s = d.Replace("moogoo")
	assert.Equal(t, "moovalue", s)

	s = d.Replace("gooboo")
	assert.Equal(t, "valueboo", s)

	s = d.Replace(`\goo`)
	assert.Equal(t, `\value`, s)

	s = d.Replace(`\foo`)
	assert.Equal(t, `\foo`, s)

	s = d.Replace(`$foo`)
	assert.Equal(t, `$foo`, s)
}

func TestMultiReplace(t *testing.T) {

	d := replacements.New()

	err := d.AddReplacement("hello", "hi")
	assert.Nil(t, err)

	err = d.AddReplacement("George", "Hank")
	assert.Nil(t, err)

	err = d.AddReplacement("Steve", "Fred")
	assert.Nil(t, err)

	assert.Equal(t, d.Len(), 3)

	before := "hello George, this is Steve paying you $50"
	after := "hi Hank, this is Fred paying you $50"

	s := d.Replace(before)
	assert.Equal(t, after, s)

}
