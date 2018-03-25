export const fib = {
	name: 'Fibonacci',

	desc: `
Fibonacci
=========

This example provides a memoized implementation of a recursive Fibonacci number generator. It also provides a recursive factorial function for the heck of it.
`,

	input: `let memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	true => + (fib (- n 1)) (fib (- n 2));
};

let ! n => switch n {
	<= 1 => 1;
	true => - n 1 -> ! -> * n;
};

fib 30
-- print
-> / 5
-- print
;`,
}

export const stream = {
	name: 'Stream',

	desc: `
Stream
======

This example demonstrates the \`stream\` module. This module provides functional iterator operations, such as map, reduce, and filter.

For a full list of functions, see [the godocs][godoc].

[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/stream
`,

	input: `let m => import 'math';
let s => import 'stream';

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
;`,
}

export const strings = {
	name: 'Strings',

	desc: `
Strings
=======

This example demonstrates the \`strings\` module. This module provides basic string operations, such as finding the index of a substring, as well as more complicated operations, such as formatting.

For a full list of functions, including an explanation of the formatting system, see [the godocs][godoc].

[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/strings
`,

	input: `let s => import 'stream';
let str => import 'strings';

s.new 'abc' 'bcd' 'cde'
-> s.map (str.index 'cd')
-> s.collect
-- print
;

'This is the type of English up with which I will not put.'
-> str.format '{q}'
-- print
;`,
}

export const quine = {
	name: 'Quine',

	desc: `
Quine
=====

This example is an implemenation of a quine. That's about it.
`,

	input: `let str => import 'strings';
let q => "let str => import 'strings';\\nlet q => {q};\\nstr.format q q -- print;";
str.format q q -- print;`,
}

export const lambdas = {
	name: 'Lambdas',

	desc: `
Lambdas
=======

This example demonstrates lambdas by implementing an iterative Fibonacci number calculator using streams.
`,

	input: `let s => import 'stream';
let a => import 'arrays';

let fib n => s.range 1 n
	-> s.reduce [0; 1] (@ self p n => [
		at p 1;
		+ (at p 0) (at p 1);
	])
	-> at 1
	;

fib 30
-- print
;`,
}

//export const Canvas = {
//	desc: ``,
//	input: ``,
//}
