package main

import (
	"fmt"
	"time"

	"github.com/bigwhite/workerpool"
)

/*
*
task[2]: error: no idle worker in pool task[4]: error: no idle worker in pool task[5]: error: no idle worker in pool task[6]: error: no idle worker in pool task[7]: error: no idle worker in pool task[8]: error: no idle worker in pool task[9]: error: no idle worker in pool worker[001]: receive a task
worker[002]: start
worker[002]: exit
worker[001]: receive a task
worker[001]: exit
workerpool freed(preAlloc=false)
由于 Goroutine 调度的不确定性，这个结果仅仅是很多种结果的一种。我们看到，仅仅 001 这个 worker 收到了 task，其余的 worker 都因为 worker 尚未创建完毕，而返回了错误，而不是像 demo1 那样阻塞在 Schedule 调用上。
*/
func main() {
	p := workerpool.New(5, workerpool.WithPreAllocWorkers(false), workerpool.WithBlock(false))

	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			fmt.Printf("task[%d]: error: %s\n", i, err.Error())
		}
	}

	p.Free()
}
