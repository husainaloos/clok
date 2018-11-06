package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	// ErrAlreadyRunning indicates that the scheduler is already running
	ErrAlreadyRunning = errors.New("scheduler is already running")

	// ErrNotRunning indicates that the scheduler is not running
	ErrNotRunning = errors.New("scheduler is not running")
)

type tjPair struct {
	t Trigger
	j Job
}

// Scheduler is the core object that coordinates executing jobs
type Scheduler struct {
	pairs []tjPair
	wg    sync.WaitGroup
	done  chan (struct{})

	mu      sync.Mutex
	running bool
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		pairs: make([]tjPair, 0),
		done:  make(chan struct{}),
	}
}

// Schedule adds a job to the list of jobs to be executed.
// NOTE: Schedule would not schedule a job if the scheduler is already running
func (sch *Scheduler) Schedule(trigger Trigger, job Job) {
	p := tjPair{trigger, job}
	sch.pairs = append(sch.pairs, p)
}

func (sch *Scheduler) startTriggeredJob(p tjPair) {
	sch.wg.Add(1)
	go func() {
		defer sch.wg.Done()
		for {
			ft := p.t.NextFire()
			if !ft.After(time.Now()) {
				return
			}
			d := time.Until(ft)
			select {
			case <-sch.done:
				return
			case <-time.After(d):
				go func() {
					if err := p.j.Execute(); err != nil {
						log.Println(err)
					}
				}()
			}

		}
	}()
}

// Start starts the scheduler and executing jobs according to their triggers
// Start is blocking function.
func (sch *Scheduler) Start() error {
	sch.mu.Lock()
	if sch.running {
		sch.mu.Unlock()
		return ErrAlreadyRunning
	}
	sch.running = true
	sch.mu.Unlock()
	for _, p := range sch.pairs {
		sch.startTriggeredJob(p)
	}
	sch.wg.Wait()
	return nil
}

// Stop stops the scheduler
func (sch *Scheduler) Stop() {
	sch.mu.Lock()
	defer sch.mu.Unlock()
	if !sch.running {
		return
	}
	for range sch.pairs {
		sch.done <- struct{}{}
	}
	return
}
