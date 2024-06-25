package main

import (
	"fmt"
	"time"

	"github.com/bigwhite/workerpool"
)

func main() {
	p := workerpool.New(5)
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			/**
			问题：在Schedule函数追加了print函数，
			fmt.Printf("%s: Scheduling task\n", time.Now())
			发现所有的Task并不是一次性加入到pool的
			当pool的capacity=5时，会先加入6个task然后等有task结束后再Schedule剩余的task，请问这是被阻塞了吗？
			解答：tasks是不带缓冲的channel，如果5个worker都在处理task，那么没有空闲的worker去从tasks中读取，这时schedule新task会阻塞。直到有空闲worker。
			*/
			time.Sleep(time.Second * 3)
			fmt.Printf("%s: Scheduling task\n", time.Now())
		})
		if err != nil {
			println("task: ", i, "err:", err)
		}
	}

	p.Free()
}
