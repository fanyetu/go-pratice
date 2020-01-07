package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Result struct {
	data []byte
	err  error
}

func main() {
	var (
		cmd       *exec.Cmd
		ctx       context.Context
		cancelFun context.CancelFunc
		results   chan *Result
		res       *Result
	)

	results = make(chan *Result, 1000)

	// 创建一个可以关闭的context
	ctx, cancelFun = context.WithCancel(context.TODO())

	// 开启协程执行命令
	go func() {
		var (
			data []byte
			err  error
		)
		// 休眠2秒后输出hello
		cmd = exec.CommandContext(ctx, "C:\\cygwin64\\bin\\bash.exe", "-c", "sleep 2;echo hello;")

		// 执行命令并输出
		data, err = cmd.CombinedOutput()

		// 向队列中输出结果，主协程接收
		results <- &Result{
			data: data,
			err:  err,
		}
	}()

	time.Sleep(1 * time.Second)

	// 执行cancelFun，将命令杀死
	cancelFun()

	res = <-results

	fmt.Println(res.err, string(res.data))
}
