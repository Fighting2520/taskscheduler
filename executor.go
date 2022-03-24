package taskscheduler

type (
	Executor interface {
		Do(TaskFunc) error
	}

	SimpleExecutor struct {
	}
)

func NewSimpleExecutor() *SimpleExecutor {
	return &SimpleExecutor{}
}

func (se *SimpleExecutor) Do(taskFunc TaskFunc) error {
	return taskFunc()
}
