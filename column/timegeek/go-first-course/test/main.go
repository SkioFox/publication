package main

import goroutine_test "github.com/publication/column/timegeek/go-first-course/test/goroutine-test"

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
	goroutine_test.UseAgeReadChannel()
}
