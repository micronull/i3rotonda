package delayed

import (
	"context"
	"sync"
	"time"

	sw "github.com/micronull/i3rotonda/internal/app/switcher"
	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

type switcher interface {
	Add(ws wm.Workspace)
	Current() wm.Workspace
	Prev() wm.Workspace
	Next() wm.Workspace
}

type Switcher struct {
	sw     switcher
	delay  time.Duration
	cancel context.CancelFunc
	mx     sync.RWMutex
}

func New(poolSize int, delay time.Duration) *Switcher {
	s := &Switcher{
		sw:    sw.New(poolSize),
		delay: delay,
		mx:    sync.RWMutex{},
	}

	return s
}

func (s *Switcher) Add(ws wm.Workspace) {
	if s.cancel != nil {
		s.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())

	s.cancel = cancel

	go func() {
		t := time.NewTimer(s.delay)

		select {
		case <-ctx.Done():
		case <-t.C:
			s.mx.Lock()
			defer s.mx.Unlock()

			s.sw.Add(ws)
		}
	}()
}

func (s *Switcher) Current() wm.Workspace {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.sw.Current()
}

func (s *Switcher) Prev() wm.Workspace {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.sw.Prev()
}

func (s *Switcher) Next() wm.Workspace {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.sw.Next()
}
