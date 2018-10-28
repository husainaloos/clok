package main

import "time"

var past = time.Now().Add(-100 * time.Minute)

// Trigger is an interface for anything that define when a job to be executed
type Trigger interface {
	// NextFire should return the time when it will fire next. If it no longer fires, it should return a past time
	NextFire() time.Time
}

// NeverTrigger is the nil triggers that never triggers
type NeverTrigger struct {
}

// NextFire always returns a past value
func (trigger NeverTrigger) NextFire() time.Time {
	return past
}

// OneTimeTrigger is a trigger that fires only once
type OneTimeTrigger struct {
	t time.Time
}

// NextFire return the time it will fire
func (trigger OneTimeTrigger) NextFire() time.Time {
	return trigger.t
}

// RecurringTrigger is a trigger that fires every certain duration
type RecurringTrigger struct {
	d       time.Duration
	created time.Time
	stopped bool
}

// NewRecurringTrigger creates a new recurring trigger
func NewRecurringTrigger(d time.Duration) *RecurringTrigger {
	return &RecurringTrigger{
		d:       d,
		created: time.Now(),
		stopped: false,
	}
}

// NextFire returns a time when it will fire by a duration. This time starts since the creation of the trigger
func (trigger RecurringTrigger) NextFire() time.Time {
	if trigger.stopped {
		return past
	}
	now := time.Now()
	diff := time.Since(trigger.created)
	rem := diff % trigger.d
	add := trigger.d - rem
	return now.Add(add)
}

// Stop stops the trigger from firing
func (trigger *RecurringTrigger) Stop() {
	trigger.stopped = true
}
