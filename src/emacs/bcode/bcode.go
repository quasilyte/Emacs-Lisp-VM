package bcode

import "emacs/lisp"

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
