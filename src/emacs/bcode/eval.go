package bcode

import (
	"emacs/lisp"
)

// fetchB returns 8bit instruction argument at pc offset in code.
func fetchB(pc uint32, code []byte) uint32 {
	return uint32(code[pc+1])
}

// fetchW returns 16bit instruction argument at pc offset in code.
func fetchW(pc uint32, code []byte) uint32 {
	return uint32(code[pc+1]) +
		uint32(code[pc+2])<<8
}

// evalExt runs single instruction that is prefixed by OpExt byte.
//
// Moved outside of normal eval to preserve the code density
// of code that gets executed more frequently.
func evalExt(env *Env, fn *Func, sp, pc uint32) (uint32, error) {
	switch fn.code[pc] {
	case OpExtStop:
		return sp, ErrEOF

	case OpExtGoCall0, OpExtGoCall1, OpExtGoCall2, OpExtGoCall3, OpExtGoCall4, OpExtGoCall5:
		op := uint32(fn.code[pc])
		fsym := env.stack[sp-op].Symbol()
		err := env.goFuncs[fsym.FuncID](env.stack[sp-op : sp])
		if err != nil {
			return sp, err
		}
		return sp - op + 1, nil
	}

	return sp, nil
}

// eval is main byte code evaluating routine.
//
// Input arguments:
//   env - evaluation context.
//   fn - "main" function, evaluation entry point.
//   sp - stack pointer (position inside env.stack).
//
// Returns new stack pointer value along with error.
// Successful evaluation yields ErrEOF error value, not nil error.
//
// Does not catch Go panics.
func eval(env *Env, fn *Func, sp uint32) (uint32, error) {
	stack := env.stack
	frames := env.frames
	funcs := env.funcs

	// Zero frame always forces OpReturn to set pc to
	// trailing {OpExt,OpExtStop} that valid fn code should have.
	callDepth := 0
	frames[0].pc = uint32(len(fn.code) - 2)
	frames[0].fp = 0
	frames[0].fn = fn

	if safetyCheck {
		// Check that byte code really has trailing {OpExt,OpExtStop}.
		pc := frames[0].pc
		if fn.code[pc+0] != OpExt || fn.code[pc+1] != OpExtStop {
			return sp, ErrStopByte
		}
	}

	pc := uint32(0)

	for {
		switch fn.code[pc] {
		default:
			return sp, ErrBadOpcode

		case OpExt:
			var err error
			sp, err = evalExt(env, fn, sp, pc+1)
			if err != nil {
				return sp, err
			}
			pc += extOpWidth[fn.code[pc+1]]

		case OpStackRef1:
			stack[sp] = stack[sp-2]
			sp++
			pc++
		case OpStackRef2:
			stack[sp] = stack[sp-3]
			sp++
			pc++
		case OpStackRef3:
			stack[sp] = stack[sp-4]
			sp++
			pc++
		case OpStackRef4:
			stack[sp] = stack[sp-5]
			sp++
			pc++
		case OpStackRef5:
			stack[sp] = stack[sp-6]
			sp++
			pc++
		case OpStackRefB:
			n := fetchB(pc, fn.code)
			stack[sp] = stack[sp-n-1]
			sp++
			pc += 2
		case OpStackRefW:
			n := fetchW(pc, fn.code)
			stack[sp] = stack[sp-n-1]
			sp++
			pc += 3

		case OpCall0:
			callDepth++
			frames[callDepth].pc = pc
			frames[callDepth].fp = sp - 0
			frames[callDepth].fn = fn
			fn = &funcs[stack[sp-1].Symbol().FuncID]
			pc = 0
		case OpCall1:
			callDepth++
			frames[callDepth].pc = pc
			frames[callDepth].fp = sp - 1
			frames[callDepth].fn = fn
			fn = &funcs[stack[sp-2].Symbol().FuncID]
			pc = 0

		case OpCons:
			sp--
			car := stack[sp]
			cdr := stack[sp-1]
			stack[sp-1] = lisp.NewCons(car, cdr)
			pc++

		case OpDiscard:
			sp--
			pc++

		case OpReturn:
			frame := &frames[callDepth]
			stack[frame.fp-1] = stack[sp-1]
			sp = frame.fp
			fn = frame.fn
			pc = frame.pc + 1
			callDepth--

		case OpAdd1:
			// Issue#11
			x := &stack[sp-1]
			switch x.Type {
			case lisp.TypeInt:
				x.SetInt(x.Int() + 1)
			case lisp.TypeFloat:
				x.SetFloat(x.Float() + 1)
			}
			pc++

		case OpConstantW:
			stack[sp] = fn.consts[fetchW(pc, fn.code)]
			sp++
			pc += 3

		case OpGotoW:
			pc = fetchW(pc, fn.code)

		case OpGotoIfNilW:
			sp--
			top := &stack[sp]
			if lisp.Null(top) {
				pc = fetchW(pc, fn.code)
			} else {
				pc += 3
			}

		case OpDup:
			stack[sp] = stack[sp-1]
			sp++
			pc++

		case OpDiscardB:
			n := fetchB(pc, fn.code)
			sp -= n
			pc += 2

		case OpConstant0:
			stack[sp] = fn.consts[0]
			sp++
			pc++
		case OpConstant1:
			stack[sp] = fn.consts[1]
			sp++
			pc++
		case OpConstant2:
			stack[sp] = fn.consts[2]
			sp++
			pc++
		case OpConstant3:
			stack[sp] = fn.consts[3]
			sp++
			pc++
		case OpConstant4:
			stack[sp] = fn.consts[4]
			sp++
			pc++
		case OpConstant5:
			stack[sp] = fn.consts[5]
			sp++
			pc++
		case OpConstant6:
			stack[sp] = fn.consts[6]
			sp++
			pc++
		case OpConstant7:
			stack[sp] = fn.consts[7]
			sp++
			pc++
		case OpConstant8:
			stack[sp] = fn.consts[8]
			sp++
			pc++
		case OpConstant9:
			stack[sp] = fn.consts[9]
			sp++
			pc++
		case OpConstant10:
			stack[sp] = fn.consts[10]
			sp++
			pc++
		case OpConstant11:
			stack[sp] = fn.consts[11]
			sp++
			pc++
		case OpConstant12:
			stack[sp] = fn.consts[12]
			sp++
			pc++
		case OpConstant13:
			stack[sp] = fn.consts[13]
			sp++
			pc++
		case OpConstant14:
			stack[sp] = fn.consts[14]
			sp++
			pc++
		case OpConstant15:
			stack[sp] = fn.consts[15]
			sp++
			pc++
		case OpConstant16:
			stack[sp] = fn.consts[16]
			sp++
			pc++
		case OpConstant17:
			stack[sp] = fn.consts[17]
			sp++
			pc++
		case OpConstant18:
			stack[sp] = fn.consts[18]
			sp++
			pc++
		case OpConstant19:
			stack[sp] = fn.consts[19]
			sp++
			pc++
		case OpConstant20:
			stack[sp] = fn.consts[20]
			sp++
			pc++
		case OpConstant21:
			stack[sp] = fn.consts[21]
			sp++
			pc++
		case OpConstant22:
			stack[sp] = fn.consts[22]
			sp++
			pc++
		case OpConstant23:
			stack[sp] = fn.consts[23]
			sp++
			pc++
		case OpConstant24:
			stack[sp] = fn.consts[24]
			sp++
			pc++
		case OpConstant25:
			stack[sp] = fn.consts[25]
			sp++
			pc++
		case OpConstant26:
			stack[sp] = fn.consts[26]
			sp++
			pc++
		case OpConstant27:
			stack[sp] = fn.consts[27]
			sp++
			pc++
		case OpConstant28:
			stack[sp] = fn.consts[28]
			sp++
			pc++
		case OpConstant29:
			stack[sp] = fn.consts[29]
			sp++
			pc++
		case OpConstant30:
			stack[sp] = fn.consts[30]
			sp++
			pc++
		case OpConstant31:
			stack[sp] = fn.consts[31]
			sp++
			pc++
		case OpConstant32:
			stack[sp] = fn.consts[32]
			sp++
			pc++
		case OpConstant33:
			stack[sp] = fn.consts[33]
			sp++
			pc++
		case OpConstant34:
			stack[sp] = fn.consts[34]
			sp++
			pc++
		case OpConstant35:
			stack[sp] = fn.consts[35]
			sp++
			pc++
		case OpConstant36:
			stack[sp] = fn.consts[36]
			sp++
			pc++
		case OpConstant37:
			stack[sp] = fn.consts[37]
			sp++
			pc++
		case OpConstant38:
			stack[sp] = fn.consts[38]
			sp++
			pc++
		case OpConstant39:
			stack[sp] = fn.consts[39]
			sp++
			pc++
		case OpConstant40:
			stack[sp] = fn.consts[40]
			sp++
			pc++
		case OpConstant41:
			stack[sp] = fn.consts[41]
			sp++
			pc++
		case OpConstant42:
			stack[sp] = fn.consts[42]
			sp++
			pc++
		case OpConstant43:
			stack[sp] = fn.consts[43]
			sp++
			pc++
		case OpConstant44:
			stack[sp] = fn.consts[44]
			sp++
			pc++
		case OpConstant45:
			stack[sp] = fn.consts[45]
			sp++
			pc++
		case OpConstant46:
			stack[sp] = fn.consts[46]
			sp++
			pc++
		case OpConstant47:
			stack[sp] = fn.consts[47]
			sp++
			pc++
		case OpConstant48:
			stack[sp] = fn.consts[48]
			sp++
			pc++
		case OpConstant49:
			stack[sp] = fn.consts[49]
			sp++
			pc++
		case OpConstant50:
			stack[sp] = fn.consts[50]
			sp++
			pc++
		case OpConstant51:
			stack[sp] = fn.consts[51]
			sp++
			pc++
		case OpConstant52:
			stack[sp] = fn.consts[52]
			sp++
			pc++
		case OpConstant53:
			stack[sp] = fn.consts[53]
			sp++
			pc++
		case OpConstant54:
			stack[sp] = fn.consts[54]
			sp++
			pc++
		case OpConstant55:
			stack[sp] = fn.consts[55]
			sp++
			pc++
		case OpConstant56:
			stack[sp] = fn.consts[56]
			sp++
			pc++
		case OpConstant57:
			stack[sp] = fn.consts[57]
			sp++
			pc++
		case OpConstant58:
			stack[sp] = fn.consts[58]
			sp++
			pc++
		case OpConstant59:
			stack[sp] = fn.consts[59]
			sp++
			pc++
		case OpConstant60:
			stack[sp] = fn.consts[60]
			sp++
			pc++
		case OpConstant61:
			stack[sp] = fn.consts[61]
			sp++
			pc++
		case OpConstant62:
			stack[sp] = fn.consts[62]
			sp++
			pc++
		case OpConstant63:
			stack[sp] = fn.consts[63]
			sp++
			pc++
		}
	}
}
