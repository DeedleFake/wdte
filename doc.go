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
// operator, "->", the ignored chain operator, "--", or the error
// chain operator, "-|". Each piece of the chain is executed in turn,
// and the output of the previous section is passed as an argument to
// the output of the current section. In other words, in the previous
// example, the chain's execution matches the following pseudocode
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
// The "-|" chain operator is used for error handling. During the
// evaluation of a chain, if no errors have occurred, chain segements
// using "-|" are ignored completely. Unlike with "--", they are
// completely not executed. If, however, an error occurs, all chain
// segments that don't use "-|" are ignored instead. If a "-|" segment
// exists in the chain after the location that the error occurred,
// then that segment is executed next, following which normal
// execution continues, unless that segment itself returned an error.
// If no "-|" segment exists, the error is returned from the entire
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
// semicolon separated series of cases in squiggly braces. A case is
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
// A few more minor points exist as well:
//    Array literals are a semicolon list of expression surrounded by
//    square brackets. Like in compounds and switches, the last
//    semicolon is optional.
//
//    Identifier parsing rules are very loose; essentially, anything
//    that isn't ambiguous with an existing keyword, operator, or
//    other syntactic construct is allowed.
//
//    All strings are essentially heredocs, allowing newlines like
//    they're any other character. There's no difference between
//    single-quoted and double-quoted strings.
//
//    There are no boolean literals, but the standard library provides
//    true and false functions that are essentially the same thing.
//
// Embedding
//
// As previously mentioned, everything in WDTE is a function. In Go
// terms, everything in WDTE implements the Func type defined in this
// package. This includes syntactic constructs as well, such as
// compounds, switches, and chains.
//
// When a script is parsed by one of the parsing functions in this
// package, it is translated into a recursive series of Func
// implementations. The specific types that it is translated to are
// all defined in and exported by this package. For example, the
// top-level of a script, being itself a compound, results in the
// instantiation of a Compound.
//
// What this means in terms of embedding is that the only thing
// required for interaction between Go and WDTE is an interoperative
// layer of Func implementations. As a functional language, WDTE is
// stateless; there is no global interpreter state to keep track of at
// all. Systems for tracking interpreter state, should they be
// required, are provided by the repl package.
//
// When a Func is called, it is passed a Frame. A Frame keeps track of
// anything the function needs that isn't directly an argument to the
// function. This includes the scope in which the Func call should be
// evaluated. For example, the expression
//
//    func arg1 arg2
//
// translates to an instance of the FuncCall implementation of Func.
// When the FuncCall is "called", it must be given a scope which
// contains, at a minimum, "func", "arg1", and "arg2", or the call
// will fail with an error. It is through this mechanism that new
// functions can be provided to WDTE. A custom scope can be created
// with new implementations of Func inserted into it. If this scope is
// inserted into a Frame which is then passed to a call of, for
// example, the top-level compound created by parsing a script, they
// will be available during the evaluation.
//
// Example:
//    const src = `
//      let io => import 'io';
//      io.writeln io.stdout example;
//    `
//
//    c, _ := wdte.Parse(strings.NewReader(src), std.Import)
//
//    scope := std.Scope.Add("example", wdte.String("This is an example."))
//    r := c.Call(std.F().WithScope(scope))
//    if err, ok := r.(error); ok {
//      log.Fatalln(err)
//    }
//
// This will print "This is an example." to stdout.
//
// For convenience, a simple function wrapper around the single method
// required by Func is provided in the form of GoFunc. GoFunc provides
// a number of extra features, such as automatically converting panics
// into errors, but for the most part is just a simple wrapper around
// manual implementations of Func. If more automatic behavior is
// required, possibly at the cost of some runtime performance,
// functions for automatically wrapping Go functions are provided in
// the wdteutil package.
package wdte
