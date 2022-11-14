package timer

import "time"

type LimitTimer struct {
	t    *time.Timer
	wait time.Duration
}

func New(wait time.Duration) *LimitTimer {
	return &LimitTimer{
		t:    time.NewTimer(wait),
		wait: wait,
	}
}

func (l *LimitTimer) Stop() bool {
	return l.t.Stop()
}

func (l *LimitTimer) Reset() bool {
	return l.t.Reset(l.wait)
}
func (l *LimitTimer) ResetWithNewDuration(newWait time.Duration) bool {
	l.wait = newWait
	return l.t.Reset(newWait)
}

func (l *LimitTimer) C() <-chan time.Time {
	return l.t.C
}
