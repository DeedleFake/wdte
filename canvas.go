package main

import "github.com/DeedleFake/wdte"

type Drawer interface {
	wdte.Func

	Draw()
}

type drawFunc func()

func (d drawFunc) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return d
}

func (d drawFunc) Draw() {
	d()
}

func Start(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return drawFunc(func() {
		canvasCtx.Set("fillStyle", "black")
	})
}

type color wdte.String

func Color(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return color(args[0].Call(frame).(wdte.String))
}

func (c color) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return c
	}

	return drawFunc(func() {
		args[0].Call(frame).(Drawer).Draw()
		canvasCtx.Set("fillStyle", c)
	})
}

type rect struct {
	x, y wdte.Number
	w, h wdte.Number
}

func Rect(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Rect)
	}

	return rect{
		x: args[0].Call(frame).(wdte.Number),
		y: args[1].Call(frame).(wdte.Number),
		w: args[2].Call(frame).(wdte.Number),
		h: args[3].Call(frame).(wdte.Number),
	}
}

func (r rect) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return r
	}

	return drawFunc(func() {
		args[0].Call(frame).(Drawer).Draw()
		canvasCtx.Call("fillRect", r.x, r.y, r.w, r.h)
	})
}

func Draw(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) > 0 {
		args[0].Call(frame).(Drawer).Draw()
	}

	return wdte.GoFunc(Draw)
}

type Pather interface {
	wdte.Func

	Path()
}

type pathFunc func()

func (p pathFunc) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return p
}

func (p pathFunc) Path() {
	p()
}

func Path(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return pathFunc(func() {
		canvasCtx.Call("beginPath")
	})
}

type move struct {
	x, y wdte.Number
}

func Move(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Move)
	}

	return move{
		x: args[0].Call(frame).(wdte.Number),
		y: args[1].Call(frame).(wdte.Number),
	}
}

func (m move) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return m
	}

	return pathFunc(func() {
		args[0].Call(frame).(Pather).Path()
		canvasCtx.Call("moveTo", m.x, m.y)
	})
}

type line struct {
	x, y wdte.Number
}

func Line(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Line)
	}

	return line{
		x: args[0].Call(frame).(wdte.Number),
		y: args[1].Call(frame).(wdte.Number),
	}
}

func (l line) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return l
	}

	return pathFunc(func() {
		args[0].Call(frame).(Pather).Path()
		canvasCtx.Call("lineTo", l.x, l.y)
	})
}

type pathShape struct {
	path Pather
}

func Close(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Close)
	}

	return pathShape{
		path: args[0].Call(frame).(Pather),
	}
}

func (p pathShape) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return p
	}

	return drawFunc(func() {
		args[0].Call(frame).(Drawer).Draw()

		p.path.Path()
		canvasCtx.Call("fill")
	})
}

func CanvasModule() *wdte.Scope {
	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"start": wdte.GoFunc(Start),
		"color": wdte.GoFunc(Color),
		"rect":  wdte.GoFunc(Rect),
		"draw":  wdte.GoFunc(Draw),

		"path":  wdte.GoFunc(Path),
		"move":  wdte.GoFunc(Move),
		"line":  wdte.GoFunc(Line),
		"close": wdte.GoFunc(Close),
	})
}
