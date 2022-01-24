package workerpool

import (
	"fmt"
	"sync/atomic"
)

type Metrics struct {
	wp *WorkerPool
}

func (m *Metrics) MinWorkers() int {
	return m.wp.minWorkers
}

func (m *Metrics) MaxWorkers() int {
	return m.wp.maxWorkers
}

func (m *Metrics) QueueCapacity() int {
	return cap(m.wp.queue)
}

func (m *Metrics) RunningWorkers() int {
	return int(atomic.LoadInt32(&m.wp.runningWorkers))
}

func (m *Metrics) IdleWorkers() int {
	return int(atomic.LoadInt32(&m.wp.idleWorkers))
}

func (m *Metrics) SubmittedTasks() uint64 {
	return atomic.LoadUint64(&m.wp.submittedTasks)
}

func (m *Metrics) WaitingTasks() uint64 {
	return atomic.LoadUint64(&m.wp.waitingTasks)
}

func (m *Metrics) SuccessfulTasks() uint64 {
	return atomic.LoadUint64(&m.wp.successfulTasks)
}

func (m *Metrics) FailedTasks() uint64 {
	return atomic.LoadUint64(&m.wp.failedTasks)
}

func (m *Metrics) CompletedTasks() uint64 {
	return m.SuccessfulTasks() + m.FailedTasks()
}

func (m *Metrics) Snapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		minWorkers:      m.MinWorkers(),
		maxWorkers:      m.MaxWorkers(),
		queueCapacity:   m.QueueCapacity(),
		runningWorkers:  m.RunningWorkers(),
		idleWorkers:     m.IdleWorkers(),
		submittedTasks:  m.SubmittedTasks(),
		waitingTasks:    m.WaitingTasks(),
		successfulTasks: m.SuccessfulTasks(),
		failedTasks:     m.FailedTasks(),
	}
}

func (m *Metrics) String() string {
	return fmt.Sprintf(metricsStringFormat, m.MinWorkers(), m.MaxWorkers(), m.QueueCapacity(), m.RunningWorkers(), m.IdleWorkers(), m.SubmittedTasks(), m.WaitingTasks(), m.SuccessfulTasks(), m.FailedTasks(), m.CompletedTasks())
}

const metricsStringFormat = "Metrics [minWorkers=%d, maxWorkers=%d, queueCapacity=%d, runningWorkers=%d, idleWorkers=%d, submittedTasks=%d, waitingTasks=%d, successfulTasks=%d, failedTasks=%d, completedTasks=%d]"

type MetricsSnapshot struct {
	minWorkers      int
	maxWorkers      int
	queueCapacity   int
	runningWorkers  int
	idleWorkers     int
	submittedTasks  uint64
	waitingTasks    uint64
	successfulTasks uint64
	failedTasks     uint64
}

func (m *MetricsSnapshot) MinWorkers() int {
	return m.minWorkers
}

func (m *MetricsSnapshot) MaxWorkers() int {
	return m.maxWorkers
}

func (m *MetricsSnapshot) QueueCapacity() int {
	return m.queueCapacity
}

func (m *MetricsSnapshot) RunningWorkers() int {
	return m.runningWorkers
}

func (m *MetricsSnapshot) IdleWorkers() int {
	return m.idleWorkers
}

func (m *MetricsSnapshot) SubmittedTasks() uint64 {
	return m.submittedTasks
}

func (m *MetricsSnapshot) WaitingTasks() uint64 {
	return m.waitingTasks
}

func (m *MetricsSnapshot) SuccessfulTasks() uint64 {
	return m.successfulTasks
}

func (m *MetricsSnapshot) FailedTasks() uint64 {
	return m.failedTasks
}

func (m *MetricsSnapshot) CompletedTasks() uint64 {
	return m.SuccessfulTasks() + m.FailedTasks()
}

func (m *MetricsSnapshot) String() string {
	return fmt.Sprintf(metricsSnapshotStringFormat, m.MinWorkers(), m.MaxWorkers(), m.QueueCapacity(), m.RunningWorkers(), m.IdleWorkers(), m.SubmittedTasks(), m.WaitingTasks(), m.SuccessfulTasks(), m.FailedTasks(), m.CompletedTasks())
}

const metricsSnapshotStringFormat = "MetricsSnapshot [minWorkers=%d, maxWorkers=%d, queueCapacity=%d, runningWorkers=%d, idleWorkers=%d, submittedTasks=%d, waitingTasks=%d, successfulTasks=%d, failedTasks=%d, completedTasks=%d]"
