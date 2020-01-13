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
		resp   *clientv3.DeleteResponse
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

	// 删除一个kv
	if resp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		// 打印被删除之前的value
		if resp.Deleted > 0 {
			for _, kvpar := range resp.PrevKvs {
				fmt.Println(string(kvpar.Key) + ":" + string(kvpar.Value))
			}
		}
	}
}
