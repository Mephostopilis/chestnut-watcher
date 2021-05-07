// 题目: 写出下面代码输出内容。
// output:
// 打印前
// 打印中
// 打印后
// panic: 触发异常
package main

import (
	"fmt"
)

func defer_call() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}

func main() {
	defer_call()
}
