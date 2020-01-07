package iotest

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestWriter(t *testing.T) {
	proverbs := []string{
		"Channels orchestrate mutexes serialize",
		"Cgo is not Go",
		"Errors are values",
		"Don't panic",
	}
	var writer bytes.Buffer
	for _, proverb := range proverbs {
		n, err := writer.Write([]byte(proverb))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(proverb) {
			fmt.Println("failed to write data")
			os.Exit(1)
		}
	}

	fmt.Println(writer.String())
}
