package lisp

import (
	"testing"
)

func TestObjectString(t *testing.T) {
	tests := [...]struct {
		object Object
		want   string
	}{
		// Explicit indexes are useful when locating failed test.

		0: {NewInt(0), "0"},
		1: {NewInt(64), "64"},
		2: {NewInt(-64), "-64"},

		3: {NewFloat(0), "0.0"},
		4: {NewFloat(64), "64.0"},
		5: {NewFloat(-64), "-64.0"},
		6: {NewFloat(0.55), "0.55"},
		7: {NewFloat(64.55), "64.55"},
		8: {NewFloat(-64.55), "-64.55"},

		9:  {NewSymbol("nil"), "nil"},
		10: {NewSymbol(""), "##"},
		11: {NewSymbol("symbol-name"), "symbol-name"},

		12: {NewVector(nil), "[]"},
		13: {
			NewVector([]Object{NewInt(1), NewSymbol("foo")}),
			"[1 foo]",
		},
		14: {
			NewVector([]Object{
				NewFloat(0.4),
				NewVector([]Object{NewFloat(1)}),
				NewVector(nil),
				NewFloat(0.3),
			}),
			"[0.4 [1.0] [] 0.3]",
		},
		15: {
			NewVector([]Object{NewVector([]Object{NewVector(nil)})}),
			"[[[]]]",
		},
	}

	for i, tt := range tests {
		have := ObjectString(tt.object)
		if have != tt.want {
			t.Errorf("test %d:\nwant: `%s`\nhave: `%s`",
				i, tt.want, have)
		}
	}
}
