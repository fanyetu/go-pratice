package iotest

import "io"

// 包装另一个流
type alphaReader2 struct {
	reader io.Reader
}

func newAlphaReader2(reader io.Reader) *alphaReader2 {
	return &alphaReader2{reader: reader}
}

func (a *alphaReader2) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)
	if err != nil {
		return n, err
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}

	copy(p, buf)
	return n, nil
}

// 实现从流中过滤非字母的字符
type alphaReader struct {
	// 原始数据
	src string
	// 当前读取的位置
	cur int
}

func newAlphaReader(src string) *alphaReader {
	return &alphaReader{src: src}
}

// 过滤函数
func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

// 实现Reader接口的Read函数
func (a *alphaReader) Read(p []byte) (int, error) {
	// 如果当前的位置超过了源数据的长度，那么就是读取完了，返回EOF
	if a.cur >= len(a.src) {
		return 0, io.EOF
	}

	// 剩余未读的长度
	x := len(a.src) - a.cur
	n, bound := 0, 0
	if x >= len(p) {
		// 剩余长度大于缓冲区大小，所以本次可以填满缓冲区
		bound = len(p)
	} else if x < len(p) {
		// 剩余长度不够填满缓冲区
		bound = x
	}

	buf := make([]byte, bound)
	for n < bound {
		// 每次读取一个字节，执行过滤函数
		if char := alpha(a.src[a.cur]); char != 0 {
			buf[n] = char
		}
		n++
		a.cur++
	}

	// 将处理后的buf复制给p
	copy(p, buf)
	return n, nil
}
