package radix_tree

import (
	"github.com/magiconair/properties/assert"
	"log"
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

	for _, v := range r.Find("abc") {
		log.Println(v)
	}

	assertEqual(t,
		r.Find("abc"),
		[]string{"abc", "*bc", "*"},
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
