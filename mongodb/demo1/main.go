package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		ctx        context.Context
		database   *mongo.Database
		collection *mongo.Collection
	)
	// 创建连接
	if client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Println(err)
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		fmt.Println(err)
		return
	}

	// 选择数据库
	database = client.Database("my_db")

	// 选择collection
	collection = database.Collection("my_collection")

	collection = collection
}
