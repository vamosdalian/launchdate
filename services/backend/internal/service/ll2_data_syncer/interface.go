package ll2datasyncer

import "github.com/vamosdalian/launchdate-backend/internal/util"

type Syncer interface {
	Start()
	Pause()
	Resume()
	Cancel()
	Done() <-chan struct{}
}

type BaseSyncer struct {
	rl       util.RateLimiter
	work     func()
	pauseCh  chan struct{}
	resumeCh chan struct{}
	stopCh   chan struct{}
	done     chan struct{}
}

func NewBaseSyncer(rl util.RateLimiter, work func()) *BaseSyncer {
	return &BaseSyncer{
		rl:       rl,
		work:     work,
		done:     make(chan struct{}),
		pauseCh:  make(chan struct{}, 1),
		resumeCh: make(chan struct{}, 1),
		stopCh:   make(chan struct{}, 1),
	}
}

func (s *BaseSyncer) Start() {
	go s.run()
}

func (s *BaseSyncer) Pause() {
	select {
	case s.pauseCh <- struct{}{}:
	default:
	}
}

func (s *BaseSyncer) Resume() {
	select {
	case s.resumeCh <- struct{}{}:
	default:
	}
}

func (s *BaseSyncer) Cancel() {
	s.requestStop()
}

func (s *BaseSyncer) requestStop() {
	select {
	case s.stopCh <- struct{}{}:
	default:
		// Channel full, safe to ignore as stop is already requested
	}
}

func (s *BaseSyncer) Done() <-chan struct{} {
	return s.done
}

func (s *BaseSyncer) run() {
	defer close(s.done)
	for {
		select {
		case <-s.stopCh:
			return
		case <-s.rl.Wait():
			if s.work != nil {
				s.work()
			}
		case <-s.pauseCh:
			select {
			case <-s.stopCh:
				return
			case <-s.resumeCh:
			}
		}
	}
}
