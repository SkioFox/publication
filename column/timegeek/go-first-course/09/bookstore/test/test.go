package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

func Test() {
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

func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}
