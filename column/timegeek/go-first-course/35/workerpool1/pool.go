package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool") // workerpool中任务已满，没有空闲goroutine用于处理新任务
	ErrWorkerPoolFreed    = errors.New("workerpool freed")       // workerpool已终止运行
)

type Pool struct {
	capacity int // workerpool大小

	active chan struct{}
	tasks  chan Task

	wg   sync.WaitGroup
	quit chan struct{}
}

type Task func()

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

/*
*
用于创建一个 pool 类型实例，并将 pool 池的 worker 管理机制运行起来；
*/
func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	fmt.Printf("workerpool start\n")

	go p.run()

	return p
}

/*
*
创建新的 worker goroutine
*/
func (p *Pool) newWorker(i int) {
	/**
		在创建一个新的 worker goroutine 之前，newWorker 方法会先调用 p.wg.Add 方法将 WaitGroup 的等待计数加一。
	由于每个 worker 运行于一个独立的 Goroutine 中， newWorker 方法通过 go 关键字创建了一个新的 Goroutine 作为 worker。
	*/
	p.wg.Add(1)
	go func() {
		/**
		这个 defer 块的作用是确保无论 goroutine 是正常结束还是因为 panic 而结束，它都能正确地从 p.active 接收一个信号，以便通知工作池该 worker 已经结束。
		在新 worker 中，为了防止用户提交的 task 抛出 panic，进而导致整个 workerpool 受到影响，我们在 worker 代码的开始处，使用了 defer+recover 对 panic 进行捕捉，捕捉后 worker 也是
		要退出的，于是我们还通过<-p.active更新了 worker 计数器。并且一旦 worker goroutine 退出，p.wg.Done 也需要被调用，这样可以减少 WaitGroup 的 Goroutine 等待数量。

		worker一旦创建后，除了panic和quit通知退出，worker是不会退出的，也就是没有所谓“正常退出” 的情况。所以没在defer中调用<-p.active。

		*/
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", i)
		/**
		新 worker 的核心，依然是一个基于 for-select 模式的循环语句，在循环体中，新 worker 通过 select 监视 quit 和 tasks 两个 channel。和前面的 run 方法一样，当接收到来自 quit channel
		的退出“信号”时，这个 worker 就会结束运行。tasks channel 中放置的是用户通过 Schedule 方法提交的请求，新 worker 会从这个 channel 中获取最新的 Task 并运行这个 Task。
		*/
		for {
			select {
			case <-p.quit: // 通知 worker 退出
				fmt.Printf("worker[%03d]: exit\n", i)
				/**
				这一步的作用是确保在 worker 退出之前，主 goroutine 或工作池管理器已经发送了一个信号到 p.active，以协调和同步 worker 的退出。
				defer 块中的 <-p.active 只有在 worker goroutine 结束时才会执行，而 select 中的 <-p.quit 只有在收到退出信号时才会执行。这两者的执行时机是互斥的，不会重叠。
				*/
				<-p.active // 从 p.active 通道接收一个值的。这种用法一般用于信号通知或同步机制。
				return
			case t := <-p.tasks: // 用于分发任务
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) run() {
	idx := 0
	/**
	run 方法内是一个无限循环，循环体中使用 select 监视 Pool 类型实例的两个 channel：quit 和 active。这种在 for 中使用 select 监视多个 channel 的实现，在 Go 代码中十分常见，是一种惯用法。
	*/
	for {
		select {
		/**
		当接收到来自 quit channel 的退出“信号”时，这个 Goroutine 就会结束运行。
		*/
		case <-p.quit:
			return
		/**
		而当 active channel 可写时，run 方法就会创建一个新的 worker Goroutine。
		此外，为了方便在程序中区分各个 worker 输出的日志，我这里将一个从 1 开始的变量 idx 作为 worker 的编号，并把它以参数的形式传给创建 worker 的方法。
		*/
		case p.active <- struct{}{}:
			// create a new worker
			idx++
			p.newWorker(idx)
		}
	}
}

/*
*
这是 Pool 类型的一个导出方法，workerpool 包的用户通过该方法向 pool 池提交待执行的任务（Task）。
Schedule 方法的核心逻辑，是将传入的 Task 实例发送到 workerpool 的 tasks channel 中。但考虑到现在 workerpool 已经被销毁的状态，我们这里通过一个 select，检视 quit channel 是否有“信号”可读，如果有，就返回一个哨兵错误 ErrWorkerPoolFreed。如果没有，一旦
p.tasks 可写，提交的 Task 就会被写入 tasks channel，以供 pool 中的 worker 处理。
这里要注意的是，这里的 Pool 结构体中的 tasks 是一个无缓冲的 channel，如果 pool 中 worker 数量已达上限，而且 worker 都在处理 task 的状态，那么 Schedule 方法就会阻塞，直到有 worker 变为 idle 状态来读取 tasks channel，schedule 的调用阻塞才会解除。
*/
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

/*
*
用于销毁一个 pool 池，停掉所有 pool 池中的 worker；
*/
func (p *Pool) Free() {
	close(p.quit) // make sure all worker and p.run exit and schedule return error
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}
