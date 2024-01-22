package delayed_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/micronull/i3rotonda/internal/app/delayed"
)

type wm string

func (w wm) Name() string {
	return string(w)
}

func (wm) IsEmpty() bool {
	return false
}

func TestSwitcher_Prev(t *testing.T) {
	t.Parallel()

	const (
		size  = 3
		delay = time.Millisecond * 500
	)

	s := delayed.New(size, delay)

	s.Add(wm("1"))

	time.Sleep(delay + 100) // wait save "1"

	s.Add(wm("2"))
	s.Add(wm("3"))

	time.Sleep(delay / 2)

	s.Add(wm("4"))

	time.Sleep(delay) // wait save last - "4"

	assert.Equal(t, "1", s.Prev().Name())
	assert.Equal(t, "4", s.Prev().Name())
	assert.Equal(t, "1", s.Prev().Name())
	assert.Equal(t, "4", s.Prev().Name())
}
