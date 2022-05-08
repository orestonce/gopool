# go线程池使用方法

````go
package main

import (
	"fmt"
	"github.com/orestonce/gopool"
)

func main() {
	task := gopool.NewThreadPool(10) // 最大线程数
	task.AddJob(func() {
		fmt.Println("thread1")
	})
	task.AddJob(func() {
		fmt.Println("thread2")
	})
	task.CloseAndWait()
}
````