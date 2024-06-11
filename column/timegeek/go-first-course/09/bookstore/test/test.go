package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"sync"
	"time"
	"unsafe"
)

var (
	mu sync.RWMutex
)

func Test() {
	//testByte()
	//testSliceDiff()
	//testMapOrder()
	//testForRange()
	//testForNewRange()
	//testForSyncNewRange()
	//testForOrderRange()
	//testForRightOrderRange()
	//testArraySliceMapRange()
	//testSwitchCase()
	//testFunc()
	//testError()
	//testCustomError()
	//testFmt()
	//testErrIs()
	testErrAs()
}
func testByte() {
	// 测试byte字节格式
	s1 := []byte("你好，世界！")
	fmt.Println("byte:", s1) // 111111111 [228 189 160 229 165 189 239 188 140 228 184 150 231 149 140 239 188 129]  UTF-8 编码的字符串的字节序列
	str := string(s1)
	fmt.Println("string:", str) // "你好，世界！
	// test string
	var strTest = "中国人"
	fmt.Printf("the length of s = %d\n", len(strTest)) // 9
	for i := 0; i < len(strTest); i++ {
		fmt.Printf("ox%x ", strTest[i]) // oxe4 oxb8 oxad oxe5 ox9b oxbd oxe4 oxba oxba
	}

	var s = "hello"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // 将string类型变量地址显式转
	fmt.Printf("0x%x\n", hdr.Data)                     // 0x10a30e0
	p := (*[5]byte)(unsafe.Pointer(hdr.Data))          // 获取Data字段所指向的数组的指针
	//data := unsafe.StringData(s)          // 获取字符串数据指针:unsafe.StringData requires go1.20 or later (-lang was set to go1.17; check go.mod)
	//fmt.Printf("0x%x\n", data)            // 输出数据指针的地址
	//p := (*[5]byte)(unsafe.Pointer(data)) // 将数据指针转换为字节数组指针
	dumpBytesArray((*p)[:]) // [h e l l o ]   // 输出底层数组的内容
}
func testSliceDiff() {
	// 对比下面两个切片的区别
	var sl1 []int
	var sl2 = []int{}
	fmt.Print("========基本区别=========\n")
	fmt.Printf("%v,len:%d,cap:%d,addr:%p\n", sl1, len(sl1), cap(sl1), &sl1)
	fmt.Printf("%v,len:%d,cap:%d,addr:%p\n", sl2, len(sl2), cap(sl2), &sl2)
	fmt.Printf("sl1==nil:%v\n", sl1 == nil)
	fmt.Printf("sl2==nil:%v\n", sl2 == nil)
	a1 := *(*[3]int)(unsafe.Pointer(&sl1))
	a2 := *(*[3]int)(unsafe.Pointer(&sl2))
	fmt.Print("========底层区别=========\n")
	fmt.Println(a1)
	fmt.Println(a2)

	type SliceDemo struct {
		Values []int
	}
	var s5 = SliceDemo{}
	var s6 = SliceDemo{[]int{}}
	bs1, _ := json.Marshal(s5)
	bs2, _ := json.Marshal(s6)
	fmt.Print("========序列化区别=========\n")
	fmt.Println(a1)
	fmt.Println(string(bs1))
	fmt.Println(string(bs2))

	// sl1是声明，还没初始化，是nil值，底层没有分配内存空间。
	// sl2初始化了，不是nil值，底层分配了内存空间，有地址。
	// https://qcrao.com/2019/04/02/dive-into-go-slice/
	// https://tonybai.com/2022/02/15/whether-go-allocate-underlying-array-for-empty-slice/
}
func testMapOrder() {
	// 测试map遍历与顺序 => 程序逻辑千万不要依赖遍历 map 所得到的的元素次序
	m := map[int]int{
		2: 12,
		1: 11,
		3: 13,
	}
	for i := 0; i < 3; i++ {
		doIteration(m)
	}
	// 顺序
	doIterationMap(m)
	doWrite(m)
	doWrite(m)
	doWrite(m)
}
func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}
func doIteration(m map[int]int) {
	fmt.Printf("{ ")
	for k, v := range m {
		fmt.Printf("[%d, %d] ", k, v)
	}
	fmt.Printf("}\n")
}
func doIterationMap(m map[int]int) {
	mu.RLock()
	defer mu.RUnlock()
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		log.Fatalf("Error serializing to JSON: %v", err)
	}
	fmt.Println("keys => ", jsonKeys)
	fmt.Printf("keys=>%s \n", jsonKeys)
	fmt.Println("keys=>%s", string(jsonKeys))
	fmt.Printf("keys=>%s", string(jsonKeys))
	sort.SliceStable(keys, func(x, y int) bool {
		return x < y
	})
	for _, k := range keys {
		fmt.Printf("[%d, %d] ", k, m[k])
	}
	fmt.Println()
}
func doWrite(m map[int]int) {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range m {
		m[k] = v + 1
	}
}
func testForRange() {
	var m = []int{1, 2, 3, 4, 5}
	for i, v := range m {
		go func() {
			// 全部输出4 5 => Goroutine 执行的闭包函数引用了它的外层包裹函数中的变量 i、v，这样，变量 i、v 在主 Goroutine 和新启动的 Goroutine 之间实现了共享，
			// 而 i, v 值在整个循环过程中是重用的，仅有一份。在 for range 循环结束后，i = 4, v = 5，因此各个 Goroutine 在等待 3 秒后进行输出的时候，输出的是 i, v 的最终值。
			fmt.Println("testForRange：", i, v)
		}()
	}
	time.Sleep(time.Second * 5) // 这里的Sleep作用 => 确保所有 Goroutine 有时间完成输出 更好的做法是使用sync.WaitGroup
}
func testForNewRange() {
	var m = []int{1, 2, 3, 4, 5}
	{
		i, v := 0, 0
		for i, v = range m {
			go func() {
				fmt.Println("testForNewRange：", i, v) // 全部输出4 5
			}()
		}
	}
	time.Sleep(time.Second * 5)
}
func testForSyncNewRange() {
	// Goroutine 并不是在 range 循环之后才执行的。实际上，Goroutine 是在 range 循环的每次迭代中启动的。之所以看起来像是在 range 循环之后才执行，是因为 Goroutine 的调度和主 Goroutine 的执行速度之间的关系导致的。
	// 当你在 range 循环的每次迭代中启动一个 Goroutine 时，这些 Goroutine 会被调度器安排去执行。由于 Go 的调度器是并发的，并且每个 Goroutine 的启动和执行并不是立即的，它们可能在主 Goroutine 执行完 range 循环后才开始运行。
	// Goroutine 并不是在 range 循环之后才执行的，而是由于调度和执行的异步特性，导致看起来像是 range 循环结束后才执行。正确地捕获循环变量值和使用同步机制（如 sync.WaitGroup）可以确保 Goroutine 按预期执行并完成任务。
	// 方法一：传递参数给 Goroutine
	var m = []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup
	//for i, v := range m {
	//	wg.Add(1) // 增加计数器
	//	go func(i, v int) {
	//		defer wg.Done()                            // Goroutine 完成时减少计数器
	//		fmt.Println("testForRightNewRange:", i, v) // 即使我们在循环中按顺序启动 Goroutine，打印的顺序仍然可能是不确定的。这是因为每个 Goroutine 的启动和执行是并发的，具体的执行顺序由 Go 的调度器决定。
	//		// 在这个输出中，Goroutine 的执行顺序并没有严格按照我们启动的顺序进行。这是因为 Goroutine 的调度是非确定性的，有可能某个 Goroutine 比其他 Goroutine 先获得执行时间片，从而导致输出顺序看似混乱。
	//		/**
	//		影响因素
	//		几个影响 Goroutine 执行顺序的因素包括：
	//
	//		操作系统调度：操作系统调度器可能会在不同时间点切换不同的 Goroutine。
	//		CPU负载：当前系统的CPU负载会影响Goroutine的调度。
	//		Goroutine数量：较多的Goroutine会增加调度的不确定性。
	//		Go调度器：Go 运行时调度器本身会根据多种因素调度Goroutine。
	//		确保顺序输出
	//		如果确实需要按顺序执行和输出，可以通过其他方式实现，例如在主 Goroutine 中按顺序启动和等待每个 Goroutine 完成，或者使用有序的通道通信。
	//		*/
	//	}(i, v)
	//}
	// 方法二：使用局部变量或者说闭包特性
	for i, v := range m {
		i, v := i, v // 重新声明并定义局部变量
		wg.Add(1)    // 增加计数器
		go func() {
			defer wg.Done()
			fmt.Println("testForSyncNewRange:", i, v)
		}()
	}
	wg.Wait() // 等待所有 Goroutine 完成
}
func testForOrderRange() {
	/**
	几个影响 Goroutine 执行顺序的因素包括：
		操作系统调度：操作系统调度器可能会在不同时间点切换不同的 Goroutine。
		CPU负载：当前系统的CPU负载会影响Goroutine的调度。
		Goroutine数量：较多的Goroutine会增加调度的不确定性。
		Go调度器：Go 运行时调度器本身会根据多种因素调度Goroutine。
	*/
	var m = []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup
	ch := make(chan string, len(m))

	for i, v := range m {
		wg.Add(1)
		go func(i, v int) {
			defer wg.Done()
			ch <- fmt.Sprintf("testForOrderRange: %d %d", i, v)
		}(i, v)
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Read and print messages from the channel in order
	for msg := range ch {
		fmt.Println(msg) // 这里的顺序还是随机的，取决于goroutine的调度 => 虽然每个 Goroutine 按顺序发送消息到通道，但是由于 Goroutine 是并发执行的，其执行顺序是不确定的，因此接收到的消息顺序也可能是不确定的。
	}
}
func testForRightOrderRange() {
	var m = []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup
	// 使用一个切片 results 来存储每个 Goroutine 的输出，并确保按顺序打印这些输出。
	results := make([]string, len(m))

	for i, v := range m {
		wg.Add(1)
		go func(i, v int) {
			defer wg.Done()
			results[i] = fmt.Sprintf("testForRightOrderRange: %d %d", i, v)
		}(i, v)
	}

	wg.Wait() // 等待所有 Goroutine 完成

	for _, result := range results {
		fmt.Println(result)
	}
}
func testArraySliceMapRange() {
	testArray()
	testMap()
	testMap1()
}
func testArray() {
	// 参与 for range 循环的是 range 表达式的副本
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int
	fmt.Println("original a =", a)
	//for i, v := range a {
	//for i, v := range a[:] { // 用切片代替数组
	for i, v := range &a { // 用数组指针代理数组
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}
	fmt.Println("after for range loop, r =", r)
	fmt.Println("after for range loop, a =", a)
}
func testMap() {
	// 如果我们在循环的过程中，对 map 进行了修改，那么这样修改的结果是否会影响后续迭代呢？这个结果和我们遍历 map 一样，具有随机性。
	// 我们日常编码遇到遍历 map 的同时，还需要对 map 进行修改的场景的时候，要格外小心。
	var m = map[string]int{"tony": 21, "tom": 22, "jim": 23}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "tony")
		}
		counter++
		fmt.Println(k, v) // 反复运行这个例子多次，会得到两个不同的结果。
	}
	fmt.Println("counter is ", counter) // 反复运行这个例子多次，会得到两个不同的结果。
}
func testMap1() {
	var m = map[string]int{"tony": 21, "tom": 22, "jim": 23}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			m["lucy"] = 24
		}
		counter++
		fmt.Println(k, v) // 反复运行这个例子多次，会得到两个不同的结果。
	}
	fmt.Println("counter is ", counter) // 反复运行这个例子多次，会得到两个不同的结果。
}

type I interface {
	M()
}
type T struct{}

func (T) M() {

}
func testSwitchCase() {
	// type switch
	var x interface{} = 13 // 你可以发现，在前面的 type switch 演示示例中，我们一直使用 interface{}这种接口类型的变量，Go 中所有类型都实现了 interface{}类型，所以 case 后面可以是任意类型信息。 但如果在 switch 后面使用了某个特定的接口类型 I，那么 case 后面就只能使用实现了接口类型 I 的类型了，否则 Go 编译器会报错。
	switch v := x.(type) {
	case nil:
		println("v is nil")
	case int:
		println("the type of v is int, v =", v)
	case string:
		println("the type of v is string, v =", v)
	case bool:
		println("the type of v is bool, v =", v)
	default:
		println("don't support the type")
	}
	var t T
	var i I = t
	/**
	在前面的 type switch 演示示例中，我们一直使用 interface{}这种接口类型的变量，Go 中所有类型都实现了 interface{}类型，所以 case 后面可以是任意类型信息。
	但如果在 switch 后面使用了某个特定的接口类型 I，那么 case 后面就只能使用实现了接口类型 I 的类型了，否则 Go 编译器会报错。
	*/
	switch i.(type) {
	case T:
		println("it is type T")
		//case int:
		//	println("it is type int")
		//case string:
		//	println("it is type string")
		/**
		在这个例子中，我们在 type switch 中使用了自定义的接口类型 I。那么，理论上所有 case 后面的类型都只能是实现了接口 I 的类型。但在这段代码中，只有类型 T 实现了接口类型 I，Go 原生类型 int 与 string 都没有实现接口 I，于是在编译上述代码时，编译器会报出如下错误信息：
		impossible type switch case: i (type I) cannot have dynamic type int
		impossible type switch case: i (type I) cannot have dynamic type string
		*/
	}
	// 跳不出循环的 break(不带 label 的 break 语句中断执行并跳出的，是同一函数内 break 语句所在的最内层的 for、switch 或 select。)
	var sl = []int{5, 19, 6, 3, 8, 12}
	var firstEven int = -1
	// find first even number of the interger slice
	for i := 0; i < len(sl); i++ {
		switch sl[i] % 2 {
		case 0:
			firstEven = sl[i]
			break
		case 1:
			// do nothing
		}

	}
	// 这就是 Go 中 break 语句与 switch 分支结合使用会出现一个“小坑”。和我们习惯的 C 家族语言中的 break 不同，Go 语言规范中明确规定，不带 label 的 break 语句中断执行并跳出的，是同一函数内 break 语句所在的最内层的 for、switch 或 select。
	// 所以，上面这个例子的 break 语句实际上只跳出了 switch 语句，并没有跳出外层的 for 循环，这也就是程序未按我们预期执行的原因。
	println(firstEven) // 切片中的第一个偶数是 6，而输出的结果却成了切片的最后一个偶数 12
	// 优化使用loop
	var s2 = []int{5, 19, 6, 3, 8, 12}
	var firstEven1 int = -1
	// find first even number of the interger slice
loop:
	for i := 0; i < len(s2); i++ {
		switch s2[i] % 2 {
		case 0:
			firstEven1 = s2[i]
			break loop
		case 1:
			// do nothing
		}
	}
	println(firstEven1) // 6
}
func myAppend(sl []int, elems ...int) []int {
	fmt.Printf("%T\n", elems) // []int
	if len(elems) == 0 {
		println("no elems to append")
		return sl
	}
	sl = append(sl, elems...)
	return sl
}
func testFunc() {
	sl := []int{1, 2, 3}
	sl = myAppend(sl) // no elems to append
	fmt.Println(sl)   // [1 2 3]
	sl = myAppend(sl, 4, 5, 6)
	fmt.Println(sl) // [1 2 3 4 5 6]
}
func doSomething(flag bool) error {
	if !flag {
		return errors.New("something went wrong")
	}
	return nil
}
func testError() {
	// Go 标准库中有很多预定义的错误类型和创建错误的方法，例如 errors.New 和 fmt.Errorf。
	err := doSomething(false)
	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Success")
	}
}

// 定义一个自定义错误类型
type MyError struct {
	Code    int
	Message string
}

// 实现 error 接口
func (e MyError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// 包装错误
func doSomething1(flag bool) error {
	if !flag {
		// 在 Go 语言中，fmt.Errorf 是一个用于格式化和创建错误对象的函数。它类似于 fmt.Sprintf，但返回的是一个实现了 error 接口的错误对象。这
		// 使得它非常适合用于创建包含详细信息的错误消息。
		return fmt.Errorf("doSomething failed: %w", MyError{Code: 500, Message: "Internal Server Error"})
	}
	return nil
}
func testCustomError() {
	err := doSomething1(false)
	if err != nil {
		fmt.Println("Error occurred:", err)
		var myErr MyError
		if errors.As(err, &myErr) {
			// fmt.Sprintf 是 Go 语言标准库中的一个函数，用于格式化字符串。它返回一个格式化后的字符串，而不像 fmt.Printf 那样直接输出结果。
			// fmt.Sprintf 类似于 C 语言中的 sprintf 函数，可以通过指定格式化动词将不同类型的变量格式化为字符串。
			/**
			主要作用
				格式化字符串：将各种类型的数据格式化为字符串。
				字符串拼接：方便地将多个值拼接成一个字符串。
				生成带有动态数据的字符串：根据动态数据生成自定义的字符串。
			*/
			fmt.Println("Custom error code:", myErr.Code)
		}
	} else {
		fmt.Println("Success")
	}
}
func testFmt() {
	// Println - 自动换行
	fmt.Println("Hello, world!")

	// Printf - 格式化输出
	fmt.Printf("Hello, %s!\n", "Alice")
	fmt.Printf("Number: %d\n", 42)

	// Print - 不换行
	fmt.Print("Hello, ")
	fmt.Print("world!\n")

	// Sprintf - 格式化字符串，不打印
	s := fmt.Sprintf("Hello, %s!", "Bob")
	fmt.Println(s)

	// Scanf - 从输入读取格式化数据
	var age int
	fmt.Print("Enter your age: ")
	fmt.Scanf("%d", &age)
	fmt.Println("Your age is", age)

	// Errorf - 创建一个格式化的错误消息
	err := fmt.Errorf("an error occurred: %s", "file not found")
	fmt.Println(err)
}

var ErrNotFound = errors.New("resource not found")

func findResource() error {
	return fmt.Errorf("failed to find resource: %w", ErrNotFound)
}

type MyError1 struct {
	Message string
}

func (e MyError1) Error() string {
	return e.Message
}

func (e MyError1) Is(target error) bool {
	_, ok := target.(MyError1)
	return ok
}
func testErrIs() {
	// errors.Is 是 Go 1.13 版本引入的一个用于错误处理的函数。它主要用于判断一个错误是否是特定类型的错误或是由特定类型的错误包装而来的。这在处理复杂错误链时非常有用。
	// 判断简单错误
	err := ErrNotFound
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error is ErrNotFound")
	}
	// 判断包装错误
	err1 := findResource()
	if errors.Is(err1, ErrNotFound) {
		fmt.Println("Error is ErrNotFound")
	} else {
		fmt.Println("Different error")
	}
	// 判断自定义错误类型
	err2 := MyError1{Message: "custom error"}
	if errors.Is(err2, MyError1{}) {
		fmt.Println("Error is of type MyError")
	}
	/**
	总结
		1. errors.Is 用于判断一个错误是否等于或包含某个特定的错误。
		2. 可以判断简单错误和包装错误。
		3. 通过 fmt.Errorf 的 %w 动词，可以将错误包装起来，以保留错误链。
		4. 自定义错误类型可以通过实现 Is 方法来支持 errors.Is 的判断。
	*/
}

type MyError2 struct {
	Message string
}

func (e MyError2) Error() string {
	return e.Message
}

func (e MyError2) Unwrap() error {
	return errors.New("underlying error")
}

func doSomething3() error {
	return fmt.Errorf("doSomething failed: %w", MyError2{Message: "something went wrong"})
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e NotFoundError) Unwrap() error {
	return errors.New("underlying not found error")
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func doSomething4() error {
	notFoundErr := NotFoundError{Message: "resource not found"}
	return fmt.Errorf("validation failed: %w", fmt.Errorf("wrapped error: %w", notFoundErr))
}
func testErrAs() {
	// errors.As 是 Go 1.13 版本引入的一个用于错误处理的函数。它主要用于将错误链中的某个错误类型提取出来，以便进一步处理。这在处理嵌套或包装的错误时非常有用。
	// errors.As 用于将错误链中的某个特定类型的错误提取出来。如果找到该类型的错误，则会将其赋值给目标变量，并返回 true；否则返回 false。
	err := doSomething3()
	var myErr MyError2
	if errors.As(err, &myErr) {
		fmt.Println("Caught MyError:", myErr)
	} else {
		fmt.Println("Different error")
	}
	// 检查嵌套错误
	err4 := doSomething4()
	var notFoundErr NotFoundError
	var validationErr ValidationError

	if errors.As(err4, &notFoundErr) {
		fmt.Println("Caught NotFoundError:", notFoundErr)
	} else if errors.As(err, &validationErr) {
		fmt.Println("Caught ValidationError:", validationErr)
	} else {
		fmt.Println("Different error")
	}
	/**
	总结
		1. errors.As 用于将错误链中的某个特定类型的错误提取出来。
		2. 如果找到该类型的错误，则会将其赋值给目标变量，并返回 true；否则返回 false。
		3. 通过 errors.As，可以更方便地处理复杂的错误链，并根据具体的错误类型执行不同的逻辑。
		4. 这种机制对于处理嵌套错误和包装错误非常有用，使得错误处理代码更加清晰和简洁。
	*/
}
