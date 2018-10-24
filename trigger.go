package main

import "time"

var past = time.Now().Add(-100 * time.Minute)

type Trigger interface {
	NextFire() time.Time
}

type Never struct {
}

func (trigger Never) NextFire() time.Time {
	return past
}

type OneTime struct {
	t time.Time
}

func (trigger OneTime) NextFire() time.Time {
	now := time.Now()
	if trigger.t.After(now) {
		return trigger.t
	}
	return past
}

type Recurring struct {
	d       time.Duration
	created time.Time
	stopped bool
}

func NewRecurring(d time.Duration) *Recurring {
	return &Recurring{
		d:       d,
		created: time.Now(),
		stopped: false,
	}
}

func (trigger Recurring) NextFire() time.Time {
	if trigger.stopped {
		return past
	}
	now := time.Now()
	diff := time.Since(trigger.created)
	rem := diff % trigger.d
	add := trigger.d - rem
	return now.Add(add)
}

func (trigger *Recurring) Stop() {
	trigger.stopped = true
}
