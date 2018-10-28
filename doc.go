// Package wdte implements the WDTE scripting language.
//
// WDTE is an embeddable, functionalish scripting language with a
// primary goal of simplicity of use from the embedding side, which is
// what this package provides.
//
// Quick Language Overview
//
// In order to understand how this package works, an overview of the
// language itself is first necessary. WDTE is functional-ish, with
// some emphasis on the "-ish". Although it generally follows a
// functional design, it is not purely functional. In WDTE, everything
// is a function in that everything can be "called", optionally with
// arguments. Value types return themselves, however, allowing them to
// be passed around.
//
// WDTE contains a construct called a "compound" which is similar to a
// function body in most languages. It is surrounded by parentheses
// and contains a semicolon separated list of expressions, with the
// last semicolon being optional. The top-level of a WDTE script is a
// compound without the parentheses. When a compound is executed, each
// expression is evaluated in turn. If any yield an error, that error
// is immediately returned from the entire compound. If not, the
// result of the last expression is returned.
//
// Example 1
//    # Declare a function called print3 that takes no arguments.
//    let print3 => (
//      let x => 3;
//      io.writeln io.stdout x;
//    );
//
// There are very few functions built-in in WDTE, but the standard
// library, found in the std directory and its subdirectories,
// contains a number of useful functions and definitions. For example,
// the stream module contains iterator functionality, which provides a
// means of looping over expressions, something which is otherwise not
// possible.
//
// Example 2
//    # Import 'stream' and 'array' and assign them to s and a,
//    # respectively. Note that an import is a compile-time operation,
//    # unlike normal functions. As such, it must be passed a string
//    # literal, not a variable.
//    let s => import 'stream';
//    let a => import 'arrays';
//
//    # Create a function called flatten that takes one argument,
//    # array.
//    let flatten array =>
//      # Create a new stream that iterates over array.
//      a.stream array
//
//      # Create a stream from the previous one that performs a flat
//      # map operation. The (@ name arg => ...) syntax is a lambda
//      # declaration.
//      -> s.flatMap (@ f v => v {
//          # If the current element of the stream, v, is an array,
//          # recursively flatten it into the stream.
//          reflect 'Array' => a.stream v -> s.flatMap f;
//        })
//
//      # Collect the previous stream into an array.
//      -> s.collect
//      ;
//
// This example also demonstrates "chains" and "switches", some
// features that seem complicated at first but quickly become second
// nature so with some practice.
//
// A chain is a series of expressions separated by either the chain
// operator, "->", or the ignored chain operator, "--". Each piece of
// the chain is executed in turn, and the output of the previous
// section is passed as an argument to the output of the current
// section. In other words, in the previous example, the chain's
// execution matches the following pseudocode
//    r1 = a.stream(array)
//    r2 = s.flatMap(<lambda>)
//    r1 = r2(r1)
//    r2 = s.collect
//    return r2(r1)
//
// A chain with a use of "--" operates in much the same way, but the
// output of the piece of the chain immediately following the operator
// is ignored, meaning that it doesn't affect the remainder of the
// chain.
//
// Chains can also have "slots" assigned to each piece. This is an
// identifier immediately following the expression of a piece of
// chain. This identifier is inserted into the scope for the remainder
// of the chain, allowing manual access to earlier sections of the
// chain. For example
//    let io => import 'io';
//    let file => import 'io/file';
//
//    let readFile path =>
//      file.open path : f
//      -> io.string
//      -- io.close f
//      ;
//
// The aforementioned switch expression is the only conditional
// provided by WDTE. It looks like an expression followed by a
// semicolon separated series of cases in squigly braces. A case is
// two expressions separated by the assignment operator, "=>". The
// original expression is first evaluated, following which each case's
// left-hand side is evaluated and the result of the original
// expression's evaluation is passed to it. If and only if this call
// results in the boolean value true, the right-hand side of that case
// is returned. If no cases match, the original expression is
// returned. For example,
//    func arg1 {
//      lhs1 arg2 => rhs1 arg3;
//      lhs2 arg4 => rhs2 arg6;
//    }
//
// This is analogous to the following pseudocode
//    check = func(arg1)
//    if lhs := lhs1(arg2); lhs(check) {
//      return rhs1(arg3)
//    }
//    if lhs := lhs2(arg4); lhs(check) {
//      return rhs2(arg6)
//    }
//    return check
//
// A few more minor points exist as well. Array literals are a
// semicolon list of expression surrounded by square brackets.
// Identifier parsing rules are very loose; essentially, anything that
// isn't ambiguous with an existing keyword, operator, or other
// syntactic construct is allowed.
package wdte
