package goroutine_test

import (
	"fmt"
	"sync"
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
for select 语句：

	使用 for select 语句可以在一个循环中处理多个通道操作。它允许在每次循环迭代时选择处理一个或多个可用的通道操作。
	可以与 case <-ch: 结合使用来读取通道中的数据，并结合其他 case 来处理其他事件或通道操作。
*/
func testForSelect() {
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
		}
	}()

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
	testForSelect()
	//testArrowUse()
}
