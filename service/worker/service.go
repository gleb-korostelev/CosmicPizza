// Package worker implements a worker pool that handles tasks concurrently.
// This package is designed to efficiently process tasks using multiple goroutines,
// managing synchronization and lifecycle of worker routines.
package worker

import (
	"context"
	"sync"

	"github.com/gleb-korostelev/CosmicPizza.git/tools/logger"
)

// Task represents a unit of work to be executed by the worker pool.
// It contains an action to be executed and a channel to signal completion of the task.
type Task struct {
	Action func(ctx context.Context) error // Action is the function that performs the task.
	Done   chan struct{}                   // Done is used to signal the completion of the task.
}

// WorkerPool manages a pool of worker goroutines that execute Tasks.
type WorkerPool struct {
	taskQueue  chan Task      // taskQueue is a channel that holds tasks to be processed by the workers.
	wg         sync.WaitGroup // wg is used to wait for all workers to finish processing before shutdown.
	maxWorkers int            // maxWorkers defines the maximum number of worker goroutines.
}

// NewWorkerPool initializes a new WorkerPool with a specified number of workers.
// maxWorkers specifies the maximum number of concurrent workers in the pool.
func NewWorkerPool(maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		taskQueue:  make(chan Task),
		maxWorkers: maxWorkers,
	}

	pool.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

// worker is a goroutine that processes Tasks from the taskQueue.
// It executes the Task's Action and signals completion via the Task's Done channel.
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		// logger.Infof("WorkerPoolService: Task sent to be executed: %+v", task)
		if err := task.Action(context.Background()); err != nil {
			logger.Infof("Error executing task: %v", err)
		}
		if task.Done != nil {
			close(task.Done)
		}
	}
}

// AddTask submits a new Task to the pool. It adds the Task to the taskQueue.
func (p *WorkerPool) AddTask(task Task) {
	p.taskQueue <- task
}

// Shutdown gracefully stops the worker pool. It closes the taskQueue and waits for all workers to finish.
func (p *WorkerPool) Shutdown() {
	close(p.taskQueue)
	p.wg.Wait()
}
