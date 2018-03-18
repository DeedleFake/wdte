wdte
====

wdte is a command-line interpreter for the WDTE scripting language. It provides a basic WDTE environment to run scripts in. Execution is starts in std.Scope with a custom importer. The importer provides full access to the standard library, as well as a few custom features.

Importer
--------

The custom importer provides two main features over std.Import:

* It provides a `cli` module which gives access to features that make sense from the command-line, such as arguments.
* It provides a means of importing scripts from the file system in much the same way that languages like Python or Ruby do.

### cli

The `cli` module provides the following functions:

#### args

```wdte
args
```

args returns an array containing the arguments passed to the interpreter on the command-line, starting with the path to the current script.

### File System Imports

An import from the file system is attempted if the import string starts with either a `.` or a `/`. If this is true, an import is attempted of the Go plugin at `<import string>.so`. If that doesn't exist, the script at `<import string>.wdte` is tried instead.

If a plugin is found, a function with the signature `func S() *wdte.Scope` is called in it. This is expected to return the module's scope.

If a script is found, the script is parsed with the same importer as the interpreter itself uses, except that any file system imports attempted by that script will be relative to that script's path.

Note that, due to limitations in the Go `plugin` package, plugin imports currently only work on Linux.
