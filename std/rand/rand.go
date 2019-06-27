// Package rand provides functions for generating and dealing with
// random numbers.
package rand

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/std/stream"
	"github.com/DeedleFake/wdte/wdteutil"
)

// A Source is a WDTE function that can create successive random
// numbers.
type Source interface {
	wdte.Func

	// Next returns the next random number from the generator.
	Next() wdte.Number
}

// Next is a WDTE function with the following signature:
//
//    next source
//
// It creates and returns the next random number from the given
// source.
func Next(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("next")

	if len(args) == 0 {
		return wdte.GoFunc(Next)
	}

	r := args[0].Call(frame).(Source)
	return r.Next()
}

type source struct {
	rand *rand.Rand
}

func (s *source) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return s
}

func (s *source) Next() wdte.Number {
	return wdte.Number(s.rand.Float64())
}

func (s *source) Reflect(name string) bool {
	return name == "Source"
}

func (s *source) String() string {
	return "<source>"
}

// Gen is a WDTE function with the following signature:
//
//    gen seed
//
// It returns a new Source that starts with the given seed.
func Gen(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("gen")

	if len(args) == 0 {
		return wdte.GoFunc(Gen)
	}

	seed := args[0].Call(frame).(wdte.Number)
	return &source{rand: rand.New(rand.NewSource(int64(seed)))}
}

type urand struct{}

func (urand) Int63() int64 {
	var r uint64
	err := binary.Read(crand.Reader, binary.LittleEndian, &r)
	if err != nil {
		panic(err)
	}
	return int64(r &^ (0x1 << 63))
}

func (urand) Seed(int64) {}

// UGen is a WDTE function with the following signature:
//
//    ugen
//
// It returns a Source that creates numbers from the operating
// system's cryptographic random number generator.
func UGen(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("ugen")
	return &source{rand: rand.New(urand{})}
}

// Stream is a WDTE function with the following signature:
//
//    stream source num
//
// It returns a Stream that yields the given number of random numbers
// from the provided Source.
func Stream(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("stream")

	if len(args) < 2 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Stream), args...)
	}

	r := args[0].Call(frame).(Source)
	rem := int(args[1].Call(frame).(wdte.Number))

	return stream.NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		if rem <= 0 {
			return nil, false
		}
		rem--

		return r.Next(), true
	})
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"gen":  wdte.GoFunc(Gen),
	"ugen": wdte.GoFunc(UGen),

	"next":   wdte.GoFunc(Next),
	"stream": wdte.GoFunc(Stream),
})

func init() {
	std.Register("rand", Scope)
}
