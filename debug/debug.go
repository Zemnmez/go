//+build debug

package debug //import "zemn.me/debug"

import "log"

func Log(fmt string, args ...interface{}) { log.Printf(fmt, args...) }
func Assert(test bool, fmt string, args ...interface{}) {
	if !test {
		panic(fmt.Sprintf(fmt, args...))
	}
}
