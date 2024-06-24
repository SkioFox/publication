//package foo
//
//import (
//	"testing"
//)
//
//// for send benchmark test
//var c1 chan string
//
//// for recv benchmark test
//var c2 chan string
//
//func init() {
//	c1 = make(chan string)
//	go func() {
//		for {
//			<-c1
//		}
//	}()
//
//	c2 = make(chan string)
//	go func() {
//		for {
//			c2 <- "hello"
//		}
//	}()
//}
//
//func send(msg string) {
//	c1 <- msg
//}
//func recv() {
//	<-c2
//}
//
///*
//*
//单接收单发送性能的基准测试
//
//	go语言中基准测试用于测量代码性能。testing.B 是标准库 testing 包中用于基准测试的类型，而 b.N 是基准测试运行的迭代次数。
//
//名字以 Benchmark 开头，这是基准测试函数的命名约定，便于 go test 命令识别。
//*/
//func BenchmarkUnbufferedChan1To1Send(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		send("hello")
//	}
//}
//func BenchmarkUnbufferedChan1To1Recv(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		recv()
//	}
//}
