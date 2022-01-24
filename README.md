# Golang WorkerPool

![Build status](https://github.com/leocarmona/workerpool/actions/workflows/main.yaml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/leocarmona/workerpool)
<a title="Release" target="_blank" href="https://github.com/leocarmona/workerpool/releases"><img src="https://img.shields.io/github/v/release/leocarmona/workerpool"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/leocarmona/workerpool"><img src="https://goreportcard.com/badge/github.com/leocarmona/workerpool"/></a>
[![GoDoc](https://pkg.go.dev/badge/github.com/leocarmona/workerpool)](https://pkg.go.dev/github.com/leocarmona/workerpool)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/leocarmona/workerpool/blob/main/LICENSE)

## Installation
To install this library in your project, do you need to run:
```
go get github.com/leocarmona/workerpool@1.0.0
```

## How to use

### All commands available
```go
package main

import (
   "fmt"
   "github.com/leocarmona/workerpool"
   "runtime/debug"
   "time"
)

func main() {
   wp := workerpool.New("io", &workerpool.Settings{
      MinWorkers:  64,
      MaxWorkers:  256,
      IdleTimeout: 15 * time.Second,
      UpScaling:   4,
      DownScaling: 2,
      Queue:       4096,
      PanicHandler: func(panicErr interface{}) {
         fmt.Printf("Worker recovered from a panic: %v\nStack trace: %s\n", panicErr, string(debug.Stack()))
      },
   })

   task := func() { fmt.Println("hi") }
   timeout := time.Second
   deadline := time.Now().Add(timeout)

   // Executions
   wp.Submit(task)
   wp.SubmitAndWait(task)
   wp.SubmitAndWaitWithTimeout(timeout, task)
   wp.SubmitAndWaitWithDeadline(deadline, task)

   wp.TrySubmit(task)
   wp.TrySubmitAndWait(task)
   wp.TrySubmitAndWaitWithTimeout(timeout, task)
   wp.TrySubmitAndWaitWithDeadline(deadline, task)

   // Manage Workers/Goroutines
   wp.Burst(1)
   wp.ScaleUp(1)
   wp.ScaleDown(1)
   wp.ReleaseIdleWorkers()

   wp.Stop()
   wp.StopAndWait()
   wp.StopAndWaitWithTimeout(timeout)
   wp.StopAndWaitWithDeadline(deadline)

   wp.Stopped()

   // Metrics
   wp.Metrics().RunningWorkers()
   wp.Metrics().IdleWorkers()
   wp.Metrics().MinWorkers()
   wp.Metrics().MaxWorkers()
   wp.Metrics().QueueCapacity()
   wp.Metrics().SubmittedTasks()
   wp.Metrics().WaitingTasks()
   wp.Metrics().SuccessfulTasks()
   wp.Metrics().FailedTasks()
   wp.Metrics().CompletedTasks()
   wp.Metrics().Snapshot()
   wp.Metrics().String()

   // Settings
   wp.MinWorkers()
   wp.MaxWorkers()
   wp.IdleTimeout()
   wp.UpScaling()
   wp.DownScaling()
   wp.QueueCapacity()
   wp.PanicHandler()

   wp.SetMinWorkers(0)
   wp.SetMaxWorkers(1)
   wp.SetIdleTimeout(15 * time.Second)
   wp.SetUpScaling(1)
   wp.SetDownScaling(1)
   wp.SetPanicHandler(func(panicErr interface{}) {
      fmt.Printf("Worker recovered from a panic: %v\nStack trace: %s\n", panicErr, string(debug.Stack()))
   })

   wp.String()
}
```

## API Reference

Full API reference is available at https://pkg.go.dev/github.com/leocarmona/workerpool
