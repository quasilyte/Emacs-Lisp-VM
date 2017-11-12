package lisp

import (
	"math/rand"
	"testing"
)

const getsetRepeats = 64

func TestObjectInt(t *testing.T) {
	for i := 0; i < getsetRepeats; i++ {
		x1 := int64(rand.Int())
		o := NewInt(x1)
		if o.Int() != x1 {
			t.Fatalf("GetInt():\nwant: %v\nhave: %v",
				x1, o.Int())
		}
		x2 := int64(rand.Int())
		o.SetInt(x2)
		if o.Int() != x2 {
			t.Fatalf("GetInt():\nwant: %v\nhave: %v",
				x2, o.Int())
		}
	}
}

func TestObjectFloat(t *testing.T) {
	for i := 0; i < getsetRepeats; i++ {
		x1 := rand.Float64()
		o := NewFloat(x1)
		if o.Float() != x1 {
			t.Fatalf("GetFloat():\nwant: %v\nhave: %v",
				x1, o.Float())
		}
		x2 := rand.Float64()
		o.SetFloat(x2)
		if o.Float() != x2 {
			t.Fatalf("GetFloat():\nwant: %v\nhave: %v",
				x2, o.Float())
		}
	}
}

func TestNull(t *testing.T) {
	tests := [...]struct {
		object Object
		want   bool
	}{
		{Nil, true},
		{T, false},
		{NewSymbol("nil"), false},
		{NewFloat(4), false},
		{NewInt(4), false},
	}

	for _, tt := range tests {
		have := Null(&tt.object)
		if have != tt.want {
			t.Errorf("%s is not nil", ObjectString(tt.object))
		}
	}
}
