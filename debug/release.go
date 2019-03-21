//+build !debug

package debug // import "zemn.me/debug"

// Log calls fmt.Logf on the arguments
func Log(string, ...interface{}) {}

// Assert panics with fmt.Sprintf of the arguments
// if the assertion is false
func Assert(assertion bool, fmt string, args ...interface{}) {}
