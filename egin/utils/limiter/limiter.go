package limiter

import (
	"container/ring"
	"fmt"
	"time"

	"github.com/daodao97/egin/utils/goroutine"
)

// 基于滑动时间窗口的限流器
type Limiter struct {
	LimitCount  int        // 限制数量
	limitBucket int        // 滑动窗口数
	head        *ring.Ring // 环形队列（链表)
	initEd      bool
}

func (l *Limiter) init() {
	if l.initEd {
		return
	}
	l.limitBucket = 10
	head := ring.New(l.limitBucket)
	for i := 0; i < l.limitBucket; i++ {
		head.Value = 0
		head = head.Next()
	}
	l.head = head
	l.initEd = true
	goroutine.Go(func() {
		// ms级别，limitBucket int = 10意味将每秒分为10份，每份100ms
		for range time.Tick(time.Millisecond * time.Duration(1000/l.limitBucket)) {
			l.head.Value = 0
			l.head = l.head.Next()
		}
	})
}

func (l *Limiter) CheckOverLimit() bool {
	l.init()
	total := 0
	l.head.Do(func(i interface{}) {
		total += i.(int)
	})
	fmt.Println("total", total, "limit", l.LimitCount, total > l.LimitCount)
	return total > l.LimitCount
}

func (l *Limiter) Incr() {
	l.init()
	l.head.Value = l.head.Value.(int) + 1
}
