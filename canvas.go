package main

import "github.com/DeedleFake/wdte"

type Shape interface {
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
		args[0].Call(frame).(Shape).Draw()
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

	return &rect{
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
		args[0].Call(frame).(Shape).Draw()
		canvasCtx.Call("fillRect", r.x, r.y, r.w, r.h)
	})
}

func Draw(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) > 0 {
		args[0].Call(frame).(Shape).Draw()
	}

	return wdte.GoFunc(Draw)
}

func CanvasModule() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"start": wdte.GoFunc(Start),
			"color": wdte.GoFunc(Color),
			"rect":  wdte.GoFunc(Rect),
			"draw":  wdte.GoFunc(Draw),
		},
	}
}
