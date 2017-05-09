package ast_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte/ast"
)

func printTree(t *testing.T, cur ast.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	switch cur := cur.(type) {
	case *ast.Term:
		t.Logf("%v%v", indent, cur)

	case *ast.NTerm:
		t.Logf("%v(%v", indent, cur)
		for _, c := range cur.Children() {
			printTree(t, c, depth+1)
		}
		t.Logf("%v)", indent)

	case *ast.Epsilon:
		t.Logf("%vε", indent)

	default:
		t.Fatalf("Unexpected node: %#v", cur)
	}
}

func TestParse(t *testing.T) {
	//const test = `"test" => t; + x y => nil;`

	const test = `
'test' => test;

fib n => switch n {
	0 => 0;
	default => + (fib (- n 1)) (fib (- n 2));
};

main => print (fib 5);
`

	root, err := ast.Parse(strings.NewReader(test))
	if err != nil {
		t.Fatal(err)
	}
	printTree(t, root, 0)
}
