package interface_test

import "fmt"

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
