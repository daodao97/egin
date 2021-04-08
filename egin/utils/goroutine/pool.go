package goroutine

import (
	"sync"
)

// Pool Goroutine Pool
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// New 新建一个协程池
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add 新增一个执行
func (p *Pool) Add(delta int) {
	// delta为正数就添加
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	// delta为负数就减少
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done 执行完成减一
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait 等待Goroutine执行完毕
func (p *Pool) Wait() {
	p.wg.Wait()
}

// example
// 这里限制5个并发
// pool := NewPool(5)
// fmt.Println("the NumGoroutine begin is:", runtime.NumGoroutine())
// for i := 0; i < 20; i++ {
// 	pool.Add(1)
// 	go func(i int) {
// 		time.Sleep(time.Second)
// 		fmt.Println("the NumGoroutine continue is:", runtime.NumGoroutine())
// 		pool.Done()
// 	}(i)
// }
// pool.Wait()
// fmt.Println("the NumGoroutine done is:", runtime.NumGoroutine())
