package main

import (
	"log"
	"sync"
	"time"
)

type scheduledPair struct {
	t  Trigger
	j  Job
	id int
}

type Scheduler struct {
	pairs []scheduledPair
	size  int
	wg    sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		pairs: make([]scheduledPair, 0),
		size:  0,
	}
}

func (scheduler *Scheduler) Schedule(trigger Trigger, job Job) {
	scheduler.pairs = append(scheduler.pairs, scheduledPair{trigger, job, scheduler.size})
	scheduler.size++
}

func (scheduler *Scheduler) startTriggeredJob(p scheduledPair) {
	defer scheduler.wg.Done()
	for {
		ft := p.t.NextFire()
		if !ft.After(time.Now()) {
			return
		}
		d := time.Until(ft)
		time.Sleep(d)
		if err := p.j.Execute(); err != nil {
			log.Println(err)
		}
	}
}

func (scheduler *Scheduler) Start() error {
	for _, p := range scheduler.pairs {
		scheduler.wg.Add(1)
		go scheduler.startTriggeredJob(p)
	}
	scheduler.wg.Wait()
	return nil
}
