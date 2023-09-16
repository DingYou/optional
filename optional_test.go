package optional

import (
	"testing"
)

func TestOptional_String(t *testing.T) {
	// a := make([]string, 0, 16)
	var a []string
	t.Log(a == nil)
	o := OfNullable(a)
	t.Log(o.String())

	i := Of(1)
	t.Log(i.String())
}
