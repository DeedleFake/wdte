// @format

export const fib = {
	name: 'Fibonacci',

	desc: `
Fibonacci
=========

This example provides a memoized implementation of a recursive Fibonacci number generator. It also provides a recursive factorial function for the heck of it.
`,

	input: `let memo fib n => n {>= 2 => + (fib (- n 1)) (fib (- n 2))};

let ! n => n {
	<= 1 => 1;
	true => - n 1 -> ! -> * n;
};

fib 30
-- io.writeln io.stdout
-> / 5
-- io.writeln io.stdout
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

io.writeln io.stdout 'Map and filter:';
s.range 0 (* m.pi 2) (/ m.pi 2)
-> s.map m.sin
-> s.filter (>= 0)
-> s.map (io.writeln io.stdout)
-> s.drain
;

io.writeln io.stdout 'Reduce:';
s.range 1 5
-> s.reduce 1 *
-- io.writeln io.stdout
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
-- io.writeln io.stdout
;

'This is the type of English up with which I will not put.'
-> str.format '{q}'
-- io.writeln io.stdout
;`,
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
-- io.writeln io.stdout
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
let q => "let str => import 'strings';\\nlet q => {q};\\nstr.format q q -- io.writeln io.stdout;";
str.format q q -- io.writeln io.stdout;`,
}

export const hundredDoors = {
	name: '100 Doors',

	desc: `
100 Doors
=========

The [100 doors problem](https://www.rosettacode.org/wiki/100_doors), as presented by Rosetta Code, is as follows:

There are 100 doors that are all closed. You walk past the doors in the same direction 100 times. On the first pass, you toggle the state of every door, opening closed doors and closing open doors. On the second pass, you toggle every second door. On the third you toggle every third door. Etc.

This example simulates this scenario, printing out the final state of the doors.
`,

	input: `let a => import 'arrays';
let s => import 'stream';

let toggle doors m =>
	a.stream doors
	-> s.enumerate
	-> s.map (@ s n => [+ (at n 0) 1; at n 1])
	-> s.map (@ s n => n {
			(@ s n => == (% (at n 0) m) 0) => ! (at n 1);
			true => at n 1;
		})
	-> s.collect
	;

s.range 100
-> s.map false
-> s.collect : doors
-> s.range 1 100
-> s.reduce doors toggle
-> a.stream
-> s.map (@ s n => 0 {
		n => 'Open';
		true => 'Closed';
	} -- io.writeln io.stdout)
-> s.drain
;`,
}

//export const Canvas = {
//	desc: ``,
//	input: ``,
//}
