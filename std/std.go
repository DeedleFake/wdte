package std

import (
	"fmt"
	"math"

	"github.com/DeedleFake/wdte"
)

func save(f wdte.Func, saved ...wdte.Func) wdte.Func {
	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		return f.Call(frame, append(args, saved...)...)
	})
}

// Plus is a WDTE function with the following signatures:
//
//    + a ...
//    (+ a) ...
//
// Returns the sum of a and the rest of its arguments.
func Plus(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Plus), args...)
	}

	frame = frame.Sub("+")

	var sum wdte.Number
	for _, arg := range args {
		arg = arg.Call(frame)
		if _, ok := arg.(error); ok {
			return arg
		}
		sum += arg.(wdte.Number)
	}
	return sum
}

// Minus is a WDTE with the following signatures:
//
//    - a b
//    (- b) a
//
// Returns a minus b.
func Minus(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Minus), args...)
	}

	frame = frame.Sub("-")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return a1.(wdte.Number) - a2.(wdte.Number)
}

// Times is a WDTE function with the following signatures:
//
//    * a ...
//    (* a) ...
//
// Returns the product of a and its other arguments.
func Times(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Times), args...)
	}

	frame = frame.Sub("*")

	p := wdte.Number(1)
	for _, arg := range args {
		arg = arg.Call(frame)
		if _, ok := arg.(error); ok {
			return arg
		}
		p *= arg.(wdte.Number)
	}
	return p
}

// Div is a WDTE function with the following signatures:
//
//    / a b
//    (/ b) a
//
// Returns a divided by b.
func Div(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Div), args...)
	}

	frame = frame.Sub("/")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return a1.(wdte.Number) / a2.(wdte.Number)
}

// Mod is a WDTE function with the following signatures:
//
//    % a b
//    (% b) a
//
// Returns a mod b.
func Mod(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Mod), args...)
	}

	frame = frame.Sub("%")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return wdte.Number(math.Mod(
		float64(a1.(wdte.Number)),
		float64(a2.(wdte.Number)),
	))
}

// Equals is a WDTE function with the following signatures:
//
//    == a b
//    (== b) a
//
// Returns true if a equals b. If a implements wdte.Comparer, the
// equality check is done using that implementation. If a does not but
// b does, b's implementation is used. If neither does, a direct Go
// equality check is used.
func Equals(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Equals), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, _ := cmp.Compare(a2)
		return wdte.Bool(c == 0)
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, _ := cmp.Compare(a1)
		return wdte.Bool(c == 0)
	}

	return wdte.Bool(a1 == a2)
}

// Less is a WDTE function with the following signatures:
//
//    < a b
//    (< b) a
//
// Returns true if a is less than b. Comparison rules are the same as
// those used for Equals, with the exception that the argument used
// must not only implement wdte.Comparer but that that implementation
// must support ordering.
func Less(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Less), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c < 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c > 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// Greater is a WDTE function with the following signatures:
//
//    > a b
//    (> b) a
//
// Returns true if a is greater than b. Comparison rules are the same
// as those used for Equals, with the exception that the argument used
// must not only implement wdte.Comparer but that that implementation
// must support ordering.
func Greater(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Greater), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c > 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c < 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// LessEqual is a WDTE function with the following signatures:
//
//    <= a b
//    (<= b) a
//
// Returns true if a is less than or equal to b. Comparison rules are
// the same as those used for Equals, with the exception that the
// argument used must not only implement wdte.Comparer but that that
// implementation must support ordering.
func LessEqual(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(LessEqual), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c <= 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c >= 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// GreaterEqual is a WDTE function with the following signatures:
//
//    >= a b
//    (>= b) a
//
// Returns true if a is greater than or equal to b. Comparison rules
// are the same as those used for Equals, with the exception that the
// argument used must not only implement wdte.Comparer but that that
// implementation must support ordering.
func GreaterEqual(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(GreaterEqual), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c >= 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c <= 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

const (
	// True is a WDTE function with the following signature:
	//
	//    true
	//
	// As you can probably guess, it returns a boolean true.
	True wdte.Bool = true

	// False is a WDTE function with the following signature:
	//
	//    false
	//
	// Returns a boolean false. This is rarely necessary as most
	// built-in functionality considers any value other than a boolean
	// true to be false, but it's provided for completeness.
	False wdte.Bool = false
)

// And is a WDTE function with the following signature:
//
//    && ...
//
// Returns true if all of its arguments are true.
func And(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("&&")

	switch len(args) {
	case 0:
		return wdte.GoFunc(And)
	}

	for _, arg := range args {
		arg = arg.Call(frame)
		if arg != wdte.Bool(true) {
			return wdte.Bool(false)
		}
	}

	return wdte.Bool(true)
}

// Or is a WDTE function with the following signature:
//
//    || ...
//
// Returns true if any of its arguments are true.
func Or(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("||")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Or)
	}

	for _, arg := range args {
		arg = arg.Call(frame)
		if arg == wdte.Bool(true) {
			return wdte.Bool(true)
		}
	}

	return wdte.Bool(false)
}

// Not is a WDTE function with the following signature:
//
//    ! a
//
// Returns true if a is not true or false if a is not true.
func Not(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("!")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Not)
	}

	return wdte.Bool(args[0].Call(frame) != wdte.Bool(true))
}

// Len is a WDTE function with the following signature:
//
//    len a
//
// Returns the length of a if a implements wdte.Lenner, or false if it
// doesn't.
func Len(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("len")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Len)
	}

	if lenner, ok := args[0].Call(frame).(wdte.Lenner); ok {
		return wdte.Number(lenner.Len())
	}

	return wdte.Bool(false)
}

// At is a WDTE function with the following signatures:
//
//    at a i
//    (at i) a
//
// Returns the ith index of a. a is assumed to implement wdte.Atter.
func At(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(At), args...)
	}

	frame = frame.Sub("at")

	at := args[0].Call(frame).(wdte.Atter)
	i := args[1].Call(frame)

	ret, ok := at.At(i)
	if !ok {
		return &wdte.Error{
			Frame: frame,
			Err:   fmt.Errorf("index %v out of range", i),
		}
	}

	return ret
}

// Collect is a WDTE function with the following signature:
//
//    collect compound
//
// Collect takes a compound as its argument and returns the scope
// collected from executing that compound. The argument must be a
// compound literal or the function will fail. Assigning a compound to
// an ID and then passing that ID will not work.
//
// It surrounds the returned scope with a bound called "collect".
func Collect(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Collect)
	}

	frame = frame.Sub("collect")

	sf := args[0].(*wdte.ScopedFunc)
	s, _ := sf.Func.(wdte.Compound).Collect(frame.WithScope(sf.Scope))
	return s
}

// Known is a WDTE function with the following signature:
//
//    known scope
//
// Returns an array containing known identifiers in the given scope
// sorted alphabetically.
func Known(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Known)
	}

	frame = frame.Sub("known")

	s := args[0].Call(frame).(*wdte.Scope)
	k := s.Known()

	ret := make(wdte.Array, 0, len(k))
	for _, id := range k {
		ret = append(ret, wdte.String(id))
	}

	return ret
}

// Sub is a WDTE function with the following signatures:
//
//    sub scope id val
//    (sub val) scope id
//    (sub id val) scope
//
// Sub returns a subscope of scope with the value val bound to the ID
// id. It puts a "collect" lower bound on the scope but does not add
// any upper bounds.
func Sub(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 2 {
		return save(wdte.GoFunc(Sub), args...)
	}

	frame = frame.Sub("sub")

	s := args[0].Call(frame).(*wdte.Scope)
	id := wdte.ID(args[1].Call(frame).(wdte.String))
	v := args[2]

	return s.Add(id, v)
}

// Scope is a scope containing the functions in this package.
//
// This scope is primarily useful for bootstrapping an environment for
// running scripts in. To use it, simply pass a frame containing it or
// a subscope of it to a function call. In many cases, a client can
// simply call F to obtain such a frame.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"+": wdte.GoFunc(Plus),
	"-": wdte.GoFunc(Minus),
	"*": wdte.GoFunc(Times),
	"/": wdte.GoFunc(Div),
	"%": wdte.GoFunc(Mod),

	"==":    wdte.GoFunc(Equals),
	"<":     wdte.GoFunc(Less),
	">":     wdte.GoFunc(Greater),
	"<=":    wdte.GoFunc(LessEqual),
	">=":    wdte.GoFunc(GreaterEqual),
	"true":  True,
	"false": False,
	"&&":    wdte.GoFunc(And),
	"||":    wdte.GoFunc(Or),
	"!":     wdte.GoFunc(Not),

	"len":     wdte.GoFunc(Len),
	"at":      wdte.GoFunc(At),
	"collect": wdte.GoFunc(Collect),
	"known":   wdte.GoFunc(Known),
	"sub":     wdte.GoFunc(Sub),
})

// F returns a top-level frame that has S as its scope.
func F() wdte.Frame {
	return wdte.F().WithScope(Scope)
}
