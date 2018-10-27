package main

import (
	"log"
	"time"
)

func main() {
	s := NewScheduler()
	t1 := OneTimeTrigger{time.Now().Add(5 * time.Second)}
	t2 := OneTimeTrigger{time.Now().Add(10 * time.Second)}
	t3 := NewRecurringTrigger(1 * time.Second)
	j1 := LogJob{"job1 executing"}
	j2 := LogJob{"job2 executing"}
	j3 := LogJob{"job3 executing"}
	s.Schedule(t1, j1)
	s.Schedule(t2, j2)
	s.Schedule(t3, j3)

	go func() {
		time.Sleep(20 * time.Second)
		t3.Stop()
		log.Println("t3 stopped")
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("finished running the scheduler")
}
