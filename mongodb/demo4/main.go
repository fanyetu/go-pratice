package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// bson格式{"$lt": timestamp}
type BeforeTimeCond struct {
	Before int64 `bson:"$lt"`
}

// bson格式{"timePoint.startTime": {"$lt“: timestamp}}
type DeleteCond struct {
	BeforeTime BeforeTimeCond `bson:"timePoint.startTime"`
}

// 删除操作
func main() {
	var (
		ctx      context.Context
		client   *mongo.Client
		err      error
		database *mongo.Database
		coll     *mongo.Collection
		cond     *DeleteCond
		result   *mongo.DeleteResult
	)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	coll = database.Collection("log")

	// 删除当前之间之前的记录
	cond = &DeleteCond{BeforeTime: BeforeTimeCond{Before: time.Now().Unix()}}

	if result, err = coll.DeleteMany(context.TODO(), cond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除记录数：", result.DeletedCount)
}
