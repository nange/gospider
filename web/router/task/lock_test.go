package task

import (
	"fmt"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	for i := 0; i < 3; i++ {
		taskID := uint64(1)
		go op1(taskID)
	}
	time.Sleep(5 * time.Second)
}

func op1(taskID uint64) {
	if taskLock.IsRunning(taskID) {
		fmt.Println("task is running...")
		return
	}
	defer taskLock.Complete(taskID)
	fmt.Println("run task")
	time.Sleep(1 * time.Second)
}
