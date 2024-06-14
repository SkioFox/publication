package main

import "github.com/SkioFox/publication/column/timegeek/go-first-course/27/instrument_trace"

func foo() {
	defer trace.Trace()()
	bar()
}

func bar() {
	defer trace.Trace()()
}

func main() {
	defer trace.Trace()()
	foo()
}
