package goroutine_test

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func Test() {
	/**
	Go 语言通过go关键字+函数/方法的方式创建一个 goroutine。创建后，新 goroutine 将拥有独立的代码执行流，并与创建它的 goroutine 一起被 Go 运行时调度。
	*/
	go fmt.Println("I am a goroutine")
	var c = make(chan int)
	go func(a, b int) {
		c <- a + b
		close(c) // Close the channel when done 计算完后关闭通道
	}(3, 4)
	// Receive data from the channel and print it
	//for val := range c {
	//	//fmt.Println(val)
	//	fmt.Printf("val is:%v\n", val)
	//}
	//fmt.Printf("c is:%v\nc type is:%T\nc val is:%v\n", c, c, <-c)
	//v := <-c
	//s := <-c
	//fmt.Printf("v is:%v\ns is:%v\n", v, s) // v is:7 s is:0 => 从一个已经关闭的通道中接收值会返回通道类型的零值。这就是为什么在你的示例中，第二次接收值是 0。为了避免这种情况，你通常需要在接收值时检查通道是否已经关闭。
	// 使用 "comma ok" 习惯用法，你可以检查从通道接收的值是否有效，以及通道是否已经关闭。
	// 第一次从通道接收值
	value1, ok1 := <-c
	if ok1 {
		fmt.Println("First value received:", value1) // 输出: First value received: 7
	} else {
		fmt.Println("Channel closed, no value received")
	}

	// 第二次从通道接收值
	value2, ok2 := <-c
	if ok2 {
		fmt.Println("Second value received:", value2)
	} else {
		fmt.Println("Channel closed, no value received") // 输出: Channel closed, no value received
	}
	/**
	创建 goroutine 后，go 关键字不会返回 goroutine id 之类的唯一标识 goroutine 的 id，你也不要尝试去得到这样的 id 并依赖它。
	另外，和线程一样，一个应用内部启动的所有 goroutine 共享进程空间的资源，如果多个 goroutine 访问同一块内存数据，将会存在竞争，我们需要进行 goroutine 间的同步。
	goroutine 的使用代价很低，Go 官方也推荐你多多使用 goroutine。而且，多数情况下，我们不需要考虑对 goroutine 的退出进行控制：goroutine 的执行函数的返回，就意味着 goroutine 退出。
	goroutine 执行的函数或方法即便有返回值，Go 也会忽略这些返回值。所以，如果你要获取 goroutine 执行后的返回值，你需要另行考虑其他方法，比如通过 goroutine 间的通信来实现。
	channel发送数据和主协程中range读取都是按照顺序的
		顺序发送: 在协程中，值是按顺序发送到通道中的。
		顺序接收: 在主协程中，range 循环按顺序接收从通道发送过来的值。
		通道的 FIFO 特性: Go 的通道是先入先出（FIFO）的，因此发送到通道中的值会按顺序被接收。
	*/
}

/*
*
使用 goroutine 处理并发读取：注意竞态条件问题
*/
func TestGoroutineConcurrentRead() {
	ch := make(chan int)

	// 启动一个 goroutine 向通道发送数据 然后关闭通道
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	// 使用 WaitGroup 来等待所有读取操作完成
	var wg sync.WaitGroup
	numReaders := 2
	wg.Add(numReaders)

	// 启动多个 goroutine 并发读取通道数据
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			for value := range ch {
				fmt.Printf("Goroutine %d received: %d\n", id, value)
			}
		}(i)
	}

	// 等待所有读取操作完成
	wg.Wait()
}
func readFromChannel(ch <-chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		value, ok := <-ch
		if !ok {
			fmt.Printf("Goroutine %d: channel closed\n", id)
			return
		}
		/**
		每个 goroutine 接收到的 ch 数据不一样，并且会变化，是由于竞态条件（Race Condition）导致的。竞态条件发生在多个 goroutine 并发读取和写入共享变量时，导致程序的行为无法确定，因为读取和写入操作的顺序是不确定的。
		在上述代码中，每个 goroutine 启动后都会从通道 ch 中接收数据。如果你的程序在多个 goroutine 之间共享了 ch 通道，并且在主 goroutine 中往通道中发送数据，那么就可能会出现竞态条件。
		为了避免竞态条件，可以采取以下措施之一：
			使用互斥锁：在访问共享变量之前，先获取互斥锁，然后释放互斥锁。这样可以确保在同一时间只有一个 goroutine 能够访问共享变量，从而避免竞态条件。
			使用通道进行同步：在多个 goroutine 之间使用通道进行同步操作，这样可以确保在一个 goroutine 完成操作之前，另一个 goroutine 不会开始操作。
			避免共享状态：设计程序时尽量避免共享状态，可以通过将状态封装在单个 goroutine 中，然后使用通道进行通信来避免竞态条件。
		*/
		fmt.Printf("Goroutine %d received: %d\n", id, value)
	}
}
func readFromChannelSelect(ch <-chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				fmt.Printf("Goroutine %d: channel closed\n", id)
				return
			}
			fmt.Printf("Goroutine %d received: %d\n", id, value)
		default:
		}
	}
}

/*
*

	使用 select 语句进行非阻塞读取：注意竞态条件问题
*/
func TestGoroutineConcurrentReadUseSelect() {
	ch := make(chan int)
	var wg sync.WaitGroup
	// 启动一个 goroutine 向通道发送数据
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	// 启动两个 goroutine 并发读取通道数据
	numReaders := 2
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go readFromChannelSelect(ch, i+1, &wg) // 通过传递 sync.WaitGroup 对象的指针，可以确保在函数内部对其进行修改时影响到原始对象，并且减少了内存开销。
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All goroutines finished")
}
func readFromChannelMux(ch <-chan int, id int, wg *sync.WaitGroup, mu *sync.Mutex) {
	// 这里不使用互斥锁就不会死锁
	defer wg.Done()
	for {
		//mu.Lock() // 获取互斥锁
		value, ok := <-ch
		//mu.Unlock() // 释放互斥锁
		if !ok {
			fmt.Printf("Goroutine %d: channel closed\n", id)
			return
		}
		fmt.Printf("Goroutine %d received: %d\n", id, value)
	}
}

/*
*

	死锁例子
*/
func TestGoroutineRaceConditionForReadAndWrite() {
	ch := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 启动一个 goroutine 向通道发送数据
	/**
	发送数据的 goroutine 中使用了互斥锁 mu，导致发送数据和接收数据的操作在同一时刻只能由一个 goroutine 完成，从而造成了死锁。
	要解决这个问题，你可以将发送数据和接收数据的操作分开，或者不使用互斥锁。
	*/
	go func() {
		for i := 0; i < 5; i++ {
			//mu.Lock() // 获取互斥锁
			ch <- i
			//mu.Unlock() // 释放互斥锁
		}
		close(ch)
	}()

	// 启动两个 goroutine 并发读取通道数据
	numReaders := 2
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go readFromChannelMux(ch, i+1, &wg, &mu)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All goroutines finished")
}
func writeToChannel(ch chan<- int, mu *sync.Mutex) {
	/**
	在 writeToChannel 函数中使用互斥锁导致死锁的原因是因为该函数在一个无限循环中不断尝试获取互斥锁，并在获取到锁之后写入数据到通道，但由于通道没有机会在这个循环中被读取，因此通道会一直保持满状态，而无法给写入的数据腾出空间，因此导致了死锁。
	具体来说，在 writeToChannel 函数中的互斥锁是没有必要的，因为写入通道不会出现并发写入的情况。通常来说，互斥锁是用来保护共享资源的并发访问，如果没有并发访问，使用互斥锁反而会增加额外的开销，并且可能引发死锁等问题。
	在修复代码的时候，可以直接删除 writeToChannel 函数中的互斥锁，因为通道本身就是并发安全的，不需要额外的同步机制来保护并发写入操作。
	*/
	for i := 0; i < 5; i++ {
		//mu.Lock() // 获取互斥锁
		ch <- i
		//mu.Unlock() // 释放互斥锁
	}
	close(ch) // close可以避免竞态条件？
}
func TestDeadLock() {
	ch := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 启动一个 goroutine 向通道发送数据
	go writeToChannel(ch, &mu)

	// 启动两个 goroutine 并发读取通道数据
	numReaders := 2
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go readFromChannelMux(ch, i+1, &wg, &mu)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All goroutines finished")
}

/*
*
使用 for range 循环从通道中读取数据时，当通道关闭时，循环会自动退出。这种方式非常简洁明了，适用于需要连续读取通道中的数据，直到通道关闭的情况
*/
func testForRange() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// 启动一个 goroutine 不断向通道发送数据
	wg.Add(1)
	go func() {
		defer close(ch)
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	// 使用 for range 循环读取通道数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		for value := range ch {
			fmt.Println("Received:", value)
		}
		fmt.Println("All data received")
	}()

	wg.Wait()
}

/*
*
这里会死锁
for select 语句：

	使用 for select 语句可以在一个循环中处理多个通道操作。它允许在每次循环迭代时选择处理一个或多个可用的通道操作。
	可以与 case <-ch: 结合使用来读取通道中的数据，并结合其他 case 来处理其他事件或通道操作。
*/
func testForSelectLockError() {
	ch1 := make(chan int)
	ch2 := make(chan string)
	var wg sync.WaitGroup

	// 启动两个 goroutine 向不同的通道发送数据
	/**
	死锁的原因：
		问题出现在两个 goroutine 向通道发送数据的过程中。当其中一个 goroutine 先发送完数据并关闭通道时，另一个 goroutine 仍在尝试发送数据到通道，但通道已经关闭，这样会导致发送数据的 goroutine 发生 panic。
	解决方法
		可以将 wg.Wait() 放置在发送数据的 goroutine 中，以确保所有数据发送完毕后再等待数据接收 goroutine 的完成。
	*/
	wg.Add(2)
	go func() {
		defer close(ch1)
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
	}()
	go func() {
		defer close(ch2)
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch2 <- fmt.Sprintf("Message %d", i)
		}
	}()

	// 使用 for select 语句同时读取多个通道数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case value1, ok := <-ch1:
				/**
				当其中一个发送 goroutine 完成并关闭通道时，接收数据的 goroutine 会继续在 for select 循环中等待另一个通道的数据。如果此时没有新的数据到来，且所有发送 goroutine 都已经完成并关闭了通道，接收数据的 goroutine 会一直阻塞，导致程序死锁。
				所以这里需要对通道进行判空并跳出接收的循环
				*/
				if !ok {
					fmt.Println("Channel 1 closed")
					return
				}
				fmt.Println("Received from Channel 1:", value1)
			case value2, ok := <-ch2:
				if !ok {
					fmt.Println("Channel 2 closed")
					return
				}
				fmt.Println("Received from Channel 2:", value2)
			}
			//if ch1 == nil && ch2 == nil {
			//	break
			//}
		}
	}()

	wg.Wait()
	fmt.Println("All data received")
}
func testForSelectNoLock() {
	//ch1 := make(chan int)
	//ch2 := make(chan string)
	//var wg sync.WaitGroup
	//// 启动两个 goroutine 向不同的通道发送数据
	//wg.Add(2)
	//go func() {
	//	defer close(ch1)
	//	defer wg.Done()
	//	for i := 0; i < 5; i++ {
	//		ch1 <- i
	//	}
	//}()
	//go func() {
	//	defer close(ch2)
	//	defer wg.Done()
	//	for i := 0; i < 5; i++ {
	//		ch2 <- fmt.Sprintf("Message %d", i)
	//	}
	//}()
	//// 使用 for select 语句同时读取多个通道数据
	//for {
	//	select {
	//	case value1, ok := <-ch1:
	//		if !ok {
	//			ch1 = nil // 重点：设置为 nil 以防止再次选择，死锁的原因
	//		} else {
	//			fmt.Println("Received from Channel 1:", value1)
	//		}
	//	case value2, ok := <-ch2:
	//		if !ok {
	//			ch2 = nil // 重点：设置为 nil 以防止再次选择，死锁的原因
	//		} else {
	//			fmt.Println("Received from Channel 2:", value2)
	//		}
	//	}
	//
	//	// 如果两个通道都已经关闭，跳出循环
	//	if ch1 == nil && ch2 == nil {
	//		break
	//	}
	//}
	//wg.Wait()
	//fmt.Println("All data received")
	// 更加完善的写法
	ch1 := make(chan int)
	ch2 := make(chan string)
	var wg sync.WaitGroup

	// 启动两个 goroutine 向不同的通道发送数据
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
		close(ch1) // 发送完成后关闭通道
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch2 <- fmt.Sprintf("Message %d", i)
		}
		close(ch2) // 发送完成后关闭通道
	}()

	// 使用一个 goroutine 读取两个通道的数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		for ch1 != nil || ch2 != nil { // select之前先判空
			select {
			case value1, ok := <-ch1:
				if !ok {
					fmt.Println("Channel 1 closed")
					ch1 = nil // 设置为 nil 以防止再次选择
				} else {
					fmt.Println("Received from Channel 1:", value1)
				}
			case value2, ok := <-ch2:
				if !ok {
					fmt.Println("Channel 2 closed")
					ch2 = nil // 设置为 nil 以防止再次选择
				} else {
					fmt.Println("Received from Channel 2:", value2)
				}
			}
		}
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All data received")
}

/*
*
<- 操作符：
使用 <- 操作符可以单独读取通道中的数据，但不能关闭通道。通常与 ok 模式结合使用，以检查通道是否已关闭。
*/
func testArrowUse() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// 启动一个 goroutine 向通道发送数据
	wg.Add(1)
	go func() {
		defer close(ch)
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	// 使用 <- 操作符单独读取通道数据
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				break
			}
			fmt.Println("Received:", value)
		}
		break
	}

	wg.Wait()
	fmt.Println("All data received")
}
func UseAgeReadChannel() {
	//testForRange()
	//testForSelectLockError()
	//testArrowUse()
	testForSelectNoLock()
}
func spawn(f func() error) <-chan error {
	c := make(chan error)
	go func() {
		c <- f()
	}()
	return c
}
func TestChannelConnect() {
	/**
	这个示例在 main goroutine 与子 goroutine 之间建立了一个元素类型为 error 的 channel，子 goroutine 退出时，会将它执行的函数的错误返回值写入这个 channel，main goroutine 可以通过读取 channel 的值来获取子 goroutine 的退出状态。
	*/
	c := spawn(func() error {
		time.Sleep(2 * time.Second)
		return errors.New("timeout")
	})
	fmt.Println(<-c)
}
func deadloop() {
	for {

	}
}

/*
*
goroutine调度
*/
func TestGoroutiineScheduling() {
	/**
	1、可以定时 1 秒间隔地不断看到“I got scheduled!”输出；
	2、main 函数中，go deadloop() 语句前添加 runtime.GOMAXPROCS(1)，即可使 main goroutine 在创建 deadloop goroutine 之后无法继续得到调度。
	*/
	//runtime.GOMAXPROCS(1)
	go deadloop()
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("I got scheduled!")
	}
}

func workerTest(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	interval, ok := args[0].(int)
	if !ok {
		return
	}

	time.Sleep(time.Second * (time.Duration(interval)))
}

func spawnTest(f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	go func() {
		f(args...)
		c <- struct{}{}
	}()
	return c
}

/*
*
测试goroutine通信
*/
func TestGoroutingSingal() {
	done := spawnTest(workerTest, 10)
	println("spawn a worker goroutine")
	<-done
	println("worker done")
}

var OK = errors.New("ok")

func workerWithError(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	interval, ok := args[0].(int)
	if !ok {
		return errors.New("invalid interval arg")
	}

	time.Sleep(time.Second * (time.Duration(interval)))
	return OK
}

func spawnWithError(f func(args ...interface{}) error, args ...interface{}) chan error {
	c := make(chan error)
	go func() {
		c <- f(args...)
	}()
	return c
}

/*
*
我们将 channel 中承载的类型由struct{}改为了error，这样 channel 承载的信息就不仅仅是一个“信号”了，还携带了“有价值”的信息：新 goroutine 的结束状态。
*/
func TestGoroutingSingalWithError() {
	done := spawnWithError(workerWithError, 5)
	println("spawn worker1")
	err := <-done
	fmt.Println("worker1 done:", err)
	done = spawnWithError(workerWithError)
	println("spawn worker2")
	err = <-done
	fmt.Println("worker2 done:", err)
}
func spawnGroup(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			name := fmt.Sprintf("worker-%d:", i)
			f(args...)
			println(name, "done")
			wg.Done() // worker done!
		}(i)
	}

	go func() {
		wg.Wait()
		c <- struct{}{}
	}()

	return c
}

/*
*
有些场景中，goroutine 的创建者可能会创建不止一个 goroutine，并且需要等待全部新 goroutine 退出。我们可以通过 Go 语言提供的sync.WaitGroup实现等待多个 goroutine 退出的模式
*/
func TestGoroutingSingalWithSyncWaitGroup() {
	done := spawnGroup(5, workerTest, 3)
	println("spawn a group of workers")
	<-done
	println("group workers done")
}

/*
*
超时退出
在下述代码中，我们通过一个定时器(time.Timer)设置了超时等待时间，并通过select原语同时 timer 和done channel，哪个先返回数据就执行哪个 case 分支。
*/
func TestGoroutineTimeout() {
	done := spawnGroup(5, workerTest, 30)
	println("spawn a group of workers")

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-done:
		println("group workers done")
	}
}

/*
*
前面的几个场景中，goroutine 的创建者都是在被动地等待着新 goroutine 的退出。
但很多时候，goroutine 创建者需要主动通知那些新 goroutine 退出，尤其是当 main goroutine 作为创建者时。main goroutine 退出意味着 Go 程序的终止，
而粗暴地直接让 main goroutine 退出的方式可能会导致业务数据的损坏、不完整或丢失。我们可以通过“notify-and-wait（通知并等待）”模式来满足这一场景的要求。
虽然这一模式也不能完全避免“损失”，但是它给了各个 goroutine 一个“挽救数据”的机会，可以尽可能地减少损失的程度。
*/
func worker(j int) {
	time.Sleep(time.Second * (time.Duration(j)))
}

/*
*
示例代码中，使用创建模式创建 goroutine 的spawn函数返回的 channel 的作用发生了变化，从原先的只是用于新 goroutine 发送退出“信号”给创建者，变成了一个双向的数据通道：
既承载创建者发送给新 goroutine 的“退出信号”，也承载新 goroutine 返回给创建者的“退出状态”。
*/
func spawnNotifyWait(f func(int)) chan string {
	quit := make(chan string)
	go func() {
		var job chan int // 模拟job channel
		for {
			select {
			case j := <-job:
				f(j)
			case <-quit:
				quit <- "ok"
			}
		}
	}()
	return quit
}

/*
*通知并等待一个 goroutine 退出
 */
func TestNotifyWait() {
	quit := spawnNotifyWait(worker)
	println("spawn a worker goroutine")

	time.Sleep(5 * time.Second)

	// 通知新创建的goroutine退出
	println("notify the worker to exit...")
	quit <- "exit"

	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	select {
	case status := <-quit:
		println("worker done:", status)
	case <-timer.C:
		println("wait worker exit timeout")
	}
}
func spawnGroupNotifyWait(n int, f func(int)) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done() // 保证wg.Done在goroutine退出前被执行
			name := fmt.Sprintf("worker-%d:", i)
			for {
				j, ok := <-job
				if !ok {
					println(name, "done")
					return
				}
				// do the job
				worker(j)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job) // 广播给所有新goroutine
		wg.Wait()
		quit <- struct{}{}
	}()

	return quit
}

/*
*
通知并等待多个 goroutine 退出
下面是“通知并等待多个 goroutine 退出”的场景。Go 语言的 channel 有一个特性，那就是当使用 close 函数关于 channel 时，所有阻塞到该 channel 上的 goroutine 都会得到“通知”，
我们就利用这一特性实现满足这一场景的模式
*/
func TestNotifyWait1() {
	quit := spawnGroupNotifyWait(5, worker)
	println("spawn a group of workers")

	time.Sleep(1 * time.Second)
	// notify the worker goroutine group to exit
	println("notify the worker group to exit...")
	quit <- struct{}{}

	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-quit:
		println("group workers done")
	}
}
