package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr *cronexpr.Expression
		err  error
		now  time.Time
		next time.Time
	)

	// 解析cron表达式
	if expr, err = cronexpr.Parse("*/7 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	// 当前时间
	now = time.Now()

	// 通过expr计算处下一次执行时间
	next = expr.Next(now)

	fmt.Println(now, next)

	// 间隔时间后执行方法
	time.AfterFunc(next.Sub(now), func() {
		fmt.Println("执行方法:", time.Now())
	})

	time.Sleep(15 * time.Second)
}
