package util

import "time"

type RateLimiter interface {
	Allow() bool
	Wait() <-chan struct{}
	Close()
}

type RateLimit struct {
	du time.Duration
	t  *time.Ticker
	c  chan struct{}
}

func NewRateLimit(td time.Duration) RateLimiter {
	rl := &RateLimit{
		du: td,
		t:  time.NewTicker(td),
		c:  make(chan struct{}, 1),
	}
	rl.c <- struct{}{}
	go rl.Start()
	return rl
}

func (r *RateLimit) Start() {
	for range r.t.C {
		select {
		case r.c <- struct{}{}:
		default:
		}
	}
}

func (r *RateLimit) Allow() bool {
	select {
	case <-r.c:
		return true
	default:
		return false
	}
}

func (r *RateLimit) Wait() <-chan struct{} {
	return r.c
}

func (r *RateLimit) Close() {
	r.t.Stop()
}
