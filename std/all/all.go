// Package all is a convenience package that imports the entire
// standard library, thus registering it with std.Import.
package all

//go:generate bash gen.bash -o gen.go -p all
