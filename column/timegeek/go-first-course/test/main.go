package main

import channel_test "github.com/publication/column/timegeek/go-first-course/test/channel-test"

func main() {
	//fmt.Println("test")
	//interface_test.TestInterFace()
	//interface_test.TestInterFace1()
	//interface_test.TestInterFace2()
	//interface_test.TestinterFaceDuck()
	//interface_test.TestNilErr()
	//interface_test.TestErrNilDiffNil()
	//interface_test.TestBoxing()
	/**
	类似node koa中的经典洋葱模型
	*/
	//interface_test.TestMiddlewareModel()
	//goroutine_test.Test()
	//goroutine_test.TestGoroutineConcurrentRead()
	//goroutine_test.TestGoroutineConcurrentReadUseSelect()
	/**
	下面两个例子都是一样的会死锁
	*/
	//goroutine_test.TestGoroutineRaceConditionForReadAndWrite()
	//goroutine_test.TestDeadLock()
	//goroutine_test.UseAgeReadChannel()
	//goroutine_test.TestChannelConnect()
	//goroutine_test.TestGoroutiineScheduling()
	//channel_test.TestChannelType()
	//channel_test.TestChannelUnBuffered()
	//channel_test.TestChannelUnBuffered1()
	//channel_test.TestChannelUnBuffered1Lock()
	//channel_test.TestChannelUnBuffered1Lock1()
	//channel_test.TestChannelUseCountingSemaphore()
	//channel_test.TestCondition()
	//channel_test.TestCondition1()
	channel_test.TestDeadLock2()
}
