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
//	c1 = make(chan string, 100)
//	for i := 0; i < 10; i++ {
//		go func() {
//			for {
//				<-c1
//			}
//		}()
//		go func() {
//			for {
//				c1 <- "hello"
//			}
//		}()
//	}
//
//	c2 = make(chan string, 100)
//	for i := 0; i < 10; i++ {
//		go func() {
//			for {
//				c2 <- "hello"
//			}
//		}()
//		go func() {
//			for {
//				<-c2
//			}
//		}()
//	}
//}
//
//func send(msg string) {
//	c1 <- msg
//}
//func recv() {
//	<-c2
//}
//
//func BenchmarkBufferedChanNToNSendCap100(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		send("hello")
//	}
//}
//func BenchmarkBufferedChanNToNRecvCap100(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		recv()
//	}
//}
