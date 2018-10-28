package main

import "time"

var past = time.Now().Add(-100 * time.Minute)

type Trigger interface {
	NextFire() time.Time
}

type NeverTrigger struct {
}

func (trigger NeverTrigger) NextFire() time.Time {
	return past
}

type OneTimeTrigger struct {
	t time.Time
}

func (trigger OneTimeTrigger) NextFire() time.Time {
	return trigger.t
}

type RecurringTrigger struct {
	d       time.Duration
	created time.Time
	stopped bool
}

func NewRecurringTrigger(d time.Duration) *RecurringTrigger {
	return &RecurringTrigger{
		d:       d,
		created: time.Now(),
		stopped: false,
	}
}

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

func (trigger *RecurringTrigger) Stop() {
	trigger.stopped = true
}
