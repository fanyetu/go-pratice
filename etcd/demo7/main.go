package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config            clientv3.Config
		client            *clientv3.Client
		err               error
		kv                clientv3.KV
		getResp           *clientv3.GetResponse
		watchStartVersion int64
		watcher           clientv3.Watcher
		ctx               context.Context
		cancelFun         context.CancelFunc
		watchChan         <-chan clientv3.WatchResponse
		watchResp         clientv3.WatchResponse
		event             *clientv3.Event
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

	// 启动协程，模拟数据的修改删除
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "hello jobs")

			kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	// 获取一下当前的值，并从当前值之后进行监听
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// 从当前reversion之后的一个reversion进行监听
	watchStartVersion = getResp.Header.Revision + 1

	// 创建watcher
	watcher = clientv3.NewWatcher(client)

	// 创建一个可以取消的context
	ctx, cancelFun = context.WithCancel(context.TODO())
	go func() {
		// 5秒后，关闭context，关闭监听
		time.AfterFunc(5*time.Second, func() {
			cancelFun()
		})
	}()

	// 监听key数据变化
	watchChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartVersion))

	// 循环获取chan中的数据
	for watchResp = range watchChan {
		// 循环判断event的类型
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改了：", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了：", event.Kv.ModRevision)
			}
		}
	}
}
