package main

var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}

/*
*
对于无缓冲 channel 而言，我们得到以下结论：

	发送动作一定发生在接收动作完成之前；
	接收动作一定发生在发送动作完成之前。

下面的代码可以保证main输出的变量 a 的值为"hello, world"，因为函数 f 中的 channel 接收动作发生在主 goroutine 对 channel 发送动作完成之前，而a = "hello, world"语句又发生在 channel 接收动作之前，因此主 goroutine 在 channel 发送操作完成后看到的变量 a 的值一定是"hello, world"，而不是空字符串。
*/
func main() {
	go f()
	c <- 5
	println(a)
}
