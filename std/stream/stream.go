// Package stream provides WDTE functions for manipulating streams of
// data.
package stream

import "github.com/DeedleFake/wdte"

type Stream interface {
	wdte.Func

	Next(frame []wdte.Func) (wdte.Func, bool)
}

type array struct {
	a wdte.Array
	i int
}

func New(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(New)

	case 1:
		a1 := args[0].Call(frame)
		if a1, ok := a1.(wdte.Array); ok {
			return &array{a: a1}
		}
	}

	return &array{a: args}
}

func (a *array) Call(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	return a
}

func (a *array) Equals(other wdte.Func) bool {
	panic("Not implemented.")
}

func (a *array) Next(frame []wdte.Func) (wdte.Func, bool) {
	if a.i >= len(a.a) {
		return nil, false
	}

	r := a.a[a.i]
	a.i++
	return r, true
}

type rng struct {
	i wdte.Number
	m wdte.Number
	s wdte.Number
}

func Range(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Range)

	case 1:
		return &rng{
			m: args[0].Call(frame).(wdte.Number),
			s: 1,
		}

	case 2:
		return &rng{
			i: args[0].Call(frame).(wdte.Number),
			m: args[1].Call(frame).(wdte.Number),
			s: 1,
		}
	}

	return &rng{
		i: args[0].Call(frame).(wdte.Number),
		m: args[1].Call(frame).(wdte.Number),
		s: args[2].Call(frame).(wdte.Number),
	}
}

func (r *rng) Call(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	return r
}

func (r *rng) Equals(other wdte.Func) bool {
	panic("Not implemented.")
}

func (r *rng) Next(frame []wdte.Func) (wdte.Func, bool) {
	if r.i >= r.m {
		return nil, false
	}

	n := r.i
	r.i += r.s
	return n, true
}

type mapper struct {
	m wdte.Func
	s Stream
}

func Map(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	return &mapper{m: args[0].Call(frame)}
}

func (m *mapper) Call(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			m.s = a
		}
	}

	return m
}

func (m *mapper) Equals(other wdte.Func) bool {
	panic("Not implemented.")
}

func (m *mapper) Next(frame []wdte.Func) (wdte.Func, bool) {
	if m.s == nil {
		return nil, false
	}

	n, ok := m.s.Next(frame)
	if !ok {
		return nil, false
	}

	return m.m.Call(frame, n), true
}

func Collect(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Collect)
	}

	a, ok := args[0].Call(frame).(Stream)
	if !ok {
		return args[0]
	}

	r := wdte.Array{}
	for {
		n, ok := a.Next(frame)
		if !ok {
			break
		}

		r = append(r, n.Call(frame))
	}

	return r
}

func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"new":     wdte.GoFunc(New),
			"range":   wdte.GoFunc(Range),
			"map":     wdte.GoFunc(Map),
			"collect": wdte.GoFunc(Collect),
		},
	}
}
