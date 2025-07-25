package main

import (
	"log"
	"sync"
	"time"
)

// Task abstraction to be used by workers in a Worker Pool
type Task struct {
	Task_id uint32
}

// This is the WorkerPool itself
type WorkerPool struct {
	tasks   chan Task      // chan to transfer tasks into Worker pool
	wg      sync.WaitGroup //
	workers uint32         // size representation of the number of workers that we have
}

// "Constructor" for WorkerPool, creating a new WorkerPool entity with specified number of workers
func CreateWP(workers_val uint32) *WorkerPool {
	return &WorkerPool{
		// creating a new channel of type Task and assigning it to the tasks Channel of
		tasks: make(chan Task),
		// assigning number of workers to worker pool
		workers: workers_val,
	}
}

// Starts Worker pool
func (wp *WorkerPool) StartWP() {
	var i uint32
	for i = 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// A routine abstraction for a task to be done by any Worker from Worker Pool
func (wp *WorkerPool) worker(id uint32) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		log.Printf("LOG : Worker %d --- Task %d", id, task)
		// TASK GOES HERE, NOW IT IS A SIMULATION
		time.Sleep(10 * time.Millisecond)
	}
}

// Adds task to execution queue
func (wp *WorkerPool) AddTaskWP(task Task) {
	wp.tasks <- task
}

// Closes tasks chanel and awaits all tasks completion by workers
func (wp *WorkerPool) CloseWP() {
	close(wp.tasks)
	wp.wg.Wait()
}

func main() {
	worker_pool := CreateWP(6) // Create worker pool
	worker_pool.StartWP()      // Start worker pool
	var i uint32
	for i = 1; i <= 100; i++ {
		worker_pool.AddTaskWP(Task{Task_id: i}) // Initialise it with tasks
	}
	worker_pool.CloseWP() // End worker pool activity after all tasks are done
}
