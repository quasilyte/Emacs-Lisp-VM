package bcode

import (
	"emacs/lisp"
)

// MasterEnv holds data that is shared by multiple Env objects.
type MasterEnv struct {
	funcs   []Func
	goFuncs []GoFunc
}

// Env is a context that can be used to perform code evaluation.
// It is not thread/goroutine safe,
// each asynchronously running code should have it's own Env.
type Env struct {
	// stack holds local values that are pushed and popped
	// during evaluation.
	//
	// Shared across all activation records.
	//
	// len(stack) limits total data stack space available.
	stack []lisp.Object

	// frames saves previous activation record evaluation state.
	// When function returns, this information is used to continue
	// code execution.
	//
	// len(frames) limits call depth.
	frames []callFrame

	// MasterEnv holds information that is not required
	// to be bound to particular execution thread.
	*MasterEnv
}

// GoFunc is a type for Go functions that are callable
// from byte code via special opcode.
//
// The input slice argument contains at least one argument:
// function symbol that was used to make a call and
// function positional arguments.
//
// Return value must be placed in args[0].
// If args[0] is not changed, caller will have function
// symbol returned.
// Assign lisp.Nil explicitly if conventional void-like
// behavior is desired.
//
// Function can return a non-nil error which will
// trigger throw-like effect from Emacs Lisp point of view.
// Using Go panic directly may provide worse error message
// for the callee.
// Precise panic behavior inside byte code evaluation
// context provided elsewhere.
type GoFunc func(args []lisp.Object) error

// Func is compiled Emacs Lisp function object.
//
// Properties that are not related to evaluation are
// stored outside Func.
//
// Func object is safe to be executed inside multiple
// goroutines/interpreters.
type Func struct {
	code   []byte
	consts []lisp.Object
}

// callFrame holds single function call activation record data.
// Used during function return to restore interpreter state
// that can continue execution from the point right after the invocation.
type callFrame struct {
	// pc holds position inside fn code before function call
	// that spawned this frame.
	pc uint32

	// fp holds data stack index that is used to clear
	// stack upon function return.
	// It is also used to store function result properly.
	fp uint32

	// fn is a function that created this frame.
	// In other words, it is a caller function.
	fn *Func
}
