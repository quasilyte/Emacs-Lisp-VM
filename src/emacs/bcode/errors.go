package bcode

import (
	"errors"
)

// Simple errors that do not provide much context information,
// but can be compared directly.
//
// Should be treated as constants.
var (
	// ErrEOF is "not an error" condition that signals that whole
	// byte code input has been consumed (evaluated).
	ErrEOF = errors.New("byte code EOF")

	// ErrStopByte reports about malformed byte code sequences
	// that miss trailing []byte{OpExt, OpExtStop} bytes.
	ErrStopByte = errors.New("code misses trailing {OpExt,OpExtStop} bytes")

	// ErrBadOpcode reports invalid/unsupported opcode that was
	// about to be evaluated.
	//
	// This should be a rich type of error with PC offset
	// and opcode values stored inside.
	// It may as well require call trace attached,
	// that is currently unimplemented (Issue#12).
	ErrBadOpcode = errors.New("found unexpected opcode")
)
