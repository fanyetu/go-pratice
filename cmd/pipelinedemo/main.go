package main

import (
	"fmt"
	"go-pratice/pipeline"
)

func main() {
	p := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(13, 3, 7, 6, 15, 2, 33, 4)),
		pipeline.InMemSort(pipeline.ArraySource(2, 5, 6, 8, 1, 21)),
	)

	// 第一种循环方式
	//for {
	//	if num, ok := <-p; ok {
	//		fmt.Print(num)
	//	} else {
	//		break
	//	}
	//}

	// 第二种循环方式
	for num := range p {
		fmt.Println(num)
	}
}
