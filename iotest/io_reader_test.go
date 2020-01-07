package iotest

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestAlphaReader2_Read(t *testing.T) {
	reader := newAlphaReader2(strings.NewReader("Hello! It's 9am, where is the sun?"))
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

func TestAlphaReader_Read(t *testing.T) {
	reader := newAlphaReader("Hello! It's 9am, where is the sun?")
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

// 通过流的方式读取字符串
func TestStringReader(t *testing.T) {
	reader := strings.NewReader("Elasticsearch is better than Mongodb.")
	bytes := make([]byte, 4)

	for {
		n, err := reader.Read(bytes)
		if err != nil {
			// 出现EOF错误就是当前数据读取完了
			if err == io.EOF {
				fmt.Println("EOF:", n)
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(n, string(bytes[:n]))
	}
}
