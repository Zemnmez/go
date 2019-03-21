//+build debug

package debug //import "zemn.me/debug"

import (
	"fmt"
	"log"
)

func Log(fmt string, args ...interface{}) { log.Printf(fmt, args...) }
func Assert(test bool, format string, args ...interface{}) {
	if !test {
		panic(fmt.Sprintf(fmt, args...))
	}
}
