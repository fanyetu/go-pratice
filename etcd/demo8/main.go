package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

// etcd client的op操作
func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		op     clientv3.Op
		opResp clientv3.OpResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 创建连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	// 创建put类型的op
	op = clientv3.OpPut("/cron/jobs/job8", "hello")

	// 执行op
	if opResp, err = kv.Do(context.TODO(), op); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("修改了:", opResp.Put().Header.Revision)

	// 创建get类型的op
	op = clientv3.OpGet("/cron/jobs/job8")

	// 执行op
	if opResp, err = kv.Do(context.TODO(), op); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("获取：", opResp.Get().Header.Revision, string(opResp.Get().Kvs[0].Value))

}
