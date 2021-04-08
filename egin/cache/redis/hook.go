package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	goredis "github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/daodao97/egin/egin/utils/logger"
)

var namespace = "service"
var labels = []string{"endpoint", "method"}
var reqCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "redis_request_count_total",
		Help:      "Total number of Redis key call.",
	}, labels,
)

func init() {
	prometheus.MustRegister(reqCount)
}

type loggerHook struct {
}

func (l *loggerHook) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	spew.Dump()
	newCtx := context.WithValue(ctx, "start", time.Now())
	return newCtx, nil
}

func (l *loggerHook) AfterProcess(ctx context.Context, cmd goredis.Cmder) error {
	key := fmt.Sprintf("redis.%s.%s", cmd.Args()[0], cmd.Args()[1])
	tc := time.Since(ctx.Value("start").(time.Time))
	logger.NewLogger("redis").Info(key, map[string]interface{}{
		"method": cmd.FullName(),
		"args":   cmd.Args(),
		"ums":    fmt.Sprintf("%v", tc),
	})

	endpoint := key
	method := cmd.FullName()

	lvs := []string{endpoint, method}

	// api 请求计数
	reqCount.WithLabelValues(lvs...).Inc()
	return nil
}

func (l *loggerHook) BeforeProcessPipeline(ctx context.Context, cmds []goredis.Cmder) (context.Context, error) {
	fmt.Println("logger before pipe")
	return ctx, nil
}

func (l *loggerHook) AfterProcessPipeline(ctx context.Context, cmds []goredis.Cmder) error {
	fmt.Println("logger after pipe")
	return nil
}
