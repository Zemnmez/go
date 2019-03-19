package main

import (
	"fmt"
	"time"

	"zemn.me/progress"
)

func main() {
	fmt.Println("let's go!")
	var c progress.Canvas
	c.Draw(progress.SimpleText{
		Pos:  progress.Pos{0, 0},
		Text: "hello world!",
	})

	c.Flush()

	<-time.After(4 * time.Second)
	c.Close()

}
