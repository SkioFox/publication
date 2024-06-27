package generics_test

import "fmt"

func TestGenericsConstraint[T1, T2 any](t1 T1, t2 T2) T1 {
	var a T1 // 声明变量
	var b T2
	a, b = t1, t2 // 同类型赋值
	_ = b

	f := func(t T1) {
	}
	f(a) // 传给其他函数

	p := &a // 取变量地址
	_ = p

	var i interface{} = a // 转换或赋值给interface{}类型变量
	_ = i

	c := new(T1) // 传递给预定义函数
	_ = c

	f(a) // 将变量传给其他函数

	sl := make([]T1, 0, 10) // 作为复合类型中的元素类型
	_ = sl

	j, ok := i.(T1) // 用在类型断言中
	_ = ok
	_ = j

	switch i.(type) { // 作为type switch中的case类型
	case T1:
	case T2:
	}
	return a // 从函数返回
}

// comparable.go

type foo struct {
	a int
	s string
}

type bar struct {
	a  int
	sl []string
}

func doSomething[T comparable](t T) T {
	var a T
	if a == t {
	}

	if a != t {
	}
	return a
}

func TestGenericsConstraint1() {
	doSomething(true)
	doSomething(3)
	doSomething(3.14)
	doSomething(3 + 4i)
	doSomething("hello")
	var p *int
	doSomething(p)
	doSomething(make(chan int))
	doSomething([3]int{1, 2, 3})
	doSomething(foo{})
	//doSomething(bar{}) //  bar does not implement comparable
}

// stringify.go

func Stringify0[T fmt.Stringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String())
	}
	return ret
}

type MyString0 string

func (s MyString0) String() string {
	return string(s)
}

/*
*
自定义约束
这个例子中，我们使用的是fmt.Stringer接口作为约束。一方面，这要求类型参数T的实参必须实现fmt.Stringer接口的所有方法；
另一方面，泛型函数Stringify0的实现代码中，声明的T类型实例（比如v）也仅被允许调用fmt.Stringer的String方法。
*/
func TestCstConstraint() {
	sl := Stringify0([]MyString0{"I", "love", "golang"})
	fmt.Println(sl) // 输出：[I love golang]
}

/*
*我们自定义了一个Stringer接口类型作为约束。在该类型中，我们不仅定义了String方法，还嵌入了comparable，这样在泛型函数中，我们用Stringer约束的类型参数就具备了进行相等性和不等性比较的能力了！
 */
type Stringer interface {
	ordered
	comparable
	String() string
}
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

func StringifyWithoutZero[T Stringer](s []T) (ret []string) {
	var zero T
	for _, v := range s {
		if v == zero {
			continue
		}
		ret = append(ret, v.String())
	}
	return ret
}

type MyString string

func (s MyString) String() string {
	return string(s)
}
func StringifyLessThan[T Stringer](s []T, max T) (ret []string) {
	var zero T
	for _, v := range s {
		if v == zero || v >= max {
			continue
		}
		ret = append(ret, v.String())
	}
	return ret
}
func TestCstConstraint1() {
	//sl := StringifyWithoutZero([]MyString{"I", "", "love", "", "golang"}) // 输出：[I love golang]
	sl := StringifyLessThan([]MyString{"I", "", "love", "", "golang"}, MyString("cpp")) // 输出：[I]
	fmt.Println(sl)
}

type BasicInterface interface { // 基本接口类型
	M1()
}

type NonBasicInterface interface { // 非基本接口类型
	BasicInterface
	~int | ~string // 包含类型元素
}

type MyString1 string

func (MyString1) M1() {
}

func foo1[T NonBasicInterface](a T) { // 非基本接口类型作为约束
}

func bar1[T BasicInterface](a T) { // 基本接口类型作为约束
}

func TestGenericsInterface() {
	var s = MyString1("hello")
	var bi BasicInterface = s // 基本接口类型支持常规用法
	//var nbi NonBasicInterface = s // 非基本接口不支持常规用法，导致编译器错误：cannot use type NonBasicInterface outside a type constraint: interface contains type constraints
	bi.M1()
	//nbi.M1()
	foo1(s)
	bar1(s)
}

type Intf1 interface {
	~int | string
	F1()
	F2()
}

type Intf2 interface {
	~int | ~float64
}

type I interface {
	Intf1
	M1()
	M2()
	int | ~string | Intf2
}

func doSomething1[T I](t T) {
}

type MyInt int

func (MyInt) F1() {
}
func (MyInt) F2() {
}
func (MyInt) M1() {
}
func (MyInt) M2() {
}

/*
*
测试类型集合
*/
func TestTypeSetGenerics() {
	var a int = 11
	//doSomething1(a) // int does not implement I (missing F1 method)

	var b = MyInt(a)
	doSomething1(b) // ok
}
