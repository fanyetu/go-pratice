package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

// 过滤查询的结构体
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		ctx      context.Context
		client   *mongo.Client
		err      error
		database *mongo.Database
		coll     *mongo.Collection
		conn     *FindByJobName
		findOpt  *options.FindOptions
		cursor   *mongo.Cursor
		result   *LogRecord
	)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	coll = database.Collection("log")

	// 创建查询条件
	conn = &FindByJobName{JobName: "test1"}

	// 创建查询配置
	findOpt = options.Find()
	findOpt.SetSkip(0)  // 从第0个开始
	findOpt.SetLimit(2) // 查询2个

	if cursor, err = coll.Find(context.TODO(), conn, findOpt); err != nil {
		fmt.Println(err)
		return
	}

	// 关闭游标
	defer cursor.Close(context.TODO())

	// 遍历游标获取数据
	for cursor.Next(context.TODO()) {
		// 创建接收对象
		result = &LogRecord{}

		// 将数据反序列化
		if err = cursor.Decode(result); err != nil {
			fmt.Println(err)
			return
		}

		// 输出结果
		fmt.Println(*result)
	}
}
