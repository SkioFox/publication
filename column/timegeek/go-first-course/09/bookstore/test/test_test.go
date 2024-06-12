package test

import "testing"

func sum(max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += i
	}
	return total
}
func fooWithDefer() {
	defer func() {
		sum(10)
	}()
}
func fooWithoutDefer() {
	sum(10)
}
func BenchmarkFooWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithDefer()
	}
}
func BenchmarkFooWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithoutDefer()
	}
}

/**
执行
	go mod目录执行
 	go test -bench . ./test/test_test.go
	结果 => go version go1.21.5 darwin/amd64
		goos: darwin
		goarch: amd64
		cpu: Intel(R) Core(TM) i5-8500B CPU @ 3.00GHz
		BenchmarkFooWithDefer-6         206318150                5.834 ns/op
		BenchmarkFooWithoutDefer-6      292855430                4.114 ns/op
		PASS
		ok      command-line-arguments  3.658s
	从 Go 1.13 版本开始，Go 核心团队对 defer 性能进行了多次优化，到现在的 Go 1.17 版本，defer 的开销已经足够小了。我们看看使用 Go 1.17 版本运行上述基准测试的结果：
*/
