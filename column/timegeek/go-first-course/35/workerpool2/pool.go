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
	capacity int  // workerpool大小
	preAlloc bool // 是否在创建pool的时候，就预创建workers，默认值为：false

	// 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true
	// 如果block = false，则Schedule返回ErrNoWorkerAvailInPool
	block  bool
	active chan struct{} // 对应上图中的active channel

	tasks chan Task // 对应上图中的task channel

	wg   sync.WaitGroup // 用于在pool销毁时等待所有worker退出
	quit chan struct{}  // 用于通知各个worker退出的信号channel
}

type Task func()

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

/*
*
新版 New 函数除了接受 capacity 参数之外，还在它的参数列表中增加了一个类型为 Option 的可变长参数 opts。在 New 函数体中，我们通过一个 for 循环，将传入的 Option 运用到 Pool 类型的实例上。
新版 New 函数还会根据 preAlloc 的值来判断是否预创建所有的 worker，如果需要，就调用 newWorker 方法把所有 worker 都创建出来。newWorker 的实现与上一版代码并没有什么差异，这里就不再详说了。
*/
func New(capacity int, opts ...Option) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		block:    true,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	for _, opt := range opts {
		opt(p)
	}

	fmt.Printf("workerpool start(preAlloc=%t)\n", p.preAlloc)

	if p.preAlloc {
		// create all goroutines and send into works channel
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i + 1)
			p.active <- struct{}{}
		}
	}

	go p.run()

	return p
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", i)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			/**
			指定的WithPreAllocWorkers(false)，即不预创建。run阻塞在对tasks channel的读取上，直到后面main goroutine第一次调用Schedule写入一个task，这时run才会开始按需创建Worker。
			*/
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) returnTask(t Task) {
	go func() {
		p.tasks <- t
	}()
}

/*
*
由于 preAlloc 选项的加入，Pool 的 run 方法的实现有了变化
新版 run 方法在 preAlloc=false 时，会根据 tasks channel 的情况在适合的时候创建 worker （第 4 行~ 第 18 行)，直到 active channel 写满，才会进入到和第一版代码一样的调度逻辑中（第 20 行~ 第 29 行）。
*/
func (p *Pool) run() {
	idx := len(p.active)

	if !p.preAlloc {
		/**
		loop 标签用于标记 for 循环，并允许在循环内部使用 break loop 语句跳出该循环。这个结构有助于在特定条件下优雅地退出循环，而不仅仅是从 select 语句中跳出。
		*/
	loop:
		for t := range p.tasks {
			p.returnTask(t)
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				idx++
				p.newWorker(idx)
			default:
				break loop
			}
		}
	}

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			// create a new worker
			idx++
			p.newWorker(idx)
		}
	}
}

/*
*
Schedule 函数也因 WithBlock 选项，有了一些变化：
Schedule 在 tasks channel 无法写入的情况下，进入 default 分支。在 default 分支中， Schedule 根据 block 字段的值，决定究竟是继续阻塞在 tasks channel 上，还是返回 ErrNoIdleWorkerInPool 错误。
*/
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		/**
		问题：既然没有达到最大 worker 数，为什么不是去创建新的 worker 而是直接返回错误呢？这一点不是很理解，不应该是根据 task 自动创建预期内的 worker 直到 worker 数满了再返回没有空闲的 worker 错误吗？
		回答：首先这仅是一个demo。goroutine池有很多种实现方式，甚至基于同一种方式，比如文中的channel也有不同的策略设定。也许文中的demo设定的当p.tasks阻塞就返回错误（当p.block=false时）有些不合理。
		"让调度时不阻塞，是不是把 task 放到一个队列里排队更合理，再加上一个丢弃策略，类似 Java 中的线程池" -- 很好的提议。这也是一种设计过程的设定，在我的demo的基础上，稍作修改应该就能实现。
		*/
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoIdleWorkerInPool
	}
}

func (p *Pool) Free() {
	close(p.quit) // make sure all worker and p.run exit and schedule return error
	p.wg.Wait()
	fmt.Printf("workerpool freed(preAlloc=%t)\n", p.preAlloc)
}
