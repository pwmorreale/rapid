//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package data_test

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/pwmorreale/rapid/data"
	"github.com/stretchr/testify/assert"
)

func TestAddReplacement(t *testing.T) {

	d := data.New()

	err := d.AddReplacement("goo", "kkk")
	assert.Nil(t, err)
	assert.Equal(t, d.Len(), 1)
}

func TestAddReplacementBad(t *testing.T) {

	d := data.New()

	err := d.AddReplacement(`\((?!['"]`, "kkk")
	assert.Equal(t, "error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err.Error())
	assert.Equal(t, d.Len(), 0)
}

func TestReplace(t *testing.T) {

	d := data.New()

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

	s = d.Replace(`foo`)
	assert.Equal(t, `foo`, s)
}

func TestMultiReplace(t *testing.T) {

	d := data.New()

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

func TestLookup(t *testing.T) {

	d := data.New()

	err := d.AddReplacement("hello", "hi")
	assert.Nil(t, err)

	err = d.AddReplacement("George", "Hank")
	assert.Nil(t, err)

	found := d.Lookup("George")
	assert.Equal(t, "Hank", found)

	found = d.Lookup("Boof")
	assert.Equal(t, "", found)
}

func TestXML(t *testing.T) {

	d := data.New()

	x := `<breakfast_menu>
                <food>
                  <name>Belgian Waffles</name>
                  <price>$5.95</price>
                  <description>
                    Two of our famous Belgian Waffles with plenty of real maple syrup
                  </description>
                  <calories>650</calories>
                </food>
                <food>
                  <name>Strawberry Belgian Waffles</name>
                  <price>$7.95</price>
                  <description>
                    Light Belgian waffles covered with strawberries and whipped cream
                  </description>
                  <calories>900</calories>
                </food>
              </breakfast_menu>`

	// Find first, then second 'calories'
	v, err := d.ExtractXML("//food/calories", strings.NewReader(x))
	assert.Nil(t, err)
	assert.Equal(t, "650", v)

	v, err = d.ExtractXML("//food[2]/calories", strings.NewReader(x))
	assert.Nil(t, err)
	assert.Equal(t, "900", v)

	// Valid, but missing node...
	v, err = d.ExtractXML("//boof/calories", strings.NewReader(x))
	assert.NotNil(t, err)
	assert.Equal(t, "", v)

	// Invalid path
	v, err = d.ExtractXML("/>?boof/calories", strings.NewReader(x))
	assert.Equal(t, "/>?boof/calories has an invalid token.", err.Error())
	assert.Equal(t, "", v)

}

func TestXMLBadReader(t *testing.T) {

	d := data.New()

	r := iotest.ErrReader(errors.New("blowing chunks"))

	v, err := d.ExtractXML("//food/calories", r)
	assert.Equal(t, "blowing chunks", err.Error())
	assert.Equal(t, "", v)
}

func TestExtractJSON(t *testing.T) {

	d := data.New()

	s := `{ "color":"blue"}`
	v, err := d.ExtractJSON("color", strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, "blue", v)
}

func TestExtractJSONBadReader(t *testing.T) {

	d := data.New()

	r := iotest.ErrReader(errors.New("blowing chunks"))
	v, err := d.ExtractJSON("color", r)
	assert.Equal(t, "JSON: Read error: blowing chunks", err.Error())
	assert.Equal(t, "", v)
}

func TestExtractJSONNotFoundError(t *testing.T) {

	d := data.New()

	s := `{ "color":"blue"}`
	v, err := d.ExtractJSON("foobar", strings.NewReader(s))
	assert.Equal(t, "JSON: Not found: foobar", err.Error())
	assert.Equal(t, "", v)
}

func TestExtractRegexp(t *testing.T) {

	d := data.New()

	s := `{ "color":"blue"}`
	v, err := d.ExtractRegex("color", strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, "color", v)
}

func TestExtractRegexpBadReader(t *testing.T) {

	d := data.New()

	r := iotest.ErrReader(errors.New("blowing chunks"))
	v, err := d.ExtractRegex("color", r)
	assert.Equal(t, "REGEX: Read error: blowing chunks", err.Error())
	assert.Equal(t, "", v)
}

func TestExtractRegexpBadCompile(t *testing.T) {

	d := data.New()

	s := `{ "color":"blue"}`
	v, err := d.ExtractRegex(`\((?!['"]`, strings.NewReader(s))
	assert.Equal(t, "error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err.Error())
	assert.Equal(t, "", v)
}

func TestExtractRegexpNotFoundError(t *testing.T) {

	d := data.New()

	s := `{ "color":"blue"}`
	v, err := d.ExtractRegex("foobar", strings.NewReader(s))
	assert.Equal(t, "REGEX: value not found for expression: foobar", err.Error())
	assert.Equal(t, "", v)
}
