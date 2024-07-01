package main

import (
	"fmt"
	"time"
)

type counter struct {
	c chan int
	i int
}

var cter counter

func InitCounter() {
	cter = counter{
		c: make(chan int),
	}

	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	fmt.Println("counter init ok")
}

func Increase() int {
	return <-cter.c
}

func init() {
	InitCounter()
}

/*
*
无缓冲 channel 具有同步特性，这让它在某些场合可以替代锁，从而使得程序更加清晰，可读性更好。下面是一个传统的基于“共享内存”+“锁”模式的 goroutine 安全的计数器的实现

在这个实现中，我们将计数器操作全部交给一个独立的 goroutine 去处理，并通过无缓冲 channel 的同步阻塞特性实现计数器的控制。
这样其他 goroutine 通过 Increase 函数试图增加计数器值的动作实质上就转化为一次无缓冲 channel 的接收动作。这种并发设计逻辑更符合 Go 语言所倡导的**“不要通过共享内存来通信，
而是通过通信来共享内存”**的原则。
*/
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			v := Increase()
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
		}(i)
	}

	time.Sleep(5 * time.Second)
}
