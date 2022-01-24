package workerpool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	testCases := []struct {
		name                   string
		mustSubmit             bool
		mustWait               bool
		hasWaitDeadline        bool
		hasPanic               bool
		waitTaskFinishDeadline bool
		operation              func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{})
	}{
		{
			name:            "Submit",
			mustSubmit:      true,
			mustWait:        false,
			hasWaitDeadline: false,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				wp.Submit(task)
				if task == nil {
					finished = true
				}

				return true, finished, nil
			},
		},
		{
			name:            "SubmitAndWait",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: false,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				wp.SubmitAndWait(task)
				return true, true, nil
			},
		},
		{
			name:            "SubmitAndWaitWithTimeout_ShouldReleaseBeforeTimeout",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithTimeout(*timeout, task)
				return true, finished, wait
			},
		},
		{
			name:                   "SubmitAndWaitWithTimeout_ShouldReleaseWhenTaskFinish",
			mustSubmit:             true,
			mustWait:               true,
			hasWaitDeadline:        true,
			hasPanic:               false,
			waitTaskFinishDeadline: true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithTimeout(*timeout, task)
				return true, finished, wait
			},
		},
		{
			name:            "SubmitAndWaitWithTimeoutWithDeadline_ShouldReleaseBeforeTimeout",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return true, finished, wait
			},
		},
		{
			name:                   "SubmitAndWaitWithTimeoutWithDeadline_ShouldReleaseWhenTaskFinish",
			mustSubmit:             true,
			mustWait:               true,
			hasWaitDeadline:        true,
			hasPanic:               false,
			waitTaskFinishDeadline: true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return true, finished, wait
			},
		},
		{
			name:            "TrySubmit",
			mustSubmit:      false,
			mustWait:        false,
			hasWaitDeadline: false,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted = wp.TrySubmit(task)
				if task == nil {
					finished = true
				}

				return submitted, finished, nil
			},
		},
		{
			name:            "TrySubmitAndWait",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: false,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted = wp.TrySubmitAndWait(task)
				if task == nil {
					finished = true
				}

				return submitted, finished, nil
			},
		},
		{
			name:            "TrySubmitAndWaitWithTimeout_ShouldReleaseBeforeTimeout",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithTimeout(*timeout, task)
				return submitted, finished, wait
			},
		},
		{
			name:                   "TrySubmitAndWaitWithTimeout_ShouldReleaseWhenTaskFinish",
			mustSubmit:             false,
			mustWait:               true,
			hasWaitDeadline:        true,
			hasPanic:               false,
			waitTaskFinishDeadline: true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithTimeout(*timeout, task)
				return submitted, finished, wait
			},
		},
		{
			name:            "TrySubmitAndWaitWithDeadline_ShouldReleaseBeforeTimeout",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        false,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return submitted, finished, wait
			},
		},
		{
			name:                   "TrySubmitAndWaitWithDeadline_ShouldReleaseWhenTaskFinish",
			mustSubmit:             false,
			mustWait:               true,
			hasWaitDeadline:        true,
			hasPanic:               false,
			waitTaskFinishDeadline: true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return submitted, finished, wait
			},
		},
		// ====================================================================================================
		// Panic Test cases
		// ====================================================================================================
		{
			name:            "Submit",
			mustSubmit:      true,
			mustWait:        false,
			hasWaitDeadline: false,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				wp.Submit(task)
				if task == nil {
					finished = true
				}

				return true, finished, nil
			},
		},
		{
			name:            "SubmitAndWait",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: false,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				wp.SubmitAndWait(task)
				return true, true, nil
			},
		},
		{
			name:            "SubmitAndWaitWithTimeout",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithTimeout(*timeout, task)
				return true, finished, wait
			},
		},
		{
			name:            "SubmitAndWaitWithTimeoutWithDeadline",
			mustSubmit:      true,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				finished, wait = wp.SubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return true, finished, wait
			},
		},
		{
			name:            "TrySubmit",
			mustSubmit:      false,
			mustWait:        false,
			hasWaitDeadline: false,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted = wp.TrySubmit(task)
				if task == nil {
					finished = true
				}

				return submitted, finished, nil
			},
		},
		{
			name:            "TrySubmitAndWait",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: false,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted = wp.TrySubmitAndWait(task)
				if task == nil {
					finished = true
				}

				return submitted, finished, nil
			},
		},
		{
			name:            "TrySubmitAndWaitWithTimeout",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithTimeout(*timeout, task)
				return submitted, finished, wait
			},
		},
		{
			name:            "TrySubmitAndWaitWithDeadline",
			mustSubmit:      false,
			mustWait:        true,
			hasWaitDeadline: true,
			hasPanic:        true,
			operation: func(wp *WorkerPool, task func(), timeout *time.Duration) (submitted bool, finished bool, wait <-chan struct{}) {
				submitted, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(*timeout), task)
				return submitted, finished, wait
			},
		},
	}

	for _, testCase := range testCases {
		var testName string

		if testCase.hasPanic {
			testName = fmt.Sprintf("%s_ShouldRecoverFromPanic", testCase.name)
		} else {
			testName = testCase.name
		}

		t.Run(testName, func(t *testing.T) {
			wp := New("wp", &Settings{
				MinWorkers: 1,
				MaxWorkers: 4,
				Queue:      1,
				UpScaling:  1,
			})

			// task duration
			taskDuration := 30 * time.Millisecond
			taskDurationMin := taskDuration - 5*time.Millisecond // max of 3 milliseconds of overhead
			taskDurationMax := taskDuration + 5*time.Millisecond // max of 3 milliseconds of overhead

			// timeout for non-blocking submissions
			submitTimeout := 1 * time.Millisecond

			// test panic
			expectedPanicErr := fmt.Errorf("panic error on %s", testCase.name)
			var panicErr interface{}

			// used when the operation has a timeout to wait for the task to finish
			var waitTimeout time.Duration
			if testCase.waitTaskFinishDeadline {
				waitTimeout = taskDuration * 2
			} else {
				waitTimeout = taskDuration / 2
			}

			waitTimeoutMin := waitTimeout - 5*time.Millisecond
			waitTimeoutMax := waitTimeout + 5*time.Millisecond

			// task results
			var submittedAt time.Time
			submittedChan := make(chan struct{})
			taskFinished := false
			taskFinishedChan := make(chan struct{})

			task := func() {
				defer func() {
					taskFinished = true
					close(taskFinishedChan)
				}()

				submittedAt = time.Now()
				close(submittedChan)
				time.Sleep(taskDuration)

				if testCase.hasPanic {
					panic(expectedPanicErr)
				}
			}

			// test nil task submit
			submitted, finished, wait := testCase.operation(wp, nil, &waitTimeout)
			assert.True(t, submitted, fmt.Sprintf("[%s] should mark nil task as submitted", testCase.name))
			assert.True(t, finished, fmt.Sprintf("[%s] should mark nil task as finished", testCase.name))
			assert.Nil(t, wait, fmt.Sprintf("[%s] should return nil channel on when task is nil", testCase.name))

			// test no running workers
			//assert.Equal(t, 0, wp.Metrics().IdleWorkers())
			//assert.Equal(t, 0, wp.Metrics().RunningWorkers())

			// start submit operation
			start := time.Now()
			submitted, finished, wait = testCase.operation(wp, task, &waitTimeout)
			waitDuration := time.Since(start)

			assert.True(t, submitted, fmt.Sprintf("[%s] should submit task", testCase.name))

			// test scale up
			//assert.Equal(t, 1, wp.Metrics().IdleWorkers())
			//assert.Equal(t, 2, wp.Metrics().RunningWorkers())

			// test non-blocking submit
			<-submittedChan
			submittedIn := submittedAt.Sub(start)
			assert.LessOrEqual(t, submittedIn, submitTimeout, fmt.Sprintf("[%s] should be non-blocking (target: %s | timeout: %s)", testCase.name, submitTimeout, submittedIn))

			if testCase.mustWait {
				if testCase.hasWaitDeadline {

					if testCase.waitTaskFinishDeadline {
						assert.GreaterOrEqual(t, waitDuration, taskDurationMin, fmt.Sprintf("[%s] should wait between [%s and %s] but was %s", testCase.name, taskDurationMin, taskDurationMax, waitDuration))
						assert.LessOrEqual(t, waitDuration, taskDurationMax, fmt.Sprintf("[%s] should wait between [%s and %s] but was %s", testCase.name, taskDurationMin, taskDurationMax, waitDuration))
						assert.True(t, taskFinished, "Task should be finished")
						assert.True(t, finished, "Task should be finished")
					} else {
						assert.GreaterOrEqual(t, waitDuration, waitTimeoutMin, fmt.Sprintf("[%s] should wait between [%s and %s] but was %s", testCase.name, waitTimeoutMin, waitTimeoutMax, waitDuration))
						assert.LessOrEqual(t, waitDuration, waitTimeoutMax, fmt.Sprintf("[%s] should wait between [%s and %s] but was %s", testCase.name, waitTimeoutMin, waitTimeoutMax, waitDuration))
						assert.False(t, taskFinished, "Task should not be finished yet")
						assert.False(t, finished, "Task should not be finished yet")
					}

					assert.NotNil(t, wait, "Should have wait channel")

					<-wait
					finishedIn := time.Since(start)

					assert.True(t, taskFinished, fmt.Sprintf("[%s] task should be finished after read wait channel", testCase.name))
					assert.LessOrEqual(t, finishedIn, taskDurationMax, fmt.Sprintf("[%s] task should be finished in %s but was finished in %s", testCase.name, taskDurationMax, finishedIn))

					select {
					case <-wait:
						// everything ok because task is finished
					default:
						assert.FailNow(t, "After task finished, wait channel should always be non-blocking (closed channel)")
					}
				} else {
					assert.Nil(t, wait, fmt.Sprintf("[%s] should not have wait channel on return", testCase.name))
					assert.True(t, taskFinished, fmt.Sprintf("[%s] should wait task to finish", testCase.name))

					finishedIn := time.Since(start)
					assert.LessOrEqual(t, finishedIn, taskDurationMax, fmt.Sprintf("[%s] task should be finished in %s but was finished in %s", testCase.name, taskDurationMax, finishedIn))
				}
			} else {
				assert.False(t, taskFinished, fmt.Sprintf("[%s] should not wait task to finish", testCase.name))
				<-taskFinishedChan
				assert.Nil(t, wait, fmt.Sprintf("[%s] should not have wait channel on return", testCase.name))
				assert.True(t, taskFinished, fmt.Sprintf("[%s] should wait task to finish", testCase.name))

				finishedIn := time.Since(start)
				assert.LessOrEqual(t, finishedIn, taskDurationMax, fmt.Sprintf("[%s] task should be finished in %s but was finished in %s", testCase.name, taskDurationMax, finishedIn))
			}

			if testCase.hasPanic {
				// recovered from default panic

				// now, check custom panic handler
				wp.SetPanicHandler(func(pErr interface{}) {
					panicErr = pErr
				})

				submittedChan = make(chan struct{})
				taskFinishedChan = make(chan struct{})

				testCase.operation(wp, task, &waitTimeout)
				<-taskFinishedChan

				assert.Equal(t, expectedPanicErr, panicErr, fmt.Sprintf("[%s] not recovered from panic", testCase.name))
			}

			// test try submissions
			if !testCase.mustSubmit {
				// keep the queue and workers busy
				for i := 0; i < wp.QueueCapacity()+wp.MaxWorkers()+1; i++ {
					wp.TrySubmit(func() {
						time.Sleep(50 * time.Millisecond)
					})
				}

				// task results
				submittedChan = make(chan struct{})
				taskFinished = false
				taskFinishedChan = make(chan struct{})

				start := time.Now()
				submitted, finished, wait = testCase.operation(wp, task, &waitTimeout)
				waitDuration := time.Since(start)

				assert.False(t, submitted, fmt.Sprintf("[%s] should not submit task", testCase.name))
				assert.False(t, finished, fmt.Sprintf("[%s] should not finish task", testCase.name))
				assert.Nil(t, wait, fmt.Sprintf("[%s] should not have wait channel because it's should not be submitted", testCase.name))
				assert.LessOrEqual(t, waitDuration, submitTimeout, fmt.Sprintf("[%s] should be non-blocking when queue is full (target: %s | timeout: %s)", testCase.name, submitTimeout, waitDuration))
			}
		})
	}
}

func TestWorkerPool_Burst(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 4,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	wp.Burst(-1)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	wp.Burst(0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	wp.Burst(4)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 8, wp.Metrics().IdleWorkers())
	assert.Equal(t, 8, wp.Metrics().RunningWorkers())

	wp.Burst(2)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 10, wp.Metrics().IdleWorkers())
	assert.Equal(t, 10, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_ScaleUp(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleUp := wp.ScaleUp(-1)
	assert.Equal(t, scaleUp, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleUp = wp.ScaleUp(0)
	assert.Equal(t, scaleUp, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleUp = wp.ScaleUp(4)
	assert.Equal(t, scaleUp, 4)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 8, wp.Metrics().IdleWorkers())
	assert.Equal(t, 8, wp.Metrics().RunningWorkers())

	scaleUp = wp.ScaleUp(12)
	assert.Equal(t, scaleUp, 8)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())

	scaleUp = wp.ScaleUp(1)
	assert.Equal(t, scaleUp, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_ScaleDown(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleUp := wp.ScaleUp(12)
	assert.Equal(t, scaleUp, 12)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())

	scaleDown := wp.ScaleDown(-1)
	assert.Equal(t, scaleDown, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())

	scaleDown = wp.ScaleDown(0)
	assert.Equal(t, scaleDown, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())

	scaleDown = wp.ScaleDown(4)
	assert.Equal(t, scaleDown, 4)
	time.Sleep(5 * time.Millisecond) // wait for scale down
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 12, wp.Metrics().IdleWorkers())
	assert.Equal(t, 12, wp.Metrics().RunningWorkers())

	scaleDown = wp.ScaleDown(12)
	assert.Equal(t, scaleDown, 8)
	time.Sleep(5 * time.Millisecond) // wait for scale down
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleDown = wp.ScaleDown(1)
	assert.Equal(t, scaleDown, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_ReleaseIdleWorkers(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	scaleUp := wp.ScaleUp(12)
	assert.Equal(t, scaleUp, 12)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 16, wp.Metrics().IdleWorkers())
	assert.Equal(t, 16, wp.Metrics().RunningWorkers())

	idleWorkers := wp.ReleaseIdleWorkers()
	assert.Equal(t, idleWorkers, 12)
	time.Sleep(5 * time.Millisecond) // wait for release idle workers
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	idleWorkers = wp.ReleaseIdleWorkers()
	assert.Equal(t, idleWorkers, 0)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_Stop(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	wp.Stop()

	// wait workers to be stopped
	time.Sleep(5 * time.Millisecond)
	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait := wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(0), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(0), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_StopAndWait(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	done := false
	taskDuration := 50 * time.Millisecond
	taskDurationMin := taskDuration - 2*time.Millisecond
	taskDurationMax := taskDuration + 2*time.Millisecond

	wp.Submit(func() {
		time.Sleep(taskDuration)
		done = true
	})

	stoppedAt := time.Now()
	wp.StopAndWait()
	stoppedWaitDuration := time.Since(stoppedAt)

	assert.GreaterOrEqual(t, stoppedWaitDuration, taskDurationMin, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.LessOrEqual(t, stoppedWaitDuration, taskDurationMax, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.True(t, done)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait := wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_StopAndWaitWithTimeout_StopBeforeTimeout(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	done := false
	taskDuration := 50 * time.Millisecond
	taskDurationMin := taskDuration - 2*time.Millisecond
	taskDurationMax := taskDuration + 2*time.Millisecond

	stopTimeout := taskDuration * 2

	wp.Submit(func() {
		time.Sleep(taskDuration)
		done = true
	})

	stoppedAt := time.Now()
	finished, wait := wp.StopAndWaitWithTimeout(stopTimeout)
	stoppedWaitDuration := time.Since(stoppedAt)

	assert.True(t, finished)
	assert.NotNil(t, wait)
	assert.GreaterOrEqual(t, stoppedWaitDuration, taskDurationMin, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.LessOrEqual(t, stoppedWaitDuration, taskDurationMax, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.True(t, done)

	select {
	case <-wait:
	// ok, closed channels always return
	default:
		assert.FailNow(t, "After worker pool stopped, wait channel should always be non-blocking (closed channel)")
	}

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_StopAndWaitWithTimeout_StopAfterTimeout(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	done := false
	taskDuration := 50 * time.Millisecond
	taskDurationMin := taskDuration - 2*time.Millisecond
	taskDurationMax := taskDuration + 2*time.Millisecond

	stopTimeout := taskDuration / 2
	stopTimeoutMin := stopTimeout - 2*time.Millisecond
	stopTimeoutMax := stopTimeout + 2*time.Millisecond

	submittedAt := time.Now()
	wp.Submit(func() {
		time.Sleep(taskDuration)
		done = true
	})

	stoppedAt := time.Now()
	finished, wait := wp.StopAndWaitWithTimeout(stopTimeout)
	stoppedWaitDuration := time.Since(stoppedAt)

	assert.False(t, finished)
	assert.NotNil(t, wait)
	assert.GreaterOrEqual(t, stoppedWaitDuration, stopTimeoutMin, fmt.Sprintf("should wait between [%s and %s] but was %s", stopTimeoutMin, stopTimeoutMax, stoppedWaitDuration))
	assert.LessOrEqual(t, stoppedWaitDuration, stopTimeoutMax, fmt.Sprintf("should wait between [%s and %s] but was %s", stopTimeoutMin, stopTimeoutMax, stoppedWaitDuration))
	assert.False(t, done)

	<-wait
	taskWaitDuration := time.Since(submittedAt)
	assert.GreaterOrEqual(t, taskWaitDuration, taskDurationMin, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, taskWaitDuration))
	assert.LessOrEqual(t, taskWaitDuration, taskDurationMax, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, taskWaitDuration))
	assert.True(t, done)

	select {
	case <-wait:
	// ok, closed channels always return
	default:
		assert.FailNow(t, "After worker pool stopped, wait channel should always be non-blocking (closed channel)")
	}

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_StopAndWaitWithDeadline_StopBeforeTimeout(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	done := false
	taskDuration := 50 * time.Millisecond
	taskDurationMin := taskDuration - 2*time.Millisecond
	taskDurationMax := taskDuration + 2*time.Millisecond

	stopTimeout := taskDuration * 2

	wp.Submit(func() {
		time.Sleep(taskDuration)
		done = true
	})

	stoppedAt := time.Now()
	finished, wait := wp.StopAndWaitWithDeadline(time.Now().Add(stopTimeout))
	stoppedWaitDuration := time.Since(stoppedAt)

	assert.True(t, finished)
	assert.NotNil(t, wait)
	assert.GreaterOrEqual(t, stoppedWaitDuration, taskDurationMin, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.LessOrEqual(t, stoppedWaitDuration, taskDurationMax, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, stoppedWaitDuration))
	assert.True(t, done)

	select {
	case <-wait:
	// ok, closed channels always return
	default:
		assert.FailNow(t, "After worker pool stopped, wait channel should always be non-blocking (closed channel)")
	}

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_StopAndWaitWithDeadline_StopAfterTimeout(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
		Queue:      16,
	})

	done := false
	taskDuration := 50 * time.Millisecond
	taskDurationMin := taskDuration - 2*time.Millisecond
	taskDurationMax := taskDuration + 2*time.Millisecond

	stopTimeout := taskDuration / 2
	stopTimeoutMin := stopTimeout - 2*time.Millisecond
	stopTimeoutMax := stopTimeout + 2*time.Millisecond

	submittedAt := time.Now()
	wp.Submit(func() {
		time.Sleep(taskDuration)
		done = true
	})

	stoppedAt := time.Now()
	finished, wait := wp.StopAndWaitWithDeadline(time.Now().Add(stopTimeout))
	stoppedWaitDuration := time.Since(stoppedAt)

	assert.False(t, finished)
	assert.NotNil(t, wait)
	assert.GreaterOrEqual(t, stoppedWaitDuration, stopTimeoutMin, fmt.Sprintf("should wait between [%s and %s] but was %s", stopTimeoutMin, stopTimeoutMax, stoppedWaitDuration))
	assert.LessOrEqual(t, stoppedWaitDuration, stopTimeoutMax, fmt.Sprintf("should wait between [%s and %s] but was %s", stopTimeoutMin, stopTimeoutMax, stoppedWaitDuration))
	assert.False(t, done)

	<-wait
	taskWaitDuration := time.Since(submittedAt)
	assert.GreaterOrEqual(t, taskWaitDuration, taskDurationMin, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, taskWaitDuration))
	assert.LessOrEqual(t, taskWaitDuration, taskDurationMax, fmt.Sprintf("should wait between [%s and %s] but was %s", taskDurationMin, taskDurationMax, taskWaitDuration))
	assert.True(t, done)

	select {
	case <-wait:
	// ok, closed channels always return
	default:
		assert.FailNow(t, "After worker pool stopped, wait channel should always be non-blocking (closed channel)")
	}

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.Submit(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWait(func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithTimeout(1*time.Second, func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	assertPanic(t, func() {
		wp.SubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	})

	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit := wp.TrySubmit(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit = wp.TrySubmitAndWait(func() {})
	assert.False(t, submit)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithTimeout(1*time.Second, func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())

	submit, finished, wait = wp.TrySubmitAndWaitWithDeadline(time.Now().Add(1*time.Second), func() {})
	assert.False(t, submit)
	assert.False(t, finished)
	assert.Nil(t, wait)
	assert.Equal(t, uint64(1), wp.Metrics().SubmittedTasks())
	assert.Equal(t, uint64(1), wp.Metrics().CompletedTasks())
	assert.Equal(t, 0, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_Stopped(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.False(t, wp.Stopped())

	wp.stopped = true
	assert.True(t, wp.Stopped())
}

func TestWorkerPool_Metrics(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, reflect.ValueOf(wp.Metrics()).Pointer(), reflect.ValueOf(wp.metrics).Pointer())
}

func TestWorkerPool_MinWorkers(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, 0, wp.MinWorkers())
	assert.Equal(t, 0, wp.Metrics().MinWorkers())

	wp = New("wp", &Settings{
		MinWorkers: 0,
	})
	assert.Equal(t, 0, wp.MinWorkers())
	assert.Equal(t, 0, wp.Metrics().MinWorkers())

	wp = New("wp", &Settings{
		MinWorkers: -1,
	})
	assert.Equal(t, 0, wp.MinWorkers())
	assert.Equal(t, 0, wp.Metrics().MinWorkers())

	wp = New("wp", &Settings{
		MinWorkers: 4,
	})
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 2,
	})
	assert.Equal(t, 2, wp.MinWorkers())
	assert.Equal(t, 2, wp.Metrics().MinWorkers())
	assert.Equal(t, 2, wp.MaxWorkers())
	assert.Equal(t, 2, wp.Metrics().MaxWorkers())
}

func TestWorkerPool_SetMinWorkers(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersCreated := wp.SetMinWorkers(-1)
	assert.Equal(t, 0, workersCreated)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersCreated = wp.SetMinWorkers(8)
	assert.Equal(t, 4, workersCreated)
	assert.Equal(t, 8, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 8, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 8, wp.Metrics().IdleWorkers())
	assert.Equal(t, 8, wp.Metrics().RunningWorkers())

	workersCreated = wp.SetMinWorkers(20)
	assert.Equal(t, 12, workersCreated)
	assert.Equal(t, 20, wp.MinWorkers())
	assert.Equal(t, 20, wp.MaxWorkers())
	assert.Equal(t, 20, wp.Metrics().MinWorkers())
	assert.Equal(t, 20, wp.Metrics().MaxWorkers())
	assert.Equal(t, 20, wp.Metrics().IdleWorkers())
	assert.Equal(t, 20, wp.Metrics().RunningWorkers())

	workersCreated = wp.SetMinWorkers(0)
	assert.Equal(t, 0, workersCreated)
	assert.Equal(t, 0, wp.MinWorkers())
	assert.Equal(t, 20, wp.MaxWorkers())
	assert.Equal(t, 0, wp.Metrics().MinWorkers())
	assert.Equal(t, 20, wp.Metrics().MaxWorkers())
	assert.Equal(t, 20, wp.Metrics().IdleWorkers())
	assert.Equal(t, 20, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_MaxWorkers(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, runtime.NumCPU(), wp.MaxWorkers())
	assert.Equal(t, runtime.NumCPU(), wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MaxWorkers: 0,
	})
	assert.Equal(t, runtime.NumCPU(), wp.MaxWorkers())
	assert.Equal(t, runtime.NumCPU(), wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MinWorkers: -1,
	})
	assert.Equal(t, runtime.NumCPU(), wp.MaxWorkers())
	assert.Equal(t, runtime.NumCPU(), wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MaxWorkers: 4,
	})
	assert.Equal(t, 0, wp.MinWorkers())
	assert.Equal(t, 0, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 2,
	})
	assert.Equal(t, 2, wp.MinWorkers())
	assert.Equal(t, 2, wp.Metrics().MinWorkers())
	assert.Equal(t, 2, wp.MaxWorkers())
	assert.Equal(t, 2, wp.Metrics().MaxWorkers())

	wp = New("wp", &Settings{
		MinWorkers: 2,
		MaxWorkers: 4,
	})
	assert.Equal(t, 2, wp.MinWorkers())
	assert.Equal(t, 2, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
}

func TestWorkerPool_SetMaxWorkers(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers: 4,
		MaxWorkers: 16,
	})

	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersToBePurged := wp.SetMaxWorkers(-1)
	assert.Equal(t, 0, workersToBePurged)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(0)
	assert.Equal(t, 0, workersToBePurged)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 16, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 16, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(20)
	assert.Equal(t, 0, workersToBePurged)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 20, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 20, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	wp.Burst(8)
	assert.Equal(t, 0, workersToBePurged)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 20, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 20, wp.Metrics().MaxWorkers())
	assert.Equal(t, 12, wp.Metrics().IdleWorkers())
	assert.Equal(t, 12, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(12)
	assert.Equal(t, 0, workersToBePurged)
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 12, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 12, wp.Metrics().MaxWorkers())
	assert.Equal(t, 12, wp.Metrics().IdleWorkers())
	assert.Equal(t, 12, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(8)
	assert.Equal(t, 4, workersToBePurged)
	time.Sleep(5 * time.Millisecond) // wait for workers to be purged
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 8, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 8, wp.Metrics().MaxWorkers())
	assert.Equal(t, 8, wp.Metrics().IdleWorkers())
	assert.Equal(t, 8, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(4)
	assert.Equal(t, 4, workersToBePurged)
	time.Sleep(5 * time.Millisecond) // wait for workers to be purged
	assert.Equal(t, 4, wp.MinWorkers())
	assert.Equal(t, 4, wp.MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().MinWorkers())
	assert.Equal(t, 4, wp.Metrics().MaxWorkers())
	assert.Equal(t, 4, wp.Metrics().IdleWorkers())
	assert.Equal(t, 4, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(3)
	assert.Equal(t, 1, workersToBePurged)
	time.Sleep(5 * time.Millisecond) // wait for workers to be purged
	assert.Equal(t, 3, wp.MinWorkers())
	assert.Equal(t, 3, wp.MaxWorkers())
	assert.Equal(t, 3, wp.Metrics().MinWorkers())
	assert.Equal(t, 3, wp.Metrics().MaxWorkers())
	assert.Equal(t, 3, wp.Metrics().IdleWorkers())
	assert.Equal(t, 3, wp.Metrics().RunningWorkers())

	workersToBePurged = wp.SetMaxWorkers(1)
	assert.Equal(t, 2, workersToBePurged)
	time.Sleep(5 * time.Millisecond) // wait for workers to be purged
	assert.Equal(t, 1, wp.MinWorkers())
	assert.Equal(t, 1, wp.MaxWorkers())
	assert.Equal(t, 1, wp.Metrics().MinWorkers())
	assert.Equal(t, 1, wp.Metrics().MaxWorkers())
	assert.Equal(t, 1, wp.Metrics().IdleWorkers())
	assert.Equal(t, 1, wp.Metrics().RunningWorkers())
}

func TestWorkerPool_IdleTimeout(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, 15*time.Second, wp.idleTimeout)
	assert.Equal(t, 15*time.Second, wp.IdleTimeout())

	wp = New("wp", &Settings{
		IdleTimeout: 0,
	})

	assert.Equal(t, 15*time.Second, wp.idleTimeout)
	assert.Equal(t, 15*time.Second, wp.IdleTimeout())

	wp = New("wp", &Settings{
		IdleTimeout: -1,
	})

	assert.Equal(t, 15*time.Second, wp.idleTimeout)
	assert.Equal(t, 15*time.Second, wp.IdleTimeout())

	wp = New("wp", &Settings{
		IdleTimeout: 20 * time.Millisecond,
	})

	assert.Equal(t, 20*time.Millisecond, wp.idleTimeout)
	assert.Equal(t, 20*time.Millisecond, wp.IdleTimeout())
}

func TestWorkerPool_SetIdleTimeout(t *testing.T) {
	wp := New("wp", &Settings{
		MinWorkers:  4,
		MaxWorkers:  16,
		Queue:       16,
		IdleTimeout: 1 * time.Minute,
		DownScaling: 8,
	})

	wp.SetIdleTimeout(0)
	assert.Equal(t, 15*time.Second, wp.idleTimeout)
	assert.Equal(t, 15*time.Second, wp.IdleTimeout())

	wp.SetIdleTimeout(-1)
	assert.Equal(t, 15*time.Second, wp.idleTimeout)
	assert.Equal(t, 15*time.Second, wp.IdleTimeout())

	wp.Burst(12)
	wp.SetIdleTimeout(100 * time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, wp.idleTimeout)
	assert.Equal(t, 100*time.Millisecond, wp.IdleTimeout())

	time.Sleep(105 * time.Millisecond)
	assert.Equal(t, 8, wp.metrics.IdleWorkers())
	assert.Equal(t, 8, wp.metrics.RunningWorkers())

	time.Sleep(105 * time.Millisecond)
	assert.Equal(t, 4, wp.metrics.IdleWorkers())
	assert.Equal(t, 4, wp.metrics.RunningWorkers())

	time.Sleep(105 * time.Millisecond)
	assert.Equal(t, 4, wp.metrics.IdleWorkers())
	assert.Equal(t, 4, wp.metrics.RunningWorkers())
}

func TestWorkerPool_UpScaling(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp = New("wp", &Settings{
		UpScaling: -1,
	})
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp = New("wp", &Settings{
		UpScaling: 0,
	})
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp = New("wp", &Settings{
		UpScaling: 1,
	})
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp = New("wp", &Settings{
		UpScaling: 4,
	})
	assert.Equal(t, 4, wp.upScaling)
	assert.Equal(t, 4, wp.UpScaling())
}

func TestWorkerPool_SetUpScaling(t *testing.T) {
	wp := New("wp", &Settings{
		UpScaling: 4,
	})

	wp.SetUpScaling(-1)
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp.SetUpScaling(0)
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp.SetUpScaling(1)
	assert.Equal(t, 1, wp.upScaling)
	assert.Equal(t, 1, wp.UpScaling())

	wp.SetUpScaling(4)
	assert.Equal(t, 4, wp.upScaling)
	assert.Equal(t, 4, wp.UpScaling())
}

func TestWorkerPool_DownScaling(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp = New("wp", &Settings{
		DownScaling: -1,
	})
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp = New("wp", &Settings{
		DownScaling: 0,
	})
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp = New("wp", &Settings{
		DownScaling: 1,
	})
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp = New("wp", &Settings{
		DownScaling: 4,
	})
	assert.Equal(t, 4, wp.downScaling)
	assert.Equal(t, 4, wp.DownScaling())
}

func TestWorkerPool_SetDownScaling(t *testing.T) {
	wp := New("wp", &Settings{
		DownScaling: 4,
	})

	wp.SetDownScaling(-1)
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp.SetDownScaling(0)
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp.SetDownScaling(1)
	assert.Equal(t, 1, wp.downScaling)
	assert.Equal(t, 1, wp.DownScaling())

	wp.SetDownScaling(4)
	assert.Equal(t, 4, wp.downScaling)
	assert.Equal(t, 4, wp.DownScaling())
}

func TestWorkerPool_QueueCapacity(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, 0, len(wp.queue))
	assert.Equal(t, 0, cap(wp.queue))
	assert.Equal(t, 0, wp.QueueCapacity())

	wp = New("wp", &Settings{
		Queue: -1,
	})
	assert.Equal(t, 0, len(wp.queue))
	assert.Equal(t, 0, cap(wp.queue))
	assert.Equal(t, 0, wp.QueueCapacity())

	wp = New("wp", &Settings{
		Queue: 0,
	})
	assert.Equal(t, 0, len(wp.queue))
	assert.Equal(t, 0, cap(wp.queue))
	assert.Equal(t, 0, wp.QueueCapacity())

	wp = New("wp", &Settings{
		Queue: 1,
	})
	assert.Equal(t, 0, len(wp.queue))
	assert.Equal(t, 1, cap(wp.queue))
	assert.Equal(t, 1, wp.QueueCapacity())

	wp = New("wp", &Settings{
		Queue: 4,
	})
	assert.Equal(t, 0, len(wp.queue))
	assert.Equal(t, 4, cap(wp.queue))
	assert.Equal(t, 4, wp.QueueCapacity())
}

func TestWorkerPool_PanicHandler(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())

	wp = New("wp", &Settings{
		PanicHandler: nil,
	})
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())

	panicHandler := func(panicError interface{}) {}

	wp = New("wp", &Settings{
		PanicHandler: panicHandler,
	})
	assert.Equal(t, reflect.ValueOf(panicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(panicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())
}

func TestWorkerPool_SetPanicHandler(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())

	panicHandler := func(panicError interface{}) {}
	wp.SetPanicHandler(panicHandler)
	assert.Equal(t, reflect.ValueOf(panicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(panicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())

	wp.SetPanicHandler(nil)
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.panicHandler).Pointer())
	assert.Equal(t, reflect.ValueOf(defaultPanicHandler).Pointer(), reflect.ValueOf(wp.PanicHandler()).Pointer())
}

func TestWorkerPool_String(t *testing.T) {
	wp := New("wp", &Settings{})
	assert.Equal(t, "WorkerPool [name=wp, minWorkers=0, maxWorkers="+strconv.Itoa(runtime.NumCPU())+", idleTimeout=15s, upScaling=1, downScaling=1, queueCapacity=0]", wp.String())

	wp = New("wp", &Settings{
		MinWorkers:  1,
		MaxWorkers:  2,
		IdleTimeout: 3 * time.Second,
		UpScaling:   4,
		DownScaling: 5,
		Queue:       6,
	})
	assert.Equal(t, "WorkerPool [name=wp, minWorkers=1, maxWorkers=2, idleTimeout=3s, upScaling=4, downScaling=5, queueCapacity=6]", wp.String())
}

func assertPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r == nil {
			assert.FailNow(t, "the code did not panic")
		}
	}()

	fn()
}
