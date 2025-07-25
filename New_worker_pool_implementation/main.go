package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task abstraction to be used by workers in a Worker Pool
type Task[T any] struct {
	Task_id uint32
	Data    T
	Process func(T)
}

// This is the WorkerPool itself
type WorkerPool[T any] struct {
	tasks   chan Task[T]   // chan to transfer tasks into Worker pool
	wg      sync.WaitGroup //
	workers uint32         // size representation of the number of workers that we have
}

// "Constructor" for WorkerPool, creating a new WorkerPool entity with specified number of workers
func CreateWP[T any](workers_val uint32, bufferSize uint32) *WorkerPool[T] {
	return &WorkerPool[T]{
		// creating a new channel of type Task and assigning it to the tasks Channel of
		tasks: make(chan Task[T], bufferSize),
		// assigning number of workers to worker pool
		workers: workers_val,
	}
}

// Starts Worker pool
func (wp *WorkerPool[T]) StartWP(ctx context.Context) {
	var i uint32
	for i = 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i, ctx)
	}
}

// A routine abstraction for a task to be done by any Worker from Worker Pool
func (wp *WorkerPool[T]) worker(id uint32, ctx context.Context) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.tasks:
			if !ok {
				fmt.Printf("Worker %d shutting down\n", id)
				return
			}
			fmt.Printf("Worker %d processing task %d\n", id, task.Task_id)
			task.Process(task.Data)
		case <-ctx.Done():
			fmt.Printf("Worker %d cancelled via context\n", id)
			return
		}
	}
}

// Adds task to execution queue
func (wp *WorkerPool[T]) AddTaskWP(task Task[T]) {
	wp.tasks <- task
}

// Closes tasks chanel and awaits all tasks completion by workers
func (wp *WorkerPool[T]) CloseWP() {
	close(wp.tasks)
	wp.wg.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool := CreateWP[string](5, 100)
	pool.StartWP(ctx)

	var i uint32
	for i = 1; i <= 100; i++ {
		task := Task[string]{
			Task_id: i,
			Data:    fmt.Sprintf("Task data %d", i),
			Process: func(data string) {
				fmt.Printf("Processing data: %s\n", data)
				// THIS IS IMITATION OF WORK
				time.Sleep(10 * time.Millisecond)
			},
		}
		pool.AddTaskWP(task)
	}

	select {
	case <-ctx.Done():
		fmt.Println("Main: context timeout or cancellation")
	case <-time.After(3 * time.Second):
		fmt.Println("Main: tasks completed")
	}

	pool.CloseWP()
}
