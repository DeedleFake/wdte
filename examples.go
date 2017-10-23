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
# The standard library is available for importing, along with the
# following special modules:
#
# * 'canvas' (See the 'Canvas' example.)
# * 'io/file' (This makes no sense in a browser, so it's disabled.)
# * 'io' (Disabled pending https://github.com/gopherjs/gopherjs/issues/705.)
#
# In addition, a print function is provided which uses the Go fmt
# package to create a string representation of its arguments. This
# string is printed to the output pane and then returned.

memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

memo ! n => switch n {
	<= 1 => 1;
	default => - n 1 -> ! -> * n;
};

main =>
	fib 30
	-- print
	-> / 5
	-- print
	;`,

	"stream": `# This example demonstrates the stream module. This module allows
# classic functional iterator operations, such as map, reduce, and
# filter.

'math' => m;
'stream' => s;

main => (
	print 'Map and filter:';
	s.range 0 (* m.pi 2) (/ m.pi 2)
	-> s.map m.sin
	-> s.filter (>= 0)
	-> s.map print
	-> s.drain
	;

	print 'Reduce:';
	s.range 1 5
	-> s.reduce 1 *
	-- print
	;
);`,

	"strings": `# This example demonstrates the strings module. This module provides
# basic string operations, such as finding the index of a substring,
# as well as more complicated operations, such as formatting.

'stream' => s;
'strings' => str;

main => (
	s.new 'abc' 'bcd' 'cde'
	-> s.map (str.index 'cd')
	-> s.collect
	-- print
	;

	'This is the type of English up with which I will not put.'
	-> str.format '{q}'
	-- print
	;
);`,

	"lambda": `# This example demonstrates lambdas by implementing an iterative
# Fibonacci number calculator using streams.

'stream' => s;
'arrays' => a;

fib n => s.range 1 n
	-> s.reduce [0; 1] (@ self p n => [
		a.at p 1;
		+ (a.at p 0) (a.at p 1);
	])
	-> a.at 1
	;

main =>
	fib 30
	-- print
	;`,

	"canvas": `# This example demonstrates the canvas module. This module is a module
# implemented just for this playground. Importing it automatically
# puts the playground in canvas mode, which allows for drawing to the
# output pane. It also redirects the normal output into the error pane.

'math' => m;
'canvas' => c;

# circleRec is the recursive part of circle.
circleRec r cx cy res t p => switch t {
	> (* 2 m.pi) => p -> c.close;
	default => p
             -> c.line (+ cx (* (m.cos t) r)) (+ cy (* (m.sin t) r))
             -> circleRec r cx cy res (+ t (/ m.pi res));
};

# circle returns a circular path centered at (cx, cy) with radius r
# and resolution res.
circle r cx cy res => c.path
                  -> c.move (+ cx r) cy
                  -> circleRec r cx cy res (/ m.pi res);

main => (
	c.start
	-> c.color 'purple'
	-> c.rect 10 10 100 50
	-> c.draw;

	c.start
	-> c.color 'pink'
	-> (
		c.path
		-> c.move 10 50
		-> c.line 30 30
		-> c.line 50 30
		-> c.line 100 100
		-> c.close
	)
	-> c.draw;

	c.start
	-> c.color 'red'
	-> circle 30 100 100 6
	-> c.draw;
);`,
}
