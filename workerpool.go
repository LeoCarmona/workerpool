package workerpool

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	name string

	minWorkers   int
	maxWorkers   int
	idleTimeout  time.Duration
	upScaling    int
	downScaling  int
	queue        chan func()
	panicHandler func(panicErr interface{})

	metrics *Metrics

	runningWorkers  int32
	idleWorkers     int32
	submittedTasks  uint64
	waitingTasks    uint64
	successfulTasks uint64
	failedTasks     uint64

	// properties
	stopOnce sync.Once
	workerWG sync.WaitGroup
	stopped  bool

	// events
	idleTimeoutChanged chan struct{}

	// signals
	stoppedSignal chan struct{}
	purgeSignal   chan struct{}
}

func New(name string, settings *Settings) *WorkerPool {
	wp := &WorkerPool{
		name:               name,
		minWorkers:         settings.MinWorkers,
		maxWorkers:         settings.MaxWorkers,
		idleTimeout:        settings.IdleTimeout,
		upScaling:          settings.UpScaling,
		downScaling:        settings.DownScaling,
		panicHandler:       settings.PanicHandler,
		stoppedSignal:      make(chan struct{}),
		purgeSignal:        make(chan struct{}),
		idleTimeoutChanged: make(chan struct{}, 1),
	}

	wp.metrics = &Metrics{wp: wp}

	if wp.maxWorkers < 1 {
		if wp.minWorkers > 0 {
			wp.maxWorkers = wp.minWorkers
		} else {
			wp.maxWorkers = runtime.NumCPU()
		}
	}

	if wp.minWorkers > wp.maxWorkers {
		wp.minWorkers = wp.maxWorkers
	} else if wp.minWorkers < 0 {
		wp.minWorkers = 0
	}

	if wp.idleTimeout < 1 {
		wp.idleTimeout = defaultIdleTimeout
	}

	if wp.upScaling < 1 {
		wp.upScaling = 1
	}

	if wp.downScaling < 1 {
		wp.downScaling = 1
	}

	if settings.Queue < 0 {
		wp.queue = make(chan func())
	} else {
		wp.queue = make(chan func(), settings.Queue)
	}

	if wp.panicHandler == nil {
		wp.panicHandler = defaultPanicHandler
	}

	wp.createWorkers(wp.minWorkers).Wait()

	go wp.monitor()

	return wp
}

// Submit sends a task to run asynchronous in the WorkerPool.
// If the queue is full, it will block until a task is dispatched to a worker.
func (wp *WorkerPool) Submit(task func()) {
	wp.submit(task, true, false, nil)
}

func (wp *WorkerPool) SubmitAndWait(task func()) {
	wp.submit(task, true, true, nil)
}

func (wp *WorkerPool) SubmitAndWaitWithTimeout(timeout time.Duration, task func()) (finished bool, wait <-chan struct{}) {
	return wp.SubmitAndWaitWithDeadline(time.Now().Add(timeout), task)
}

func (wp *WorkerPool) SubmitAndWaitWithDeadline(deadline time.Time, task func()) (finished bool, wait <-chan struct{}) {
	_, finished, wait = wp.submit(task, true, true, &deadline)
	return
}

// TrySubmit Attempts to send a task to run asynchronous in the WorkerPool.
// If the queue has capacity, the task runs asynchronously and returns true. Otherwise, the task will not run and will return false.
func (wp *WorkerPool) TrySubmit(task func()) (submitted bool) {
	submitted, _, _ = wp.submit(task, false, false, nil)
	return
}

func (wp *WorkerPool) TrySubmitAndWait(task func()) (submitted bool) {
	submitted, _, _ = wp.submit(task, false, true, nil)
	return
}

func (wp *WorkerPool) TrySubmitAndWaitWithTimeout(timeout time.Duration, task func()) (submitted bool, finished bool, wait <-chan struct{}) {
	return wp.TrySubmitAndWaitWithDeadline(time.Now().Add(timeout), task)
}

func (wp *WorkerPool) TrySubmitAndWaitWithDeadline(deadline time.Time, task func()) (submitted bool, finished bool, wait <-chan struct{}) {
	submitted, finished, wait = wp.submit(task, false, true, &deadline)
	return
}

func (wp *WorkerPool) Burst(workers int) {
	if workers < 1 {
		return
	}

	wp.createWorkers(workers)
}

func (wp *WorkerPool) ScaleUp(workers int) int {
	if workers < 1 {
		return 0
	}

	availableWorkersToBeCreated := wp.maxWorkers - wp.metrics.RunningWorkers()
	if availableWorkersToBeCreated < 1 {
		return 0
	}

	if workers > availableWorkersToBeCreated {
		workers = availableWorkersToBeCreated
	}

	wp.createWorkers(workers)
	return workers
}

func (wp *WorkerPool) ScaleDown(workers int) int {
	if workers < 1 {
		return 0
	}

	maxWorkersToBePurged := wp.metrics.RunningWorkers() - wp.minWorkers
	if maxWorkersToBePurged < 1 {
		return 0
	}

	if workers > maxWorkersToBePurged {
		workers = maxWorkersToBePurged
	}

	for i := 0; i < workers; i++ {
		wp.purgeSignal <- struct{}{}
	}

	return workers
}

func (wp *WorkerPool) ReleaseIdleWorkers() int {
	workersToBePurged := wp.metrics.IdleWorkers() - wp.minWorkers
	if workersToBePurged < 1 {
		return 0
	}

	for i := 0; i < workersToBePurged; i++ {
		wp.purgeSignal <- struct{}{}
	}

	return workersToBePurged
}

func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		wp.stopped = true
		close(wp.stoppedSignal)
	})
}

func (wp *WorkerPool) StopAndWait() {
	wp.Stop()
	wp.workerWG.Wait()
}

func (wp *WorkerPool) StopAndWaitWithTimeout(timeout time.Duration) (finished bool, wait <-chan struct{}) {
	return wp.StopAndWaitWithDeadline(time.Now().Add(timeout))
}

func (wp *WorkerPool) StopAndWaitWithDeadline(deadline time.Time) (finished bool, wait <-chan struct{}) {
	wp.Stop()

	done := make(chan struct{})
	go func() {
		wp.workerWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		// stopped before deadline
		return true, done
	case <-time.After(time.Until(deadline)):
		// deadline expired, then unlock
		return false, done
	}
}

func (wp *WorkerPool) Stopped() bool {
	return wp.stopped
}

func (wp *WorkerPool) Metrics() *Metrics {
	return wp.metrics
}

func (wp *WorkerPool) MinWorkers() int {
	return wp.minWorkers
}

func (wp *WorkerPool) SetMinWorkers(minWorkers int) int {
	if minWorkers < 0 {
		return 0
	}

	if minWorkers > wp.maxWorkers {
		wp.maxWorkers = minWorkers
	}

	workersToBeCreated := minWorkers - wp.minWorkers
	wp.minWorkers = minWorkers

	if workersToBeCreated < 1 {
		return 0
	}

	wp.createWorkers(workersToBeCreated)
	return workersToBeCreated
}

func (wp *WorkerPool) MaxWorkers() int {
	return wp.maxWorkers
}

func (wp *WorkerPool) SetMaxWorkers(maxWorkers int) int {
	if maxWorkers < 1 {
		return 0
	}

	if wp.minWorkers > maxWorkers {
		wp.minWorkers = maxWorkers
	}

	workersToBePurged := wp.metrics.RunningWorkers() - maxWorkers
	wp.maxWorkers = maxWorkers

	if workersToBePurged < 1 {
		return 0
	}

	for i := 0; i < workersToBePurged; i++ {
		wp.purgeSignal <- struct{}{}
	}

	return workersToBePurged
}

func (wp *WorkerPool) IdleTimeout() time.Duration {
	return wp.idleTimeout
}

func (wp *WorkerPool) SetIdleTimeout(idleTimeout time.Duration) {
	if idleTimeout < 1 {
		idleTimeout = defaultIdleTimeout
	}

	wp.idleTimeout = idleTimeout
	wp.idleTimeoutChanged <- struct{}{}
}

func (wp *WorkerPool) UpScaling() int {
	return wp.upScaling
}

func (wp *WorkerPool) SetUpScaling(upScaling int) {
	if upScaling < 1 {
		upScaling = 1
	}

	wp.upScaling = upScaling
}

func (wp *WorkerPool) DownScaling() int {
	return wp.downScaling
}

func (wp *WorkerPool) SetDownScaling(downScaling int) {
	if downScaling < 1 {
		downScaling = 1
	}

	wp.downScaling = downScaling
}

func (wp *WorkerPool) QueueCapacity() int {
	return cap(wp.queue)
}

func (wp *WorkerPool) PanicHandler() func(panicErr interface{}) {
	return wp.panicHandler
}

func (wp *WorkerPool) SetPanicHandler(panicHandler func(panicErr interface{})) {
	if panicHandler == nil {
		panicHandler = defaultPanicHandler
	}

	wp.panicHandler = panicHandler
}

func (wp *WorkerPool) String() string {
	return fmt.Sprintf(workerPoolStringFormat, wp.name, wp.minWorkers, wp.maxWorkers, wp.idleTimeout.String(), wp.upScaling, wp.downScaling, cap(wp.queue))
}

func (wp *WorkerPool) createWorkers(workers int) *sync.WaitGroup {
	wp.workerWG.Add(workers)
	atomic.AddInt32(&wp.runningWorkers, int32(workers))
	atomic.AddInt32(&wp.idleWorkers, int32(workers))

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go wp.worker(&wg)
	}

	return &wg
}

func (wp *WorkerPool) worker(wg *sync.WaitGroup) {
	defer func() {
		atomic.AddInt32(&wp.runningWorkers, -1)
		atomic.AddInt32(&wp.idleWorkers, -1)
		wp.workerWG.Done()
	}()

	wg.Done()

	for {
		select {
		case <-wp.stoppedSignal:
			// Worker Pool is stopped, finish
			return
		case <-wp.purgeSignal:
			// Purge worker
			return
		case task := <-wp.queue:
			atomic.AddInt32(&wp.idleWorkers, -1)
			wp.executeTask(task)
			atomic.AddInt32(&wp.idleWorkers, 1)
		}
	}
}

func (wp *WorkerPool) executeTask(task func()) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			atomic.AddUint64(&wp.failedTasks, 1)
			wp.panicHandler(panicErr)
		}
	}()

	atomic.AddUint64(&wp.waitingTasks, ^uint64(0))
	task()
	atomic.AddUint64(&wp.successfulTasks, 1)
}

func (wp *WorkerPool) submit(task func(), mustSubmit bool, mustWait bool, waitDeadline *time.Time) (submitted bool, finished bool, doneSignal chan struct{}) {
	if task == nil {
		return true, true, nil
	}

	if wp.stopped {
		if mustSubmit {
			panic(fmt.Errorf("WorkerPool %s is stopped and is no longer accepting new tasks", wp.name))
		}

		return false, false, nil
	}

	if mustWait {
		doneSignal = make(chan struct{})

		originalTask := task
		task = func() {
			defer func() {
				finished = true
				close(doneSignal)
			}()

			originalTask()
		}
	}

	atomic.AddUint64(&wp.submittedTasks, 1)
	atomic.AddUint64(&wp.waitingTasks, 1)

	// attempt to submit task without blocking
	select {
	case wp.queue <- task:
		// submitted, continue
		submitted = true
	default:
		// queue is full and no idle workers available, continue
	}

	// attempt to scale if no workers are available
	workersScaled := 0
	if wp.metrics.IdleWorkers() == 0 {
		workersScaled = wp.ScaleUp(wp.upScaling)
	}

	// try to resubmit task
	if !submitted {
		// if no workers are scaled, try to scale a worker to process this task
		if workersScaled == 0 {
			workersScaled = wp.ScaleUp(1)
		}

		if mustSubmit || workersScaled > 0 {
			// enqueue task if must submit or at least one worker is created to process this task
			wp.queue <- task
		} else {
			// no workers can be created, decrement counters and return
			atomic.AddUint64(&wp.submittedTasks, ^uint64(0))
			atomic.AddUint64(&wp.waitingTasks, ^uint64(0))
			return false, false, nil
		}
	}

	if mustWait {
		if waitDeadline == nil {
			<-doneSignal
			return true, true, doneSignal
		} else {
			select {
			case <-doneSignal:
				return true, true, doneSignal
			case <-time.After(time.Until(*waitDeadline)):
				return true, finished, doneSignal
			}
		}
	}

	return true, finished, nil
}

func (wp *WorkerPool) monitor() {
	ticker := time.NewTicker(wp.idleTimeout)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-wp.idleTimeoutChanged:
			ticker.Stop()
			ticker = time.NewTicker(wp.idleTimeout)
		case <-ticker.C:
			wp.handlePurgeIdleWorkers()
		case <-wp.stoppedSignal:
			return
		}
	}
}

func (wp *WorkerPool) handlePurgeIdleWorkers() {
	workersToBePurged := wp.downScaling
	idleWorkers := wp.metrics.IdleWorkers()

	if workersToBePurged > idleWorkers {
		workersToBePurged = idleWorkers
	}

	maxPurge := wp.metrics.RunningWorkers() - wp.minWorkers
	if workersToBePurged > maxPurge {
		workersToBePurged = maxPurge
	}

	for i := 0; i < workersToBePurged; i++ {
		wp.purgeSignal <- struct{}{}
	}
}

func defaultPanicHandler(panicErr interface{}) {
	fmt.Printf("Worker recovered from a panic: %v\nStack trace: %s\n", panicErr, string(debug.Stack()))
}

const (
	workerPoolStringFormat = "WorkerPool [name=%s, minWorkers=%d, maxWorkers=%d, idleTimeout=%s, upScaling=%d, downScaling=%d, queueCapacity=%v]"
	defaultIdleTimeout     = 15 * time.Second
)
