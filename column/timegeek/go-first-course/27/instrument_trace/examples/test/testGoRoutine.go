package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var (
	mu sync.Mutex
	m  = make(map[uint64]int)
)

func TestGoRoutine() {
	defer Trace()()
	var wg sync.WaitGroup
	wg.Add(1)
	/**
	当使用 go 关键字启动一个新的 goroutine 时，后面需要跟一个函数调用，这个函数调用可以是普通函数、方法，也可以是一个匿名函数。匿名函数在定义后可以直接调用，形成一个自执行函数。
	defer也是一样 经常用自执行的匿名函数作为函数表达式跟在后面 defer是在该函数注册到defered函数栈的时候进行求值(遇到函数返回函数时就能追踪到调用函数栈， 因为注册的时候求值会执行该表达式，返回一个函数，该函数在函数退出前执行)，并在函数退出前执行。
	*/
	go func() {
		A2()
		wg.Done()
	}()
	A1()
	wg.Wait()
}

/*
*
单 Goroutine 改为多 Goroutine 并发的，这样才能验证支持多 Goroutine 的新版 Trace 函数是否好用
Goroutine ID 我们可以快速确认某一行输出是属于哪个 Goroutine 的。
*/
func Trace() func() {
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
	gid := curGoroutineID()
	/**
	从上面代码看到，我们在出入口输出的跟踪信息中加入了 Goroutine ID 信息，我们输出的 Goroutine ID 为 5 位数字，如果 ID 值不足 5 位，则左补零，这一切都是 Printf 函数的格式控制字符串“%05d”帮助我们实现的。
	这样对齐 Goroutine ID 的位数，为的是输出信息格式的一致性更好。如果你的 Go 程序中 Goroutine 的数量超过了 5 位数可以表示的数值范围，也可以自行调整控制字符串。
	*/
	//fmt.Printf("g[%05d]: enter: [%s]\n", gid, name)
	//return func() { fmt.Printf("g[%05d]: exit: [%s]\n", gid, name) }

	/**
	使用了一个 map 类型变量 m 来保存每个 Goroutine 当前的缩进信息：m 的 key 为 Goroutine 的 ID，值为缩进的层次。然后，考虑到 Trace 函数可能在并发环境中运行，
	考虑“map 不支持并发写”的注意事项，我们增加了一个 sync.Mutex 实例 mu 用于同步对 m 的写操作。
	*/
	mu.Lock()
	indents := m[gid]    // 获取当前gid对应的缩进层次
	m[gid] = indents + 1 // 缩进层次+1后存入map
	mu.Unlock()
	printTrace(gid, name, "->", indents+1)
	return func() {
		mu.Lock()
		indents := m[gid]    // 获取当前gid对应的缩进层次
		m[gid] = indents - 1 // 缩进层次-1后存入map
		mu.Unlock()
		printTrace(gid, name, "<-", indents)
	}
}
func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "    "
	}
	fmt.Printf("g[%05d]:%s%s%s\n", id, indents, arrow, name)
}

func curGoroutineID() uint64 {
	var goroutineSpace = []byte("goroutine ")

	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	// Parse the 4707 out of "goroutine 4707 ["  b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}
func A1() {
	defer Trace()()
	B1()
}
func B1() {
	defer Trace()()
	C1()
}
func C1() {
	defer Trace()()
	D()
}
func D() {
	defer Trace()()
}
func A2() {
	defer Trace()()
	B2()
}
func B2() {
	defer Trace()()
	C2()
}
func C2() {
	defer Trace()()
	D()
}
