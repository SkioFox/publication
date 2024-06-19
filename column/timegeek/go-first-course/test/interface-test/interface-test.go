package interface_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func TestInterFace() {
	var a int64 = 13
	var i interface{} = a
	v1, ok := i.(int64)
	fmt.Printf("v1=%d, the type of v1 is %T, ok=%t\n", v1, v1, ok) // v1=13,  the type of v1 is int64, ok=true
	v2, ok := i.(string)
	fmt.Printf("v2=%s, the type of v2 is %T, ok=%t\n", v2, v2, ok) //  the type of v2 is string, ok=false
	v3 := i.(int64)
	fmt.Printf("v3=%d, the type of v3 is %T\n", v3, v3) // v3=13, the type of v3 is int64
	//v4 := i.([]int) // panic: interface conversion: interface {} is int64, not []int
	//fmt.Printf("the type of v4 is %T\n", v4)
	/**
	断言 i 的值实现了接口类型 T。如果断言成功，变量 v 的类型为 i 的值的类型，而并非接口类型 T。如果断言失败，v 的类型信息为接口类型 T，它的值为 nil
	*/
}

type MyInterface interface {
	M1()
}
type T int

func (T) M1() {
	println("T's M1")
}
func TestInterFace1() {
	var t T
	//var i interface{} = t
	// 同上 go 1.18新增的用于表示空接口类型的any关键词
	var i any = t
	v1, ok := i.(MyInterface)
	if !ok {
		panic("the value of i is not MyInterface")
	}
	v1.M1()
	fmt.Printf("v1=%d, the type of v1 is %T, ok is %t\n", v1, v1, ok) // the type of v1 is interface_test.T, ok is true
	i = int64(13)
	v2, ok := i.(MyInterface)                                         // 当类型断言失败时，v2 被赋值为 nil，但 v2 的类型仍然是 MyInterface。这是因为类型断言的结果 v2 是声明为接口类型 MyInterface 的一个变量，只是它的值为 nil。
	fmt.Printf("v2=%v, the type of v2 is %T, ok is %t\n", v2, v2, ok) // v2=<nil>, the type of v2 is <nil>, ok is false
	/**
	类型断言失败时，不会改变声明变量的类型。声明变量的类型是在编译时确定的，而类型断言操作是在运行时进行的。即使类型断言失败，变量仍然保持其声明时的类型。
	让我们再深入一步解释这个概念：
		类型断言成功：类型断言成功时，变量会被转换为指定的类型。
		类型断言失败：类型断言失败时，变量的值为 nil，但变量的类型仍然是指定的类型。
	*/
	// v2 type is MyInterface
	//v2 = 13 // cannot use 13 (constant of type int) as MyInterface value in assignment: int does not implement MyInterface (missing method M1)
}
func TestInterFace2() {
	var err error
	err = errors.New("error1")
	/**
	接口的这种“动静皆备”的特性，又带来了什么好处?
		接口类型变量在程序运行时可以被赋值为不同的动态类型变量，每次赋值后，接口类型变量中存储的动态类型信息都会发生变化，这让 Go 语言可以像动态语言（比如 Python）那样拥有使用Duck Typing（鸭子类型）的灵活性。
		所谓鸭子类型，就是指某类型所表现出的特性（比如是否可以作为某接口类型的右值），不是由其基因（比如 C++ 中的父类）决定的，而是由类型所表现出来的行为（比如类型拥有的方法）决定的。
	*/
	fmt.Printf("%T\n", err) // *errors.errorString

}

type QuackableAnimal interface {
	Quack()
}
type Duck struct{}

func (Duck) Quack() {
	println("duck quack!")
}

type Dog struct{}

func (Dog) Quack() {
	println("dog quack!")
}

type Bird struct{}

func (Bird) Quack() {
	println("bird quack!")
}
func AnimalQuackInForest(a QuackableAnimal) {
	a.Quack()
	fmt.Printf("a type is:%T\n", a)
}
func TestinterFaceDuck() {
	/**
	这个例子中，我们用接口类型 QuackableAnimal 来代表具有“会叫”这一特征的动物，而 Duck、Bird 和 Dog 类型各自都具有这样的特征，于是我们可以将这三个类型的变量赋值给 QuackableAnimal 接口类型变量 a。每次赋值，变量 a 中存储的动态类型信息都不同，Quack 方法的执行结果将根据变量 a 中存储的动态类型信息而定。
	这里的 Duck、Bird、Dog 都是“鸭子类型”，但它们之间并没有什么联系，之所以能作为右值赋值给 QuackableAnimal 类型变量，只是因为他们表现出了 QuackableAnimal 所要求的特征罢了。
	不过，与动态语言不同的是，Go 接口还可以保证“动态特性”使用时的安全性。比如，编译器在编译期就可以捕捉到将 int 类型变量传给 QuackableAnimal 接口类型变量这样的明显错误，决不会让这样的错误遗漏到运行时才被发现。
	*/
	animals := []QuackableAnimal{new(Duck), new(Dog), new(Bird)}
	for _, animal := range animals {
		AnimalQuackInForest(animal)
	}
}

type MyError struct {
	error
}

var ErrBad = MyError{
	error: errors.New("bad things happened"),
}

func bad() bool {
	return false
}

/*
*
这里返回值类型从error改为*MyError就能正常输出ok了 为什么？=> go语言中的接口类型的自动装箱机制：将任意类型赋值给一个接口类型变量就是装箱操作,接口类型的装箱实际就是创建一个 eface 或 iface 的过程
1）	returnsError() 函数不返回 error 非空接口类型，而是直接返回结构体指针 *MyError（明确的类型，阻止自动装箱）；
2）	不要直接 err != nil 这样判断，而是使用类型断言来判断：

	if e, ok := err.(*MyError); ok && e != nil {
		fmt.Printf("error occur: %+v\n", e)
		return
	}

PS：Go 的“接口”在编程中需要特别注意，必须搞清楚接口类型变量在运行时的表示，以避免踩坑！！！

这里出现了自动装箱的过程：

	返回类型是error， 是一个接口， p是*MyError， p的方法列表覆盖了error这个接口， 所以是可以赋值给error类型的变量。
	这个过程发生了隐式转换，赋值给接口类型，做装箱创建iface，
	p != nil就成了 (&tab, 0x0) != (0x0, 0x0)
*/
func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = &ErrBad
	}
	return p
}
func TestNilErr() {
	err := returnsError()
	if err != nil {
		/**
		在 Go 语言的 fmt 包中，格式化动词 %v 用于通用的值表示。具体来说：
			%v：值的默认格式。
			%+v：结构体时，会添加字段名。
		*/
		fmt.Printf("error occur: %+v\n", err)
		return
	}
	fmt.Println("ok")
}

/*
*
定义了一个新的类型 T1，它是 int 类型的别名。然后，我们为这个类型定义了一个方法 Error，它返回一个字符串，这使得 T1 类型实现了 error 接口。
将基本类型 int 定义为别名类型 T1 并为其实现 Error 方法，我们使得 T1 类型可以用作错误类型 error。这样我们就可以在需要错误处理的地方使用 T1 类型的值。
*/
type T1 int

func (t T1) Error() string {
	return "bad error"
}

/*
*
由于 eface 和 iface 是 runtime 包中的非导出结构体定义，我们不能直接在包外使用，所以也就无法直接访问到两个结构体中的数据。
不过，Go 语言提供了 println 预定义函数，可以用来输出 eface 或 iface 的两个指针字段的值。
而针对 eface 和 iface 类型的打印函数实现如下：
// $GOROOT/src/runtime/print.go

	func printeface(e eface) {
		 print("(", e._type, ",", e.data, ")")
	}

	func printeface(i iface) {
		 print("(", i.tab, ",", e.data, ")")
	}

println 函数可以输出各类接口类型变量的内部表示信息，并结合输出结果，解析接口类型变量的等值比较操作。
*/
func printNilInterface() {
	// 第一种nil接口变量
	//var i interface{}                 // 空接口类型
	//var err error                     // 非空接口类型
	//println(i)                        // (0x0,0x0)
	//println(err)                      // (0x0,0x0)
	//println("i = nil:", i == nil)     // i = nil: true
	//println("err = nil:", err == nil) // err = nil: true
	//println("i = err:", i == err)     // err = nil: true
	/**
	我们看到，无论是空接口类型还是非空接口类型变量，一旦变量值为 nil，那么它们内部表示
	均为(0x0,0x0)，也就是类型信息、数据值信息均为空。因此上面的变量 i 和 err 等值判断为
	true。结论：未赋初值的接口类型变量的值为 nil，这类变量也就是 nil 接口变量
	*/
	// 第二种：空接口类型变量
	//var eif1 interface{} // 空接口类型
	//var eif2 interface{} // 空接口类型
	//var n, m int = 17, 18
	//eif1 = n
	//eif2 = m
	//println("eif1:", eif1)
	//println("eif2:", eif2)
	//println("eif1 = eif2:", eif1 == eif2) // false
	//eif2 = 17
	//println("eif1:", eif1)
	//println("eif2:", eif2)
	//println("eif1 = eif2:", eif1 == eif2) // true
	//eif2 = int64(17)
	//println("eif1:", eif1)
	//println("eif2:", eif2)
	//println("eif1 = eif2:", eif1 == eif2) // false
	/**
		从输出结果中我们可以总结一下：对于空接口类型变量，只有 _type 和 data 所指数据内容一致的情况下，两个空接口类型变量之间才能划等号。
	另外，Go 在创建 eface 时一般会为 data 重新分配新内存空间，将动态类型变量的值复制到这块内存空间，并将 data 指针指向这块内存空间。
	因此我们多数情况下看到的 data 指针值都是不同的。
	*/
	// 第三种：非空接口类型变量
	//var err1 error // 非空接口类型
	//var err2 error // 非空接口类型
	//println("err1:", err1)
	//println("err1 = nil:", err1 == nil)
	//// 将 nil 转换为 *T1 类型并赋值给 err1
	///**
	//这种写法的主要目的是显式地将一个接口变量设置为具体类型的零值。在 Go 中，这种技巧主要用于保持接口变量的类型信息，同时确保其值为 nil，这在类型断言和错误处理等场景中非常有用。
	//*/
	//err1 = (*T1)(nil)
	//println("err1:", err1)
	//println("err1 = nil:", err1 == nil)
	//err1 = T1(5)
	//err2 = T1(6)
	//println("err1:", err1)
	//println("err2:", err2)
	//println("err1 = err2:", err1 == err2)
	//err2 = fmt.Errorf("%d\n", 5)
	//println("err1:", err1)
	//println("err2:", err2)
	//println("err1 = err2:", err1 == err2)
	// 第四种：空接口类型变量与非空接口类型变量的等值比较
	var eif interface{} = T1(5) // 将值 5 转换为 T1 类型，并赋值给 eif 变量
	var err error = T1(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
	err = T1(6)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
	/**
		你可以看到，空接口类型变量和非空接口类型变量内部表示的结构有所不同（第一个字段：_type vs. tab)，两者似乎一定不能相等。
	但 Go 在进行等值比较时，类型比较使用的是 eface 的 _type 和 iface 的 tab._type，因此就像我们在这个例子中看到的那样，当 eif 和 err 都被赋值为T(5)时，两者之间是划等号的。
	*/
}
func TestErrNilDiffNil() {
	printNilInterface()
}

type T2 struct {
	n int
	s string
}

func (T2) M1() {}
func (T2) M2() {}

type NonEmptyInterface interface {
	M1()
	M2()
}

func TestBoxing() {
	/**
	在 Go 语言中，将任意类型赋值给一个接口类型变量也是装箱操作。有了前面对接口类型变量内部表示的学习，我们知道接口类型的装箱实际就是创建一个 eface 或 iface 的过程。
	对 ei 和 i 两个接口类型变量的赋值都会触发装箱操作
	*/
	var t = T2{
		n: 17,
		s: "hello, interface",
	}
	var ei interface{}
	ei = t
	var i NonEmptyInterface
	i = t
	fmt.Println(ei)
	fmt.Println(i)
	fmt.Println("ei=i:", ei == i)

	var n int = 61
	var newei interface{} = n
	n = 62                             // n的值已经改变
	fmt.Println("data in box:", newei) // 输出仍是61
	/**
	经过装箱后，箱内的数据，也就是存放在新分配的内存空间中的数据与原变量便无瓜葛了
	*/
}

/*
*
中间件 authHandler 起到了对 HTTP 请求进行鉴权的作用
*/
func validateAuth(s string) error {
	if s != "123456" {
		return fmt.Errorf("%s", "bad auth token")
	}
	return nil
}
func greetings(w http.ResponseWriter, r *http.Request) {
	log.Println("222222")
	fmt.Fprintf(w, "Welcome!")
}

/*
*
中间件logHandler打印日志信息
*/
func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		log.Println("111111111")
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t)
		h.ServeHTTP(w, r)
	})
}
func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validateAuth(r.Header.Get("auth"))
		if err != nil {
			http.Error(w, "bad auth param", http.StatusUnauthorized)
			return
		}
		/**
		执行 h.ServeHTTP(w, r) 的作用是调用传入的 http.Handler 接口的 ServeHTTP 方法，实际上执行的是你传入的具体处理器的逻辑。=> 当执行 h.ServeHTTP(w, r) 时，实际上是调用了 greetings(w, r)。
		*/
		h.ServeHTTP(w, r)
	})
}
func TestMiddlewareModel() {
	err := http.ListenAndServe(":8080", logHandler(authHandler(http.HandlerFunc(greetings))))
	if err != nil {
		log.Printf("ServeError:", err)
	}
}
