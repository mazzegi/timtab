package bitset

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, want, have any) {
	if reflect.DeepEqual(want, have) {
		return
	}
	t.Fatalf("want %v, have %v", want, have)
}

func TestBitset(t *testing.T) {

	bs := New(8)
	bs.Set(2, true)
	bs.Set(4, true)
	bs.Set(6, true)
	assertEqual(t, false, bs.Get(1))
	assertEqual(t, true, bs.Get(2))
	assertEqual(t, false, bs.Get(3))
	assertEqual(t, true, bs.Get(4))
	assertEqual(t, false, bs.Get(5))
	assertEqual(t, true, bs.Get(6))

	bs.Set(4, false)
	assertEqual(t, false, bs.Get(1))
	assertEqual(t, true, bs.Get(2))
	assertEqual(t, false, bs.Get(3))
	assertEqual(t, false, bs.Get(4))
	assertEqual(t, false, bs.Get(5))
	assertEqual(t, true, bs.Get(6))
}
