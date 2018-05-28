WDTE Language Specification
===========================

*HEAVILY WORK IN PROGRESS* The structure of this is heavily based on [the Go spec](https://www.golang.org/doc/spec).

TODO: Put a table of contents here.

Introduction
------------

This document specifies the WDTE syntax and calling structure. It deals purely with defining expected behavior, not with implementation details.

WDTE is an small, functional-ish, embedded scripting language designed for integration into other projects. It is dynamically typed. All values in WDTE are functions, capable of being called and passed arguments, and there are no statements, only expressions. There are situations in which the result of an expression evaluation can be ignored, however, and this is why the language is only functional*-ish*.

Notation
--------

The syntax is specified using a variant of Extended Backus-Naur Form. Names of production rules are surrounded by `<>`, and the left and right sides of the rules are separated via `->`. Alternatives are specified using `|`, and epsilon is specified via `ε`. Literals are anything that isn't one of the already listed items, although `->` and `|` are legal values if they are not used as specified previously. The following literals are special:

* `number`: Any floating point value, optionally prefixed by a `-` to indicate a negative value.
* `string`: A string value. String value syntax is detailed below.
* `id`: A valid identifier. Valid identifier syntax is detailed below.

An example grammar:

```
`<subbable> -> id <sub>
<sub>       -> . <subbable>
             | ε
```

The string `this.is.an.example` is considered valid by this grammar.

### Pseudocode

In several places in this document, a form of pseudocode is used to properly describe the evaluation procedure of some syntactic constructs. The pseudocode is a simplified, typeless, declarationless variant of Go. It looks like the following:

```go
ex1 = f(3)
for i = 0; i < 5; i++ {
  ex1 += i
}
return ex1
```

Lexical Elements
----------------

WDTE syntax is broken into the following whitespace-separated tokens: Strings, numbers, identifiers, keywords. These are described in detail below, along with comments, which are not considered to be a token.

### Strings

Strings are denoted, like in most languages, by quotes. Both single and double quotes are valid, and function identically to each other, but a string must terminate using the same symbol as the one that started it. Strings are allowed to contain newlines without escaping, and some escape characters are also allowed. If a newline is escaped, it will be ignored, meaning that it won't be included in the processed string.

As of writing, the following escape sequences are supported:

* `\n`: Newline.
* `\t`: Tab.

### Numbers

Numbers are any valid floating point number, optionally containing a period, and optionally prefixed with a `-` to indicate a negative number. A `+` prefix is not allowed.

### Identifiers

Identifier rules are very, very loose. Essentially, an identifier is any string of non-whitespace Unicode characters that don't cause ambiguity other tokens. In the case of keywords, an identifier simply may not be the same as the keyword. In the case of symbols, an identifier may not contain the symbol at all. For example:

* `example`: Valid identifier.
* `+`: Valid identifier.
* `let`: Invalid. `let` is a keyword.
* `anidentifierthatcontainsletinit`: Valid identifier.
* `example2`: Valid identifier.
* `2example`: Invalid. Will parse as `2` and `example`.
* `example->with.symbols`: Invalid. Will parse as `example`, `->`, `with`, `.`, and `symbols`.

### Keywords

Keywords are broken into two types, standard and symbolic. The primary difference between the two is the method in which they affect the parsing of identifiers, as shown by the examples above.

#### Standard Keywords

```wdte
switch memo let import
```

#### Symbolic Keywords

```wdte
. -> => { }
  --    [ ]
; :     ( )
        (@
```

### Comments

Comments start with `#` and continue until the end of the line. There are no multi-line comments. `#` is valid in an identifier, but it must not be the first character of the identifier.

Expressions
-----------

### Singles

```
<single> -> number
          | string
          | <import>
          | <lambda>
          | <array>
          | <subbable>
```

A single is any expression that can be determined to be a standalone element, such as a number, string, or identifier on its own. There are several more complex types of singles, which will be explained in further detail below.

### Normal Expressions (Function Calls)

```
<expr> -> <single> <args> <slot> <chain>
<args> -> <single> <args>
        | ε
```

Function calls, confusingly referred to as `expr` for historical reasons, are a single followed by a space-separated list of zero or more singles as arguments, optionally followed by a slot and a chain. Slots and chains work in tandem, and they are further explained in the section on chains.

When a function is called, evaluation proceeds by first evaluating the single on its own. The return value of this is then called, passing it the list of arguments. The return value of this is considered to be the return value of the entire expression.

### Chains

```
<chain> -> -> <expr>
         | -- <expr>
         | ε
```

A chain is a special type of expression that chains multiple expressions together. It does so by performing the following process:

1. Evaluate the first element in the chain.
2. Evaluate the second element.
3. Evaluate the result of the first passed to the result of the second.
4. Evaluate the third element.
5. Evaluate the result of step 3 passed to the result of the third element.
6. Continue until finished.

For example, the chain `ex1 2 3 -> ex2 5 -> ex3` proceeds as follows:

```go
r1 = ex1(2, 3)
r2 = ex2(5)
r1 = r2(r1)
r2 = ex3()
r1 = r2(r1)
return r1
```

#### Ignored Elements

Chains may have elements be ignored, denoted by the use of `--` as the chain operator rather than `->`. If an element is ignored, its result is not used in further processing of the chain. If the above example had been `ex1 2 3 -- ex2 5 -> ex3`, then the processing would have been

```go
r1 = ex1(2, 3)
ex2(5) // ex2 is called, but the result is ignored.
r2 = ex3()
r1 = r2(r1)
return r1
```

#### Slots

```
<slot> -> : id
        | ε
```

An expression may specify a slot. Outside of chains, though legal, this is essentially ignored. In a chain, a slot specifies an identifier to bind the result of the chain expression to for the remainder of the chain, similar to a variable declaration. Slots may shadow existing elements in the scope, including other slots in the same chain. A slot may also be specified for an ignored element, which allows usage of the result of that expression further down the chain.

### Subs

```
<sub>      -> . <subbable>
            | ε
<subbable> -> id <sub>
            | <switch> <sub>
            | <compound> <sub>
```

A sub allows evaluation of a function inside of a subscope. It looks similar to a dot or arrow operator in many other languages, but it functions very differently in some cases.

TODO: More details. Maybe an example? Probably should explain that scopes are values.

# vim:ts=2 sw=2 et
