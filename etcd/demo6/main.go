package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

// 实现etcd的自动过期
func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		kv           clientv3.KV
		lease        clientv3.Lease
		grantResp    *clientv3.LeaseGrantResponse
		leaseId      clientv3.LeaseID
		putResp      *clientv3.PutResponse
		getResp      *clientv3.GetResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
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

	// 创建租约对象
	lease = clientv3.NewLease(client)

	// 创建一个10秒自动过期的租约
	if grantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 获取租约的id
	leaseId = grantResp.ID

	// 使用keepAlive进行自动续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 在协程中判断是否续租成功
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				// 如果chan中获取的resp为空，那么说明续租失败
				if keepResp == nil {
					fmt.Println("续租失败")
					return
				} else {
					fmt.Println("续租成功：", keepResp.ID)
				}
			}
		}
	}()

	kv = clientv3.NewKV(client)

	// put一个KV，并且和租约进行关联，实现自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "value", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("插入成功：", putResp.Header.Revision)

	// 定时获取kv，查看是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count > 0 {
			fmt.Println("获取成功：", getResp.Kvs)
		} else {
			fmt.Println("kv过期")
			break
		}
		time.Sleep(2 * time.Second)
	}
}
