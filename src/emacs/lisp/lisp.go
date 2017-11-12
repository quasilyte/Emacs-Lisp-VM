package lisp

import (
	"unsafe"
)

// Type is an Emacs Lisp type unique ID.
//
// Note that there is no way to define real
// user types in Emacs Lisp, so the set of these
// values is closed.
type Type uint64

// All Emacs Lisp type tags that are supported by this runtime.
// See Object type doc-comment for more information.
const (
	TypeInt Type = iota
	TypeFloat
	TypeSymbol
	TypeVector
	TypeCons
	TypeString
)

// Object is universal Emacs Lisp value.
// The type is bound dynamically.
// The sync of type and value is required for
// Object to function properly.
//
// This type emulates C-style union.
//
// Possible values:
//   {Type: TypeInt, Num: int64}
//   {Type: TypeFloat, Num: float64}
//   {Type: TypeSymbol, Ptr: *Symbol}
//   {Type: TypeVector, Ptr: *Vector}
//   {Type: TypeCons, Ptr: *Cons}
//   {Type: TypeString: Ptr: *String}
type Object struct {
	// Warning: Num member should always be the first,
	// because it is accessed via unsafe pointer at zero offset.

	// Non-heap 64bit value. Used for ints and floats.
	Num uintptr

	// A pointer to reference-type value.
	Ptr unsafe.Pointer

	// Type descriptor (or tag) for stored value.
	// In other words, it is Object dynamic type.
	Type Type
}

// Int returns object integer value.
// UB if o.Type is not TypeInt.
func (o *Object) Int() int64 {
	return *(*int64)(unsafe.Pointer(o))
}

// Float returns object float value.
// UB if o.Type is not TypeFloat.
func (o *Object) Float() float64 {
	return *(*float64)(unsafe.Pointer(o))
}

// Symbol returns object value as a symbol object.
// UB if o.Type is not TypeSymbol.
func (o *Object) Symbol() *Symbol {
	return (*Symbol)(o.Ptr)
}

// Vector returns object vector value.
// UB if o.Type is not TypeVector.
func (o *Object) Vector() *Vector {
	return (*Vector)(o.Ptr)
}

// Cons returns object cons value.
// UB if o.Type is not TypeCons.
func (o *Object) Cons() *Cons {
	return (*Cons)(o.Ptr)
}

// String returns object string value.
// UB if o.Type is not TypeString.
func (o *Object) String() *String {
	return (*String)(o.Ptr)
}

// SetInt updates object integer value.
// UB if o.Type is not TypeInt.
func (o *Object) SetInt(val int64) {
	*(*int64)(unsafe.Pointer(o)) = val
}

// SetFloat updates object float value.
// UB if o.Type is not TypeFloat.
func (o *Object) SetFloat(val float64) {
	*(*float64)(unsafe.Pointer(o)) = val
}

// Symbol is env-local interned string.
// A symbol name is unique, no two symbols have same name.
type Symbol struct {
	Name   string
	FuncID int
}

// Vector is a fixed-size dynamic array.
type Vector struct {
	Vals []Object
}

// Cons is a Lisp-y pair.
type Cons struct {
	Car Object // First (head for lists).
	Cdr Object // Second (tail for lists).
}

// String is like Vector, but stores chars instead of
// arbitrary Lisp objects.
type String struct {
	Chars []byte
}

// Values that are defined by default and considered immutable.
var (
	// Nil is the only false value in Emacs Lisp.
	// Basically, nil is a symbol and it's value is
	// empty list.
	Nil = NewSymbol("nil")

	// T is preffered Emacs Lisp truth value for predicates.
	// It is more or less the same as boolean "true", but
	// still has symbol type.
	T = NewSymbol("t")
)

// NewInt constructs Object initialized with integer val.
func NewInt(val int64) Object {
	o := Object{Type: TypeInt}
	o.SetInt(val)
	return o
}

// NewFloat constructs Object initialized with float val.
func NewFloat(val float64) Object {
	o := Object{Type: TypeFloat}
	o.SetFloat(val)
	return o
}

// NewSymbol returns a newly allocated uninterned symbol for given name.
// The symbol value is void.
func NewSymbol(name string) Object {
	return Object{
		Type: TypeSymbol,
		Ptr:  unsafe.Pointer(&Symbol{Name: name}),
	}
}

// NewVector returns a vector Object initialized with vals.
func NewVector(vals []Object) Object {
	return Object{
		Type: TypeVector,
		Ptr:  unsafe.Pointer(&Vector{Vals: vals}),
	}
}

// NewCons returns a cons Object initialized with {car, cdr}.
func NewCons(car, cdr Object) Object {
	return Object{
		Type: TypeCons,
		Ptr:  unsafe.Pointer(&Cons{Car: car, Cdr: cdr}),
	}
}

// NewString returns a string Object initialized with chars.
func NewString(chars []byte) Object {
	return Object{
		Type: TypeString,
		Ptr:  unsafe.Pointer(&String{Chars: chars}),
	}
}

// Bool maps Go boolean value to Emacs Lisp closest equivalents.
//
// true => t symbol
// false => nil symbol
//
// The returned object value should be treated as readonly.
func Bool(x bool) Object {
	if x {
		return T
	}
	return Nil
}

// Null only returns true for Nil.
func Null(x *Object) bool {
	return x.Type == TypeSymbol &&
		x.Ptr == Nil.Ptr
}

// Eq returns true if x and y are same Lisp objects.
//
// For integers and floats, it does proper structural comparison,
// for other types it only performs referential comparison.
//
// Issue#1.
func Eq(x, y *Object) bool {
	return *x == *y
}
