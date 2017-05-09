package ast_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte/ast"
)

func TestParse(t *testing.T) {
	ast, err := ast.Parse(strings.NewReader(`"test" => t; + x y => nil;`))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ast)
}
