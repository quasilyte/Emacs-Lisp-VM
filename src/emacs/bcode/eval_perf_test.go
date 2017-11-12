package bcode

import (
	"emacs/lisp"
	"testing"
)

func BenchmarkEvalConstant(b *testing.B) {
	var code []byte
	{
		tmpl := []byte{
			OpConstant0,
			OpConstant1,
			OpConstant2,
			OpConstant3,
			OpConstant4,
			OpConstant5,
			OpDiscardB, 6,
		}
		for i := 0; i < 80; i++ {
			code = append(code, tmpl...)
		}
		code = append(code,
			OpExt,
			OpExtStop,
		)
	}

	env := newTestEnv()

	var consts []lisp.Object
	for i := int64(0); i < 10; i++ {
		consts = append(consts, lisp.NewInt(i))
	}
	main := Func{
		code:   code,
		consts: consts,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eval(&env.Env, &main, 0)
		if err != ErrEOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalStack(b *testing.B) {
	var code []byte
	{
		tmpl := []byte{
			OpDup, OpDup,
			OpStackRefB, 0,
			OpStackRefW, 0, 0,
			OpStackRef1,
			OpStackRef2,
			OpDiscard, OpDiscard,
			OpDiscardB, 2,
			OpDiscardB, 2,
		}
		for i := 0; i < 60; i++ {
			code = append(code, tmpl...)
		}
		code = append(code,
			OpExt,
			OpExtStop,
		)
	}

	env := newTestEnv()
	env.stack[0] = lisp.NewInt(0)

	main := Func{code: code}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eval(&env.Env, &main, 1)
		if err != ErrEOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalCall(b *testing.B) {
	var code []byte
	{
		tmpl := []byte{
			OpConstantW, 0, 0, OpCall0,
			OpConstant0, OpCall0,
		}
		for i := 0; i < 25; i++ {
			code = append(code, tmpl...)
		}
		code = append(code,
			OpExt,
			OpExtStop,
		)
	}

	env := newTestEnv()

	nopE := env.AddFunc("nopE", Func{
		code: []byte{
			OpReturn,
		},
	})
	nopD := env.AddFunc("nopD", Func{
		code: []byte{
			OpConstant0,
			OpCall0,
			OpReturn,
		},
		consts: []lisp.Object{nopE},
	})
	nopC := env.AddFunc("nopC", Func{
		code: []byte{
			OpConstant0,
			OpCall0,
			OpReturn,
		},
		consts: []lisp.Object{nopD},
	})
	nopB := env.AddFunc("nopB", Func{
		code: []byte{
			OpConstant0,
			OpCall0,
			OpReturn,
		},
		consts: []lisp.Object{nopC},
	})
	nopA := env.AddFunc("nopA", Func{
		code: []byte{
			OpConstant0,
			OpCall0,
			OpReturn,
		},
		consts: []lisp.Object{nopB},
	})

	main := Func{
		code:   code,
		consts: []lisp.Object{nopA},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eval(&env.Env, &main, 0)
		if err != ErrEOF {
			b.Fatal(err)
		}
	}
}
