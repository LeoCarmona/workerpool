package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetrics_MinWorkers(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, 0, metrics.MinWorkers())

	metrics.wp.minWorkers = 1
	assert.Equal(t, 1, metrics.MinWorkers())
}

func TestMetrics_MaxWorkers(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, 0, metrics.MaxWorkers())

	metrics.wp.maxWorkers = 1
	assert.Equal(t, 1, metrics.MaxWorkers())
}

func TestMetrics_QueueCapacity(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}

	metrics.wp.queue = make(chan func())
	assert.Equal(t, 0, metrics.QueueCapacity())

	metrics.wp.queue = make(chan func(), 1)
	assert.Equal(t, 1, metrics.QueueCapacity())
}

func TestMetrics_RunningWorkers(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, 0, metrics.RunningWorkers())

	metrics.wp.runningWorkers = 1
	assert.Equal(t, 1, metrics.RunningWorkers())
}

func TestMetrics_IdleWorkers(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, 0, metrics.IdleWorkers())

	metrics.wp.idleWorkers = 1
	assert.Equal(t, 1, metrics.IdleWorkers())
}

func TestMetrics_SubmittedTasks(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, uint64(0), metrics.SubmittedTasks())

	metrics.wp.submittedTasks = 1
	assert.Equal(t, uint64(1), metrics.SubmittedTasks())
}

func TestMetrics_WaitingTasks(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, uint64(0), metrics.WaitingTasks())

	metrics.wp.waitingTasks = 1
	assert.Equal(t, uint64(1), metrics.WaitingTasks())
}

func TestMetrics_SuccessfulTasks(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, uint64(0), metrics.SuccessfulTasks())

	metrics.wp.successfulTasks = 1
	assert.Equal(t, uint64(1), metrics.SuccessfulTasks())
}

func TestMetrics_FailedTasks(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, uint64(0), metrics.FailedTasks())

	metrics.wp.failedTasks = 1
	assert.Equal(t, uint64(1), metrics.FailedTasks())
}

func TestMetrics_CompletedTasks(t *testing.T) {
	metrics := &Metrics{wp: &WorkerPool{}}
	assert.Equal(t, uint64(0), metrics.CompletedTasks())

	metrics.wp.successfulTasks = 1
	assert.Equal(t, uint64(1), metrics.CompletedTasks())

	metrics.wp.failedTasks = 1
	assert.Equal(t, uint64(2), metrics.CompletedTasks())
}

func TestMetrics_Snapshot(t *testing.T) {
	metrics := &Metrics{
		wp: &WorkerPool{
			minWorkers:      1,
			maxWorkers:      2,
			queue:           make(chan func(), 3),
			runningWorkers:  4,
			idleWorkers:     5,
			submittedTasks:  6,
			waitingTasks:    7,
			successfulTasks: 8,
			failedTasks:     9,
		},
	}

	snapshot := metrics.Snapshot()
	assert.Equal(t, 1, snapshot.minWorkers)
	assert.Equal(t, 2, snapshot.maxWorkers)
	assert.Equal(t, 3, snapshot.queueCapacity)
	assert.Equal(t, 4, snapshot.runningWorkers)
	assert.Equal(t, 5, snapshot.idleWorkers)
	assert.Equal(t, uint64(6), snapshot.submittedTasks)
	assert.Equal(t, uint64(7), snapshot.waitingTasks)
	assert.Equal(t, uint64(8), snapshot.successfulTasks)
	assert.Equal(t, uint64(9), snapshot.failedTasks)
}

func TestMetrics_String(t *testing.T) {
	metrics := &Metrics{
		wp: &WorkerPool{
			minWorkers:      0,
			maxWorkers:      0,
			queue:           make(chan func()),
			runningWorkers:  0,
			idleWorkers:     0,
			submittedTasks:  0,
			waitingTasks:    0,
			successfulTasks: 0,
			failedTasks:     0,
		},
	}

	assert.Equal(t, "Metrics [minWorkers=0, maxWorkers=0, queueCapacity=0, runningWorkers=0, idleWorkers=0, submittedTasks=0, waitingTasks=0, successfulTasks=0, failedTasks=0, completedTasks=0]", metrics.String())

	metrics = &Metrics{
		wp: &WorkerPool{
			minWorkers:      1,
			maxWorkers:      2,
			queue:           make(chan func(), 3),
			runningWorkers:  4,
			idleWorkers:     5,
			submittedTasks:  6,
			waitingTasks:    7,
			successfulTasks: 8,
			failedTasks:     9,
		},
	}
	assert.Equal(t, "Metrics [minWorkers=1, maxWorkers=2, queueCapacity=3, runningWorkers=4, idleWorkers=5, submittedTasks=6, waitingTasks=7, successfulTasks=8, failedTasks=9, completedTasks=17]", metrics.String())
}

func TestMetricsSnapshot_MinWorkers(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, 0, snapshot.MinWorkers())

	snapshot.minWorkers = 1
	assert.Equal(t, 1, snapshot.MinWorkers())
}

func TestMetricsSnapshot_MaxWorkers(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, 0, snapshot.MaxWorkers())

	snapshot.maxWorkers = 1
	assert.Equal(t, 1, snapshot.MaxWorkers())
}

func TestMetricsSnapshot_QueueCapacity(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, 0, snapshot.QueueCapacity())

	snapshot.queueCapacity = 1
	assert.Equal(t, 1, snapshot.QueueCapacity())
}

func TestMetricsSnapshot_RunningWorkers(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, 0, snapshot.RunningWorkers())

	snapshot.runningWorkers = 1
	assert.Equal(t, 1, snapshot.RunningWorkers())
}

func TestMetricsSnapshot_IdleWorkers(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, 0, snapshot.IdleWorkers())

	snapshot.idleWorkers = 1
	assert.Equal(t, 1, snapshot.IdleWorkers())
}

func TestMetricsSnapshot_SubmittedTasks(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, uint64(0), snapshot.SubmittedTasks())

	snapshot.submittedTasks = 1
	assert.Equal(t, uint64(1), snapshot.SubmittedTasks())
}

func TestMetricsSnapshot_WaitingTasks(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, uint64(0), snapshot.WaitingTasks())

	snapshot.waitingTasks = 1
	assert.Equal(t, uint64(1), snapshot.WaitingTasks())
}

func TestMetricsSnapshot_SuccessfulTasks(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, uint64(0), snapshot.SuccessfulTasks())

	snapshot.successfulTasks = 1
	assert.Equal(t, uint64(1), snapshot.SuccessfulTasks())
}

func TestMetricsSnapshot_FailedTasks(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, uint64(0), snapshot.FailedTasks())

	snapshot.failedTasks = 1
	assert.Equal(t, uint64(1), snapshot.FailedTasks())
}

func TestMetricsSnapshot_CompletedTasks(t *testing.T) {
	snapshot := &MetricsSnapshot{}
	assert.Equal(t, uint64(0), snapshot.CompletedTasks())

	snapshot.successfulTasks = 1
	assert.Equal(t, uint64(1), snapshot.CompletedTasks())

	snapshot.failedTasks = 1
	assert.Equal(t, uint64(2), snapshot.CompletedTasks())
}

func TestMetricsSnapshot_String(t *testing.T) {
	snapshot := &MetricsSnapshot{}

	assert.Equal(t, "MetricsSnapshot [minWorkers=0, maxWorkers=0, queueCapacity=0, runningWorkers=0, idleWorkers=0, submittedTasks=0, waitingTasks=0, successfulTasks=0, failedTasks=0, completedTasks=0]", snapshot.String())

	snapshot = &MetricsSnapshot{
		minWorkers:      1,
		maxWorkers:      2,
		queueCapacity:   3,
		runningWorkers:  4,
		idleWorkers:     5,
		submittedTasks:  6,
		waitingTasks:    7,
		successfulTasks: 8,
		failedTasks:     9,
	}
	assert.Equal(t, "MetricsSnapshot [minWorkers=1, maxWorkers=2, queueCapacity=3, runningWorkers=4, idleWorkers=5, submittedTasks=6, waitingTasks=7, successfulTasks=8, failedTasks=9, completedTasks=17]", snapshot.String())
}
