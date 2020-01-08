package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		resp   *clientv3.PutResponse
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

	// 获取client中的KV对象，对etcd进行操作
	kv = clientv3.NewKV(client)

	// 使用kv进行put操作，使用WithPreKV，可以查看之前的历史
	if resp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}

	// 打印当前操作的revision，revision全局递增
	fmt.Println(resp.Header.Revision)

	if resp.PrevKv != nil {
		// 打印当前key的上一个记录
		fmt.Println(string(resp.PrevKv.Value))
	}

}
