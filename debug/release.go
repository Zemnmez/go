//+build !debug

package debug // import "zemn.me/debug"

// This file does not contain debug-time code!
// check ./debug.go :)

// Log calls fmt.Logf on the arguments
func Log(string, ...interface{}) {}

// Assert panics with fmt.Sprintf of the arguments
// if the assertion is false
func Assert(assertion bool, fmt string, args ...interface{}) {}
