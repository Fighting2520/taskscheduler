package taskscheduler

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduler_Start(t *testing.T) {
	s := NewScheduler(SetTimeout(5*time.Second), SetTaskConcurrency(5), SetTaskSize(1))
	s.Start()
	for i := 0; i < 20; i++ {
		if err := s.AddTask(func() error {
			time.Sleep(1 * time.Second)
			fmt.Println("hello world")
			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	fmt.Println(12345)
	s.Finish()
	s.Wait()
	fmt.Println(s.ExecutedCount(), s.UnExecutedCount())
	fmt.Println(s.Error())
}
