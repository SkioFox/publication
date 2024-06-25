package workerpool

type Option func(*Pool)

/*
*
Schedule 调用是否阻塞
*/
func WithBlock(block bool) Option {
	return func(p *Pool) {
		p.block = block
	}
}

/*
*
是否预创建所有的 worker。
*/
func WithPreAllocWorkers(preAlloc bool) Option {
	return func(p *Pool) {
		p.preAlloc = preAlloc
	}
}
