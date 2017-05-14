package main

var examples = map[string]string{
	"fib": `# Welcome to the WDTE playground, a browser based evaluation
# environment for WDTE. This playground's features includes the
# standard function set as well as a number of importable modules.
#
# If you have never seen WDTE before and are completely confused at
# the moment, try reading the overview on the WDTE project's wiki:
# https://github.com/DeedleFake/wdte/wiki
#
# For documentation on the standard function set, see
# https://godoc.org/github.com/DeedleFake/wdte/std
#
# Importable modules:
# * 'math' (https://godoc.org/github.com/DeedleFake/wdte/std/math)
# * 'stream' (https://godoc.org/github.com/DeedleFake/wdte/std/stream)
# * 'canvas' (See the 'Canvas' example.)
#
# In addition, a print function is provided which uses the Go fmt
# package to create a string representation of its arguments. This
# string is printed to the output pane and then returned.

'math' => m;
'stream' => s;

memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

memo ! n => switch n {
	<= 1 => 1;
	default => - n 1 -> ! -> * n;
};

main => (
	fib 50 -> print;

	s.range (* m.pi -1) m.pi (/ m.pi 2)
	-> s.map m.sin
	-> s.collect
	-> print;
);`,

	"canvas": `# This example demonstrates the canvas module. This module is a module
# implemented just for this playground. Importing it automatically
# puts the playground in canvas mode, which allows for drawing to the
# output pane. It also redirects the normal output into the error pane.

'canvas' => c;

main =>
	c.start
	-> c.color 'purple'
	-> c.rect 10 10 100 50
	-> c.draw;`,
}
