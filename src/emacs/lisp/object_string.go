package lisp

import (
	"fmt"
	"strconv"
	"strings"
)

// ObjectString returns stringified representation of o.
// Output is not guaranteed to be prin1-compatible.
func ObjectString(o Object) string {
	switch o.Type {
	case TypeInt:
		return strconv.FormatInt(o.Int(), 10)

	case TypeFloat:
		s := strconv.FormatFloat(o.Float(), 'f', -1, 64)
		// Always provide fractional part.
		if strings.IndexByte(s, '.') == -1 {
			return s + ".0"
		}
		return s

	case TypeSymbol:
		if name := o.Symbol().Name; name != "" {
			return name
		}
		return "##" // Emacs Lisp notation for empty symbol name

	case TypeVector:
		return "[" + ObjectSliceString(o.Vector().Vals) + "]"

	case TypeCons:
		cons := o.Cons()
		car := ObjectString(cons.Car)
		cdr := ObjectString(cons.Cdr)
		return fmt.Sprintf("(%s . %s)", car, cdr)

	default:
		return fmt.Sprint(o)
	}
}

// ObjectSliceString maps all objects through ObjectString
// and returns space-separated result.
func ObjectSliceString(objects []Object) string {
	parts := make([]string, len(objects))
	for i, o := range objects {
		parts[i] = ObjectString(o)
	}
	return strings.Join(parts, " ")
}
