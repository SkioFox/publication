package test

import "fmt"

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
}
func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}
