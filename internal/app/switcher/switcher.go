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

	for c := len(s.pool); c > 0; c-- {
		if s.counter == 0 {
			s.counter = len(s.pool) - 1
		} else {
			s.counter--
		}

		ws := s.pool[s.counter]

		if ws.IsEmpty() {
			continue
		}

		return ws
	}

	return nil
}

func (s *Switcher) Next() wm.Workspace {
	if len(s.pool) == 0 {
		return nil
	}

	ct := s.counter

	for c := len(s.pool); c > 0; c-- {
		if ct == len(s.pool)-1 {
			ct = 0
		} else {
			ct++
		}

		ws := s.pool[ct]

		if ws.IsEmpty() {
			continue
		}

		if s.isCurrent(ws) {
			continue
		}

		s.counter = ct

		return ws
	}

	return nil
}
