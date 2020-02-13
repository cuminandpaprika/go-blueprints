package trace

import (
	"fmt"
	"io"
)

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off returns a tracer that does nothing
func Off() Tracer {
	return &nilTracer{}
}

// tracer is the interface that describes an object
// capable of tracing events through the code
// The Tracer interface defines a function Trace that
// will accept zero or more arguments of ANY type
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
