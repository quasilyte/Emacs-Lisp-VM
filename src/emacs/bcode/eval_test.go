package bcode

import (
	"emacs/lisp"
	"errors"
	"fmt"
	"testing"
)

// testEnv is convenience Env wrapper that simplifies testing code.
type testEnv struct {
	Env

	symbols map[string]lisp.Object
}

// Symbol returns fsym for given name or lisp.Nil, if name is unbound.
func (env *testEnv) Symbol(name string) lisp.Object {
	fsym, ok := env.symbols[name]
	if ok {
		return fsym
	}
	return lisp.Nil
}

// AddFunc binds name to fn and returns associated Lisp symbol.
// Function is expected to be a valid Emacs Lisp compiled function.
func (env *testEnv) AddFunc(name string, fn Func) lisp.Object {
	if _, ok := env.symbols[name]; ok {
		panic(fmt.Sprintf("`%s` fsym is already bound", name))
	}

	fsym := lisp.NewSymbol(name)
	fsym.Symbol().FuncID = len(env.funcs)

	env.symbols[name] = fsym
	env.funcs = append(env.funcs, fn)

	return fsym
}

// AddFunc binds name to fn and returns associated Lisp symbol.
// Function is expected to be non-nil Go function.
func (env *testEnv) AddGoFunc(name string, fn GoFunc) lisp.Object {
	if _, ok := env.symbols[name]; ok {
		panic(fmt.Sprintf("`%s` fsym is already bound", name))
	}

	fsym := lisp.NewSymbol(name)
	fsym.Symbol().FuncID = len(env.goFuncs)

	env.symbols[name] = fsym
	env.goFuncs = append(env.goFuncs, fn)

	return fsym
}

func newTestEnv() *testEnv {
	master := MasterEnv{
		// Functions with ID=0 must be unassigned.
		goFuncs: make([]GoFunc, 1),
		funcs:   make([]Func, 1),
	}
	return &testEnv{
		Env: Env{
			MasterEnv: &master,
			stack:     make([]lisp.Object, 128),
			frames:    make([]callFrame, 32),
		},
		symbols: make(map[string]lisp.Object),
	}
}

// testInterpreter runs byte code evaluation that is
// divided into "steps".
//
// Each step has expected stack state that should
// match actual stack state after step code is executed.
type testInterpreter struct {
	*testEnv

	t *testing.T

	stepsCode  [][]byte
	stepsState []string
}

func newTestInterpreter(t *testing.T) *testInterpreter {
	return &testInterpreter{
		testEnv: newTestEnv(),
		t:       t,
	}
}

func (interp *testInterpreter) Run(name string, consts, args []lisp.Object) {
	stackDepth := uint32(len(args))
	env := &interp.testEnv.Env

	copy(env.stack, args)

	interp.t.Run(name, func(t *testing.T) {
		for i := range interp.stepsCode {
			fn := Func{
				code:   append(interp.stepsCode[i], OpExt, OpExtStop),
				consts: consts,
			}

			var err error
			stackDepth, err = eval(env, &fn, stackDepth)
			if err != ErrEOF {
				t.Fatalf("%s: step=%d/%d: eval error: %v",
					name, i, len(interp.stepsCode), err)
				return
			}

			want := interp.stepsState[i]
			have := lisp.ObjectSliceString(env.stack[:stackDepth])
			if have != want {
				t.Fatalf("%s: step=%d/%d: state mismatch:\nhave: <%s>\nwant: <%s>",
					name, i, len(interp.stepsCode), have, want)
				return
			}
		}
	})
}

func (interp *testInterpreter) LoadSteps(toks []interface{}) {
	code := interp.stepsCode[:0]
	state := interp.stepsState[:0]

	i := 0
	for i < len(toks) {
		switch op := toks[i+0].(byte); op {
		case OpExt:
			switch op := toks[i+1].(byte); op {
			case OpExtGoCallB:
				b := byte(toks[i+2].(int))
				code = append(code, []byte{OpExt, op, b})
				state = append(state, toks[i+3].(string))
				i += 4

			case OpExtGoCallW:
				b1 := byte(toks[i+2].(int) & 0x00FF)
				b2 := byte(toks[i+2].(int) & 0xFF00)
				code = append(code, []byte{OpExt, op, b1, b2})
				state = append(state, toks[i+3].(string))
				i += 4

			default:
				code = append(code, []byte{OpExt, op})
				state = append(state, toks[i+2].(string))
				i += 3
			}

		case OpStackRefB,
			OpVarRefB,
			OpVarSetB,
			OpVarBindB,
			OpCallB,
			OpUnbindB,
			OpRgotoB,
			OpRgotoIfNilB,
			OpRgotoIfNonNilB,
			OpRgotoIfNilElsePopB,
			OpRgotoIfNonNilElsePopB,
			OpListB,
			OpConcatB,
			OpStackSetB,
			OpDiscardB:
			b := byte(toks[i+1].(int))
			code = append(code, []byte{op, b})
			state = append(state, toks[i+2].(string))
			i += 3

		case OpStackRefW,
			OpVarRefW,
			OpVarSetW,
			OpCallW,
			OpUnbindW,
			OpGotoW,
			OpGotoIfNilW,
			OpGotoIfNonNilW,
			OpGotoIfNilElsePopW,
			OpGotoIfNonNilElsePopW,
			OpConstantW,
			OpStackSetW:
			b1 := byte(toks[i+1].(int) & 0x00FF)
			b2 := byte(toks[i+1].(int) & 0xFF00)
			code = append(code, []byte{op, b1, b2})
			state = append(state, toks[i+2].(string))
			i += 3

		default:
			code = append(code, []byte{op})
			state = append(state, toks[i+1].(string))
			i += 2
		}
	}

	interp.stepsCode = code
	interp.stepsState = state
	return
}

func TestPreconditions(t *testing.T) {
	if OpExtGoCall0 != 1 {
		t.Error("evalExt depends on OpExtGoCall0=1")
	}
}

func TestEval(t *testing.T) {
	interp := newTestInterpreter(t)

	// Byte code functions.
	push10 := interp.AddFunc("push10", Func{
		code: []byte{
			OpConstant0,
			OpReturn,
		},
		consts: []lisp.Object{lisp.NewInt(10)},
	})
	add2 := interp.AddFunc("add2", Func{
		code: []byte{
			OpAdd1,
			OpAdd1,
			OpReturn,
		},
	})

	// Go functions.
	goPushNil := interp.AddGoFunc("push-nil", func(args []lisp.Object) error {
		args[0] = lisp.Nil
		return nil
	})
	goAdd10 := interp.AddGoFunc("add10", func(args []lisp.Object) error {
		args[0] = lisp.NewInt(args[1].Int() + 10)
		return nil
	})
	goFloatToInt := interp.AddGoFunc("float-to-int", func(args []lisp.Object) error {
		if len(args) != 2 { // Additional arg is fsym
			return errors.New("float-to-int expects exactly one arg")
		}
		x := args[1]
		if x.Type != lisp.TypeFloat {
			return errors.New("float-to-int expects float arg")
		}
		args[0] = lisp.NewInt(int64(x.Float()))
		return nil
	})

	// These types are defined for readability.
	type (
		consts []interface{}
		args   []interface{}
		steps  []interface{}
	)
	tests := []struct {
		name   string
		consts []interface{}
		args   []interface{}
		steps  []interface{}
	}{
		{
			"Dup",
			consts{},
			args{1},
			steps{
				OpDup, `1 1`,
				OpDup, `1 1 1`,
				OpDup, `1 1 1 1`,
			},
		},

		{
			"Discard",
			consts{},
			args{1, 2.5, 3, 4, 5, 6},
			steps{
				OpDiscard, `1 2.5 3 4 5`,
				OpDiscard, `1 2.5 3 4`,
				OpDiscardB, 2, `1 2.5`,
				OpDiscardB, 1, `1`,
				OpDiscard, ``,
			},
		},

		{
			"StackRef",
			consts{},
			args{-1, 4, 3, 2, 1, 0},
			steps{
				OpStackRef3, `-1 4 3 2 1 0 3`,
				OpStackRef1, `-1 4 3 2 1 0 3 0`,
				OpStackRefB, 7, `-1 4 3 2 1 0 3 0 -1`,
				OpStackRefW, 1, `-1 4 3 2 1 0 3 0 -1 0`,
				OpDiscardB, 7, `-1 4 3`,
				OpStackRef2, `-1 4 3 -1`,
				OpStackRef3, `-1 4 3 -1 -1`,
				OpDup, `-1 4 3 -1 -1 -1`,
				OpStackRef5, `-1 4 3 -1 -1 -1 -1`,
			},
		},

		{
			"Constant",
			consts{0, 1, 2, 3, 4, 5},
			args{},
			steps{
				OpConstant0, `0`,
				OpConstantW, 0, `0 0`,
				OpConstant2, `0 0 2`,
				OpConstant2, `0 0 2 2`,
				OpConstantW, 1, `0 0 2 2 1`,
				OpConstant1, `0 0 2 2 1 1`,
				OpDiscardB, 6, ``,
				OpConstant5, `5`,
				OpConstant4, `5 4`,
			},
		},

		{
			"Call",
			consts{push10, add2},
			args{7},
			steps{
				OpConstant0, `7 push10`,
				OpCall0, `7 10`,
				OpConstant1, `7 10 add2`,
				OpStackRef1, `7 10 add2 10`,
				OpCall1, `7 10 12`,
			},
		},

		{
			"GoCall",
			consts{goAdd10, goFloatToInt, goPushNil},
			args{20, 7.7},
			steps{
				OpConstant0, `20 7.7 add10`,
				OpStackRef2, `20 7.7 add10 20`,
				OpExt, OpExtGoCall1, `20 7.7 30`,
				OpConstant1, `20 7.7 30 float-to-int`,
				OpStackRef2, `20 7.7 30 float-to-int 7.7`,
				OpExt, OpExtGoCall1, `20 7.7 30 7`,
				OpDiscardB, 4, ``,
				OpConstant2, `push-nil`,
				OpExt, OpExtGoCall0, `nil`,
				OpConstant2, `nil push-nil`,
				OpExt, OpExtGoCall0, `nil nil`,
			},
		},
	}

	for _, tt := range tests {
		consts := promoteObjects(tt.consts)
		args := promoteObjects(tt.args)

		interp.LoadSteps(tt.steps)
		interp.Run(tt.name, consts, args)
	}
}

// promoteObject replaces value of primitive type with valid lisp.Object.
func promoteObject(x interface{}) lisp.Object {
	switch x := x.(type) {
	case lisp.Object:
		return x
	case int:
		return lisp.NewInt(int64(x))
	case float64:
		return lisp.NewFloat(x)

	default:
		panic(fmt.Sprintf("unexpected value %#v", x))
	}
}

// promoteObjects calls promoteObject for each x in xs.
func promoteObjects(xs []interface{}) []lisp.Object {
	objects := make([]lisp.Object, len(xs))
	for i, x := range xs {
		objects[i] = promoteObject(x)
	}
	return objects
}
