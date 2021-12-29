package workerpool

type (
	WorkerPool interface {
		Submit(task Task) Submit
		Name() string
		Workers() int
	}

	Task interface {
		Execute() TaskResult
		Name() string
	}

	TaskFn func() (interface{}, error)

	TaskResult interface {
		Name() string
		Value() interface{}
		Err() error
	}

	Submit interface {
		Result() TaskResult
	}
)

func New(name string, workers int, queue int) WorkerPool {
	return nil
}

func NewTask(name string, fn TaskFn) Task {
	return nil
}