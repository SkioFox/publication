package channel_test

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*
*
生产者只能向 channel 中发送数据，我们使用chan<- int作为 produce 函数的参数类型
*/
func produce(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i + 1
		time.Sleep(time.Second)
	}
	close(ch)
}

/*
*
消费者只能从 channel 中接收数据，我们使用<-chan int作为 consume 函数的参数类型
*/
func consume(ch <-chan int) {
	/**
	使用了 for range 循环语句来从 channel 中接收数据，for range 会阻塞在对 channel 的接收操作上，直到 channel 中有数据可接收或 channel 被关闭循环，才会继续向下执行。
	channel 被关闭后，for range 循环也就结束了。

	通过“comma, ok”惯用法或 for range 语句，我们可以准确地判定 channel 是否被关闭。而单纯采用n := <-ch形式的语句，我们就无法判定从 ch 返回的元素类型零值，究竟是不是因为 channel 被关闭后才返回的。
	*/
	//for n := range ch {
	//	println(n)
	//}
	for {
		n1 := <-ch
		fmt.Println(n1) // 无法停止会一直打印
	}
	//for {
	//	m, ok := <-ch
	//	if !ok {
	//		fmt.Println("Channel closed")
	//		break
	//	}
	//	fmt.Println(m, ok)
	//}
}
func TestChannelType() {
	//ch1 := make(chan<- int, 1) // 只发送channel类型
	//ch2 := make(<-chan int, 1) // 只接收channel类型
	// <-ch1 // Invalid operation: <-ch1 (receive from the send-only type chan<- int)
	//ch2 <- 13 // Invalid operation: ch2 <- 13 (send (send to the receive-only type <-chan int)

	ch := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		produce(ch)
		wg.Done()
	}()
	go func() {
		consume(ch)
		wg.Done()
	}()
	wg.Wait()

	// 多个channel使用for select
}

type signal struct{}

func worker() {
	println("worker is working...")
	time.Sleep(1 * time.Second)
}
func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("worker start to work...")
		f()
		c <- signal{}
	}()
	return c
}

/*
*
用作信号通知
*/
func TestChannelUnBuffered() {
	println("start a worker...")
	c := spawn(worker)
	/**
	在这个例子中，spawn 函数返回的 channel，被用于承载新 Goroutine 退出的“通知信号”，这个信号专门用作通知 main goroutine。main goroutine 在调用 spawn 函数后一直阻塞在对这
	个“通知信号”的接收动作上。
	*/
	<-c
	fmt.Println("worker work done!")
}

/*
*
1 对 n 的信号通知机制, 这样的信号通知机制，常被用于协调多个 Goroutine 一起工作
*/
func workern(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}
func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal // 等待开始工作的信号
			fmt.Printf("worker %d: start to work...\n", i)
			f(i) // 执行工作函数
			wg.Done()
		}(i + 1)
	}
	go func() {
		wg.Wait()               // 等待所有工作者完成工作
		c <- signal(struct{}{}) // 向信号通道发送完成信号
		close(c)                // 关闭信号通道
	}()
	return c // 返回信号通道
}

/*
*
关闭一个无缓冲 channel 会让所有阻塞在这个 channel 上的接收操作返回，从而实现了一种 1 对 n 的“广播”机制。
*/
func TestChannelUnBuffered1() {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal) // 创建一个信号通道
	c := spawnGroup(workern, 5, groupSignal)
	time.Sleep(5 * time.Second) // 模拟一些准备时间
	fmt.Println("the group of workers start to work...")
	close(groupSignal) // 关闭信号通道，通知所有工作者开始工作 =>  close一个channel后，所有阻塞在这个channel接收操作的goroutine都会收到通知，这是Go语言的channel语义就这么定义的。
	<-c                // 等待所有工作者完成工作的信号
	fmt.Println("the group of workers work done!")
}

// 替代锁机制
type counter struct {
	sync.Mutex
	i int
}

var cter counter

func Increase() int {
	cter.Lock()
	defer cter.Unlock()
	cter.i++
	return cter.i
}
func TestChannelUnBuffered1Lock() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			v := Increase()
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

type counter1 struct {
	c chan int
	i int
}

func NewCounter() *counter1 {
	cter := &counter1{
		c: make(chan int),
	}
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	return cter
}
func (cter *counter1) Increase() int {
	return <-cter.c
}

/*
*
我们将计数器操作全部交给一个独立的 Goroutine 去处理，并通过无缓冲 channel 的同步阻塞特性，实现了计数器的控制。的同步阻塞特性，实现了计数器的控制。这样其他 Goroutine 通过 Increase 函数试图增加计数器值的动作，实质上就转化为了一次无缓冲 channel 的接收动作。
这种并发设计逻辑更符合 Go 语言所倡导的“不要通过共享内存来通信，而是通过通信来共享内存”的原则。
*/
func TestChannelUnBuffered1Lock1() {
	cter := NewCounter()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			v := cter.Increase()
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("cter type is:%T\n", cter)
}

/*
*带缓冲channel用作计数信号量（counting semaphore）
 */

func TestChannelUseCountingSemaphore() {
	var active = make(chan struct{}, 3)
	var jobs = make(chan int, 10)
	go func() {
		/**
		请问计数信号量的例子中，因为jobs的容量是10，这里执行的循环不会导致阻塞，close(jobs)  应该会被执行到，那么下面的for range为什么不会终止，而可以继续运行？
		channel内部数据是排队的，即便被close，依然可以从closed channel中读取到尚未被消费的元素，直到没有可读的元素为止，才真正会变成closed状态。没数据后，
		如果再读就会得到元素类型的零值了, 对于没数据且closed状态的channel，for range会终止。
		*/
		for i := 0; i < 8; i++ {
			jobs <- i + 1
		}
		close(jobs)
	}()
	var wg sync.WaitGroup
	for j := range jobs {
		wg.Add(1)
		go func(j int) {
			active <- struct{}{}
			log.Printf("handle job: %d\n", j) // 同一时间允许最多 3 个 Goroutine 处于活动状态。
			time.Sleep(2 * time.Second)
			<-active
			wg.Done()
		}(j)
	}
	wg.Wait()
}

type signal1 struct{}

var ready1 bool

func worker1(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}
func spawnGroup1(f func(i int), num int, mu *sync.Mutex) <-chan signal1 {
	c := make(chan signal1)
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				mu.Lock()
				if !ready1 {
					mu.Unlock()
					time.Sleep(100 * time.Millisecond)
					continue
				}
				mu.Unlock()
				fmt.Printf("worker %d: start to work...\n", i)
				f(i)
				wg.Done()
				return
			}
		}(i + 1)
	}
	go func() {
		wg.Wait()
		c <- signal1(struct{}{})
	}()
	return c
}
func TestCondition() {
	fmt.Println("start a group of workers...")
	mu := &sync.Mutex{}
	c := spawnGroup1(worker1, 5, mu)
	time.Sleep(5 * time.Second) // 模拟ready前的准备工作  fmt.Println("the group of workers start to work...")
	mu.Lock()
	ready1 = true
	mu.Unlock()
	<-c
	fmt.Println("the group of workers work done!")
}

type signal2 struct{}

var ready2 bool

func worker2(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}
func spawnGroup2(f func(i int), num int, groupSignal *sync.Cond) <-chan signal2 {
	c := make(chan signal2)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			groupSignal.L.Lock()
			for !ready2 {
				groupSignal.Wait()
			}
			groupSignal.L.Unlock()
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}
	go func() {
		wg.Wait()
		c <- signal2(struct{}{})
	}()
	return c
}

/*
*
我们看到，sync.Cond实例的初始化，需要一个满足实现了sync.Locker接口的类型实例，通常我们使用sync.Mutex。
条件变量需要这个互斥锁来同步临界区，保护用作条件的数据。加锁后，各个等待条件成立的
Goroutine 判断条件是否成立，如果不成立，则调用sync.Cond的 Wait 方法进入等待状态。 Wait 方法在 Goroutine 挂起前会进行 Unlock 操作。
当 main goroutine 将ready置为 true，并调用sync.Cond的 Broadcast 方法后，各个阻塞的 Goroutine 将被唤醒，并从 Wait 方法中返回。
Wait 方法返回前，Wait 方法会再次加锁让 Goroutine 进入临界区。接下来 Goroutine 会再次对条件数据进行判定，如果条件成立，就会解锁并进入下一个工作阶段；如果条件依旧不成立，
那么会再次进入循环体，并调用 Wait 方法挂起等待。
*/
func TestCondition1() {
	fmt.Println("start a group of workers...")
	groupSignal2 := sync.NewCond(&sync.Mutex{})
	c := spawnGroup2(worker2, 5, groupSignal2)
	time.Sleep(5 * time.Second) // 模拟ready前的准备工作  fmt.Println("the group of workers start to work...")
	groupSignal2.L.Lock()
	ready2 = true
	groupSignal2.Broadcast()
	groupSignal2.L.Unlock()
	<-c
	fmt.Println("the group of workers work done!")
}
func op1(mu1, mu2 *sync.Mutex, wg *sync.WaitGroup) {
	mu1.Lock()
	time.Sleep(1 * time.Second)
	mu2.Lock()
	println("op1: do something...")
	mu2.Unlock()
	mu1.Unlock()
	wg.Done()
}
func op2(mu1, mu2 *sync.Mutex, wg *sync.WaitGroup) {
	mu2.Lock()
	time.Sleep(1 * time.Second)
	mu1.Lock()
	println("op1: do something...")
	mu1.Unlock()
	mu2.Unlock()
	wg.Done()
}

/*
*

	op1 和 op2 互相等待对方释放锁，从而形成死锁。
	具体步骤如下：

	op1 锁定 mu1。
	op2 锁定 mu2。
	op1 等待锁定 mu2，但此时 mu2 已被 op2 锁定。
	op2 等待锁定 mu1，但此时 mu1 已被 op1 锁定。
	结果就是两个 goroutine 互相等待对方释放锁，导致死锁。

解决死锁的方法
1. 固定加锁顺序：确保所有 goroutine 以相同的顺序锁定多个互斥锁。这样可以避免循环等待的发生。
2.使用 TryLock 方法：在某些编程语言和库中，可以尝试使用 TryLock 方法，避免死锁情况。Go 的 sync.Mutex 没有内置 TryLock 方法，可以通过一些更复杂的逻辑模拟，但这种方法通常不是首选。
3.减少锁的持有时间：尽量减少锁的持有时间，将耗时操作放在锁外执行。
*/
func TestDeadLock2() {
	var mu1 sync.Mutex
	var mu2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	go op1(&mu1, &mu2, &wg)
	go op2(&mu1, &mu2, &wg)
	wg.Wait()
}
