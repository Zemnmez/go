/*
Package progress contains functions for tracking and displaying
heirachical program progress

CONCEPT

Tracking progress is really important in long-running, debuggable programs.
Usually the software just dumps a bunch of data to the buffer,
or shows a throbber that provides no information on crash.

I wanted to do better than that.

With progress, progress can be tracked simply using a context.Context
and it will automatically be funneled to any available terminals
or other outputs.

WRITER BASED TRACKING

You can track progress via completion of units, or display a throbber.
However, this ain't your usual throbber! The throbber takes an io.Writer
and displays the last written line beside it. The throbber only updates
whenever a new line is written, which gives a visual indication
the program hasn't failed.

..  deleting files...
... removing subdirectories

If this process fails, instead of leaving it at that, the entire
buffered input from the io.Writer is dumped to the terminal.

XXX deleting files...
removing subdirectories
doing a thing
XXX

CONTEXT PASSING

The progress API exposes itself via the context.Context API, which
allows it to be passed along API boundaries.

	ctx := progress.WithContext(ctx)

A progress context can derive a progress tracker to send progress data.
If progress is not available, a dummy progress tracker is returned instead
to simplify code.

	func HandleThing(ctx context.Context) {
		w := progress.TextStream(ctx)
		defer w.Close()

		cmd := exec.Command("curl", "https://google.com")
		cmd.Stdout = w
		cmd.Run()
	}

	func HandleThing(ctx context.Context) {
		const lim = 1000
		ctr := progress.Counter(ctx, lim)
		defer ctr.Close()

		for i := 0; i < lim; i++ {
			ctr.Add(1)
		}
	}


Progress can be tracked heirachically, which gives important contextual
information. In this case, "handling thing..." is a counter tracker
whose HandleOtherThing subprocesses are children.

	func HandleThing(ctx context.Context) {
		const lim = 1000
		ctr := progress.Counter(ctx, lim)
		subCtx := progress.Describe("handling thing...")

		ctr := progress.Counter(ctx, lim)
		defer ctr.Close()

		for i := 0; i < lim; i++ {
			ctr.Add(1)
			HandleOtherThing(subCtx)
		}
	}


*/
package progress // import "zemn.me/progress"

type Animation []string

var Dots Animation = []string{
	"   ",
	".  ",
	".. ",
	"...",
	" ..",
	"  .",
}

type contextKeyType struct{}

var contextKey = contextKeyType{}
