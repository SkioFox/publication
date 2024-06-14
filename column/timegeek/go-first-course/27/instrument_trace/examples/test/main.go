package main

import (
	"runtime"
)

/*
*
testTrace单个线程
*/
func main() {
	// defer 关键字后面的表达式，是在将 deferred 函数注册到 deferred 函数栈的时候进行求值的。
	// 对于自定义的函数或方法，defer 可以给与无条件的支持，但是对于有返回值的自定义函数或方法，返回值会在 deferred 函数被调度执行的时候被自动丢弃。
	// 使用 defer 可以跟踪函数的执行过程
	//defer Trace("main")()
	/**
	由于 TraceNew() 中的 println("enter:", name) 语句是在匿名函数被创建时执行的，而 println("exit:", name) 是在匿名函数被调用时执行的，所以你会在输出中看到 "enter" 先于 "exit" 输出。
		TraceNew()()第一个匿名函数的创建和执行是在 defer 语句执行时发生的，而第二个匿名函数的创建是在第一个匿名函数被 defer 调用时发生的，但它的执行是在包含 defer 语句的函数执行完毕后。
	*/
	//defer TraceNew()()
	//foo()
	TestGoRoutine()
}

/**
返回函数类型的函数
*/
//	func Trace(name string) func() {
//		println("enter:", name)
//		return func() {
//			println("exit:", name)
//		}
//	}
func TraceNew() func() {
	/**
	通过 runtime.Caller 函数获得当前 Goroutine 的函数调用栈上的信息，runtime.Caller 的参数标识的是要获取的是哪一个栈帧的信息。当参数为 0 时，返回的是 Caller 函数的调用者的函数信息，在这里就是 Trace 函数。
	但我们需要的是 Trace 函数的调用者的信息，于是我们传入 1。
	*/
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	println("enter:", name)
	return func() {
		println("exit:", name)
	}
}
func foo() {
	//defer Trace("foo")()
	defer TraceNew()()
	bar()
}
func bar() {
	//defer Trace("bar")()
	defer TraceNew()()
}
