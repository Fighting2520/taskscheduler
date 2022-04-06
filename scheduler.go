package taskscheduler

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
)

type (
	Scheduler struct {
		executor  Executor
		ctx       context.Context
		cancel    context.CancelFunc
		ec        int // executed task count 已执行的任务数量
		uc        int // unexecuted task count 未执行的任务数
		err       error
		recoverFn RecoverFunc
		closed    bool
		done      chan struct{}

		taskMu   sync.RWMutex
		taskPool chan TaskFunc
		workPool chan struct{}
		log      Logger
	}
)

const (
	defaultConcurrency = 5
	defaultTaskSize    = 100
)

var ErrSchedulerClosed = errors.New("scheduler closed")

func NewScheduler(opts ...OptionFunc) *Scheduler {
	ctx, cancel := context.WithCancel(context.TODO())
	var s = Scheduler{
		executor:  NewSimpleExecutor(),
		ctx:       ctx,
		cancel:    cancel,
		recoverFn: defaultRecoverFunc,
		taskPool:  make(chan TaskFunc, defaultTaskSize),
		done:      make(chan struct{}),
		workPool:  make(chan struct{}, defaultConcurrency),
		log:       &simpleLog{},
	}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

type RecoverFunc func()

func defaultRecoverFunc() {
	if err := recover(); err != nil {
		logger := simpleLog{}
		logger.Errorf(" panic: %v", err)
		logger.Errorf(string(debug.Stack()))
	}
}

func (s *Scheduler) Start() {
	s.log.Infof("scheduler 已经启动, 当前设置任务池总数为: %d, 工作线程数为: %d\n", cap(s.taskPool), cap(s.workPool))
	var wg sync.WaitGroup
	defer s.recoverFn()
	defer func() {
		wg.Wait()
		close(s.workPool)
		s.log.Infof("Scheduler stopped\n")
		s.done <- struct{}{}
	}()
	for {
		select {
		case s.workPool <- struct{}{}:
			task, ok := <-s.taskPool
			if !ok {
				<-s.workPool
				return
			}
			wg.Add(1)
			go func() {
				defer s.recoverFn()
				defer func() {
					s.ec++
					wg.Done()
					<-s.workPool
				}()
				err := s.executor.Do(task)
				if err != nil {
					s.log.Errorf("work failed with: %s\n", err)
				}
			}()
		case <-s.ctx.Done():
			if s.cancel != nil {
				s.cancel()
			}
			s.err = s.ctx.Err()
			s.clearTasks()
			s.log.Errorf("scheduler exit with: %s\n", s.ctx.Err())
			return
		}
	}
}

func (s *Scheduler) clearTasks() {
	for range s.taskPool {
		s.uc++
	}
	s.log.Infof("未执行任务数为: %d", s.uc)
}

func (s *Scheduler) Finish() {
	s.taskMu.Lock()
	if !s.closed {
		close(s.taskPool)
		s.closed = true
	}
	s.taskMu.Unlock()
}

func (s *Scheduler) Closed() bool {
	s.taskMu.RLock()
	defer s.taskMu.RUnlock()
	return s.closed
}

func (s *Scheduler) AddTask(fn TaskFunc) error {
	if s.Closed() {
		return ErrSchedulerClosed
	}
	s.taskPool <- fn
	return nil
}

func (s *Scheduler) ExecutedCount() int {
	return s.ec
}

func (s *Scheduler) UnExecutedCount() int {
	return s.uc
}

func (s *Scheduler) Error() string {
	if s.err == nil {
		return ""
	}
	return s.err.Error()
}

func (s *Scheduler) Wait() {
	<-s.done
}
