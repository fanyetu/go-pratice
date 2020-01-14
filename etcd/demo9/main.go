package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

// 使用租约、事务、op实现分布式锁
func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		kv             clientv3.KV
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		ctx            context.Context
		cancelFun      context.CancelFunc
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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

	// 1.创建锁(创建租约，自动续租，抢锁)

	lease = clientv3.NewLease(client)

	// 创建租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// 获取到租约的id
	leaseId = leaseGrantResp.ID

	// 设置租约可取消的自动续租
	ctx, cancelFun = context.WithCancel(context.TODO())
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 程序结束的时候，关闭自动续租并释放租约
	defer cancelFun()
	defer lease.Revoke(context.TODO(), leaseId)

	// 开启协程判断续租状态
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto END
				} else {
					fmt.Println("续租成功:", keepResp.ID)
				}
			}
		}
	END:
	}()

	// 抢占key
	// 如果key的createRevision为0，那么则是还没创建，怎设置key，抢锁成功，如果不是则抢锁失败
	txn = kv.Txn(context.TODO()).If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9"))

	// 提交事务，并判断事务的执行情况
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	// 没有成功就是进入了Else，抢锁失败
	if !txnResp.Succeeded {
		fmt.Println("抢锁失败:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2.执行任务

	fmt.Println("执行任务")
	time.Sleep(10 * time.Second)

	// 3.释放锁(关闭自动续租，释放租约)
	// 在defer中操作
}
