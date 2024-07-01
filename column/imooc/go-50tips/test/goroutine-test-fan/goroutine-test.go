package goroutine_test_fan

import (
	"fmt"
	"sync"
	"time"
)

func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}

func spawnGroup(name string, num int, f func(int) (int, bool), in <-chan int) <-chan int {
	groupOut := make(chan int)
	var outSlice []chan int
	for i := 0; i < num; i++ {
		out := make(chan int)
		go func(i int) {
			name := fmt.Sprintf("%s-%d:", name, i)
			fmt.Printf("%s begin to work...\n", name)

			for v := range in {
				r, ok := f(v)
				if ok {
					out <- r
				}
			}
			close(out)
			fmt.Printf("%s work done\n", name)
		}(i)
		outSlice = append(outSlice, out)
	}

	// Fan-in
	//
	// out --\
	//        \
	// out ---- --> groupOut
	//        /
	// out --/
	//
	go func() {
		var wg sync.WaitGroup
		for _, out := range outSlice {
			wg.Add(1)
			go func(out <-chan int) {
				for v := range out {
					groupOut <- v
				}
				wg.Done()
			}(out)
		}
		wg.Wait()
		close(groupOut)
	}()

	return groupOut
}

/*
*
通过spawnGroup函数实现了"fan-out"模式，针对每个输入 channel，我们都建立多个功能相同的 goroutine 从这个共同的输入 channel 读取数据并处理，直至 channel 被关闭。
在 spawnGroup 函数的结尾处，我们将多个 goroutine 的 out channel 聚合到一个 groupOut channel 中，这就是"fan-in"模式的实现。
*/
func TestFanIn() {
	in := newNumGenerator(1, 20)
	out := spawnGroup("square", 2, square, spawnGroup("filterOdd", 3, filterOdd, in))

	time.Sleep(3 * time.Second)

	for v := range out {
		fmt.Println(v)
	}
}
