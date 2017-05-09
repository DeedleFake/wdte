package wdte_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
)

func TestModule(t *testing.T) {
	const test = `
'test' => test;

add x y => + x y;
`

	m, err := wdte.Parse(strings.NewReader(test), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)
}
