package gopool_test

import (
	"fmt"
	"github.com/orestonce/gopool"
	"sync"
	"testing"
)

func TestNewThreadPool(t *testing.T) {
	task := gopool.NewThreadPool(5)
	data := make([]int, 10)

	runningCount := 0
	var runningCountLocker sync.Mutex
	maxRunningCount := 0

	ch := make(chan int)

	for i := 0; i < 10; i++ {
		id := i
		task.AddJob(func() {
			runningCountLocker.Lock()
			runningCount++
			if maxRunningCount < runningCount {
				maxRunningCount = runningCount
				if maxRunningCount >= 5 { // 大于等于 5 时候关闭ch, 如果两次都大于等于5 close会让他panic
					close(ch)
				}
			}
			runningCountLocker.Unlock()

			<-ch
			fmt.Println("Hello world.", id)
			data[id] = id

			runningCountLocker.Lock()
			runningCount--
			runningCountLocker.Unlock()
		})
	}
	task.CloseAndWait()
	for i := 0; i < 10; i++ {
		if data[i] != i {
			t.Fatal(i, data[i])
		}
	}
	task.CloseAndWait()

	if runningCount != 0 || maxRunningCount != 5 {
		t.Fatal(runningCount, maxRunningCount)
	}
}
