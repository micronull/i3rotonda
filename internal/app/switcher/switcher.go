// Package switcher contains switching logic.
package switcher

import (
	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

type Switcher struct {
	pool    []wm.Workspace
	counter int
}

func NewSwitcher(poolSize int) *Switcher {
	return &Switcher{
		pool: make([]wm.Workspace, 0, poolSize),
	}
}

func (s *Switcher) Add(w wm.Workspace) {
	if s.isCurrent(w) {
		return
	}

	if len(s.pool) >= s.counter+1 && s.pool[s.counter].GetName() == w.GetName() {
		return
	}

	if len(s.pool) == cap(s.pool) {
		for i := 0; i < len(s.pool); i++ {
			n := i + 1

			if n < len(s.pool) {
				s.pool[i] = s.pool[n]
			} else {
				s.pool[i] = w
			}
		}
	} else {
		s.pool = append(s.pool, w)
	}

	s.counter = len(s.pool) - 1
}

func (s *Switcher) isCurrent(w wm.Workspace) bool {
	c := s.Current()

	if c == nil {
		return false
	}

	return w.GetName() == c.GetName()
}

func (s *Switcher) Current() wm.Workspace {
	if len(s.pool) == 0 {
		return nil
	}

	return s.pool[s.counter]
}

func (s *Switcher) Prev() wm.Workspace {
	if len(s.pool) == 0 {
		return nil
	}

	if s.counter == 0 {
		s.counter = len(s.pool) - 1
	} else {
		s.counter--
	}

	return s.pool[s.counter]
}

func (s *Switcher) Next() wm.Workspace {
	if len(s.pool) == 0 {
		return nil
	}

	if s.counter+1 == len(s.pool) {
		s.counter = 0
	} else {
		s.counter++
	}

	return s.pool[s.counter]
}
