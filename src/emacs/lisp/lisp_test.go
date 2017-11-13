package lisp

import (
	"math/rand"
	"testing"
)

const (
	getsetRepeats = 64
	sumCount      = 32
)

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

	xs := make([]Object, sumCount)
	for i := range xs {
		xs[i].Type = TypeInt
		xs[i].SetInt(int64(i))
	}
	sumHave := int64(0)
	sumWant := int64(0)
	for i := range xs {
		sumHave += xs[i].Int()
		sumWant += int64(i)
	}
	if sumHave != sumWant {
		t.Fatalf("sum 0..%d:\nhave: %d\nwant: %d",
			sumCount, sumHave, sumWant)
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

	xs := make([]Object, sumCount)
	for i := range xs {
		xs[i].Type = TypeFloat
		xs[i].SetFloat(float64(i))
	}
	sumHave := float64(0)
	sumWant := float64(0)
	for i := range xs {
		sumHave += xs[i].Float()
		sumWant += float64(i)
	}
	if sumHave != sumWant {
		t.Fatalf("sum 0..%d:\nhave: %f\nwant: %f",
			sumCount, sumHave, sumWant)
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
