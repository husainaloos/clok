package main

import "log"

// Job represent a job that should be executed
type Job interface {
	Execute() error
}

// LogJob is a job that logs a message. This job is typically used for testing
type LogJob struct {
	message string
}

// Execute logs a message
func (job LogJob) Execute() error {
	log.Println(job.message)
	return nil
}

// NoopJob is the nil job and it does nothing
type NoopJob struct{}

// Execute does nothing
func (job NoopJob) Execute() error { return nil }
