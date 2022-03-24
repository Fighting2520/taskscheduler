package taskscheduler

import (
	"context"
	"time"
)

type OptionFunc func(*Scheduler)

func SetTimeout(dur time.Duration) OptionFunc {
	return func(scheduler *Scheduler) {
		ctx, cancel := context.WithTimeout(scheduler.ctx, dur)
		scheduler.ctx = ctx
		scheduler.cancel = cancel
	}
}

func SetTaskConcurrency(c int) OptionFunc {
	return func(scheduler *Scheduler) {
		scheduler.workPool = make(chan struct{}, c)
	}
}

func SetRecoverFn(fn RecoverFunc) OptionFunc {
	return func(scheduler *Scheduler) {
		scheduler.recoverFn = fn
	}
}

func SetTaskSize(size int) OptionFunc {
	return func(scheduler *Scheduler) {
		scheduler.taskPool = make(chan TaskFunc, size)
	}
}

func SetExecutor(executor Executor) OptionFunc {
	return func(scheduler *Scheduler) {
		if executor != nil {
			scheduler.executor = executor
		}
	}
}

func SetLogger(logger Logger) OptionFunc {
	return func(scheduler *Scheduler) {
		if logger != nil {
			scheduler.log = logger
		}
	}
}
