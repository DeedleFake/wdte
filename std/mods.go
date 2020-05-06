package std

import (
	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

// A Memo wraps another function, caching the results of calls with
// the same arguments.
type Memo struct {
	Func wdte.Func
	Args []wdte.ID

	cache memoCache
}

func (m *Memo) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	s := frame.Scope()

	check := make([]wdte.Func, 0, len(m.Args))
	for _, id := range m.Args {
		check = append(check, s.Get(id).Call(frame))
	}

	cached, ok := m.cache.Get(check)
	if ok {
		return cached
	}

	r := m.Func.Call(frame, check...)
	m.cache.Set(check, r)
	return r
}

type memoCache struct {
	val  wdte.Func
	next map[wdte.Func]*memoCache
}

func (cache *memoCache) Get(args []wdte.Func) (wdte.Func, bool) {
	if cache == nil {
		return nil, false
	}

	if len(args) == 0 {
		return cache.val, true
	}

	if cache.next == nil {
		return nil, false
	}

	return cache.next[args[0]].Get(args[1:])
}

func (cache *memoCache) Set(args []wdte.Func, val wdte.Func) {
	if len(args) == 0 {
		cache.val = val
		return
	}

	if cache.next == nil {
		cache.next = make(map[wdte.Func]*memoCache)
	}

	n := new(memoCache)
	n.Set(args[1:], val)
	cache.next[args[0]] = n
}

func ModMemo(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(ModMemo)
	}

	frame = frame.Sub("memo")

	lambda := args[0].(*wdte.Lambda)

	argIDs := make([]wdte.ID, 0, len(lambda.Args))
	for _, arg := range lambda.Args {
		argIDs = append(argIDs, arg.IDs()...)
	}

	next := *lambda
	next.Expr = &Memo{
		Func: lambda.Expr,
		Args: argIDs,
	}
	return &next
}

func ModRev(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(ModRev)
	}

	frame = frame.Sub("rev")

	lambda := args[0].(*wdte.Lambda)

	var reverser wdte.Func
	reverser = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		if len(args) < len(lambda.Args) {
			return wdteutil.SaveArgsReverse(reverser, args...)
		}

		return lambda.Call(frame, args...)
	})
	return reverser
}
