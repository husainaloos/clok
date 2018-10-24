package main

import "log"

type Job interface {
	Execute() error
}

type LogJob struct {
	message string
}

func (job LogJob) Execute() error {
	log.Println(job.message)
	return nil
}

type NoopJob struct {
}

func (job NoopJob) Execute() error { return nil }
