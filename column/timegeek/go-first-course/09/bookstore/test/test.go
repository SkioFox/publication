package test

import (
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
}

func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}
