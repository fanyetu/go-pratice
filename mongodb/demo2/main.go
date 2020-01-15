package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func main() {
	var (
		ctx      context.Context
		client   *mongo.Client
		err      error
		database *mongo.Database
		coll     *mongo.Collection
		log      *LogRecord
		result   *mongo.InsertOneResult
	)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	coll = database.Collection("log")

	// 插入一条数据
	log = &LogRecord{
		JobName: "test1",
		Command: "echo hello;",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	if result, err = coll.InsertOne(context.TODO(), log); err != nil {
		fmt.Println(err)
		return
	}

	// 返回的ID默认是ObjectId
	id := result.InsertedID.(primitive.ObjectID)
	fmt.Println("插入数据ID：", id.Hex())
}
