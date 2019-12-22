package pipeline

import "sort"

// 归并排序的核心，归并操作
func Merge(in1, in2 <-chan int) chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		// 当一个channel有数据就可以循环
		for ok1 || ok2 {
			// 当v2不存在，或者v1小于等于v2的时候就可以输出到out中
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				// 输出完毕后，更新v1,v2的值
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}

		close(out)
	}()
	return out
}

// 在内存中排序
func InMemSort(in <-chan int) chan int {
	out := make(chan int)
	go func() {
		// 读取到内存中
		a := []int{}
		for v := range in {
			a = append(a, v)
		}

		// 排序
		sort.Ints(a)

		// 输出
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// 传入可变长数组，返回channel
func ArraySource(a ...int) chan int {
	out := make(chan int)
	go func() {
		// 循环读取数组，将其放入channel中
		for _, v := range a {
			out <- v
		}
		// 关闭channel，表示数据写入完毕
		close(out)
	}()

	return out
}
