package radix_tree

import (
	"github.com/magiconair/properties/assert"
	"slices"
	"testing"
)

func assertEqual(t *testing.T, got, want []string) {
	slices.Sort(got)
	slices.Sort(want)

	assert.Equal(t, got, want)
}

func TestRadix(t *testing.T) {
	r := NewRadix()
	r.Add("abc")
	r.Add("*cd")
	r.Add("*bc")
	r.Add("b*d")
	r.Add("*")

	assertEqual(t,
		r.Find("abc"),
		[]string{"abc", "*bc", "*"},
	)
}

func TestRadix2(t *testing.T) {
	r := NewRadix()
	r.Add("a*c")
	r.Add("*")

	assertEqual(t,
		r.Find("abbc"),
		[]string{"a*c", "*"},
	)
}

func TestRadix3(t *testing.T) {
	r := NewRadix()
	r.Add("a?c")
	r.Add("*")

	assertEqual(t,
		r.Find("abbc"),
		[]string{"*"},
	)
}

func TestRadix4(t *testing.T) {
	r := NewRadix()
	r.Add("a?c")
	r.Add("*")

	assertEqual(t,
		r.Find("abc"),
		[]string{"a?c", "*"},
	)
}

func TestDAWG(t *testing.T) {
	r := NewRadix()
	r.Add("abc")
	r.Add("*cd")
	r.Add("*bc")
	r.Add("b*d")
	r.Add("*")

	assertEqual(t,
		r.ToDAWG().Find("abc"),
		[]string{"abc", "*bc", "*"},
	)
}
