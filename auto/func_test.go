package auto_test

import (
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/auto"
)

func TestFunc(t *testing.T) {
	test := func(v int) int {
		return v + 3
	}

	f := auto.Func("test", test)
	t.Log(f.Call(wdte.F(), wdte.Number(2)))
}
