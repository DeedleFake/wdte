export default {
	desc: `Introduction
============

Welcome to the WDTE playground, a browser based evaluation environment for WDTE. This playground's features includes the standard function set as well as a number of importable modules.

If you have never seen WDTE before and are completely confused at the moment, try reading the overview on the WDTE project's wiki: https://github.com/DeedleFake/wdte/wiki

Fun Fact
--------

The WDTE interpreter has been compiled to WebAssembly for this playground, meaning that, by opening this page, you've downloaded the entire system. Congratulations.

Documentation
-------------

For documentation on the standard function set, see https://godoc.org/github.com/DeedleFake/wdte/std

The standard library is available for importing, with the exception of the \`io/file\` module. The \`io\` module is pre-inserted into the initial scope as \`io\`. There is also a \`playground\` module which provides interaction with the playground. It is detailed below.

Playground Module
-----------------

#### wdteVersion
    wdteVersion
Returns the version of WDTE that the playground is using.

#### goVersion
    goVersion
Returns the version of Go that the playground was built with.`,

	input: `io.stdout -> io.writeln 'Greetings, pocket universe.';`,
}
