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

		16: {NewString(nil), `""`},
		17: {NewString([]byte("")), `""`},
		18: {NewString([]byte("a b c")), `"a b c"`},
		19: {NewString([]byte(`"""`)), `"""""`},
		20: {NewString([]byte("\t\n")), "\"\t\n\""},

		21: {NewCons(Nil, Nil), "(nil . nil)"},
		22: {NewCons(NewInt(1), NewFloat(2.0)), "(1 . 2.0)"},
		23: {NewCons(NewCons(T, T), T), "((t . t) . t)"},
		24: {NewCons(T, NewCons(T, T)), "(t . (t . t))"},
		25: {NewCons(NewCons(T, T), NewCons(T, T)), "((t . t) . (t . t))"},

		26: {
			NewVector([]Object{
				NewCons(Nil, T),
				NewCons(
					NewVector([]Object{
						T,
						NewString([]byte("abc")),
						T,
					}),
					NewVector([]Object{
						Nil,
						NewInt(7),
					}),
				),
				T,
			}),
			`[(nil . t) ([t "abc" t] . [nil 7]) t]`,
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
