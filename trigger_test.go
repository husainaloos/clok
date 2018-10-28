package main

import (
	"testing"
	"time"
)

func Test_OneTimeTrigger(t *testing.T) {
	t1 := time.Now().Add(1 * time.Second)
	t2 := time.Now().Add(2 * time.Second)
	t3 := time.Now().Add(3 * time.Second)
	t4 := time.Now().Add(100 * time.Millisecond)
	t5 := time.Now().Add(-1 * time.Second)

	tcs := []struct {
		id     int
		input  time.Time
		output time.Time
	}{
		{1, t1, t1},
		{2, t2, t2},
		{3, t3, t3},
		{4, t4, t4},
		{5, t5, t5},
	}

	for _, tc := range tcs {
		trg := OneTimeTrigger{tc.input}
		got := trg.NextFire()
		if got != tc.output {
			t.Errorf("OneTimeTrigger.NextFire() test %d, expected=%v, found=%v", tc.id, tc.output, got)
		}
	}
}

func Test_NeverTrigger(t *testing.T) {
	trg := NeverTrigger{}
	nf := trg.NextFire()
	if nf.After(time.Now()) {
		t.Errorf("NeverTrigger.NextFire(), expected=past, found=%v(present of future)", nf)
	}
}

func Test_RecurringTrigger(t *testing.T) {
	tcs := []struct {
		d time.Duration
	}{
		{100 * time.Millisecond},
		{10 * time.Millisecond},
		{1 * time.Millisecond},
	}

	for _, tc := range tcs {
		d := tc.d
		delta := 1 * time.Nanosecond
		trg := NewRecurringTrigger(d)
		nf := trg.NextFire()
		dur := nf.Sub(time.Now())
		if dur > d+delta || dur < 0 {
			t.Errorf("RecurringTrigger.NextFire(), expected=%v, got=%v", time.Now().Add(d), nf)
		}

		time.Sleep(d)
		nf = trg.NextFire()
		dur = nf.Sub(time.Now())
		if dur > d+delta || dur < 0 {
			t.Errorf("RecurringTrigger.NextFire(), expected=%v, got=%v", time.Now().Add(d), nf)
		}

	}
}
