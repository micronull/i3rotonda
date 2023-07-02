package switcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/micronull/i3rotonda/internal/app/switcher"
)

type wm string

func (w wm) GetName() string {
	return string(w)
}

func TestSwitcher_Prev_Empty(t *testing.T) {
	t.Parallel()

	s := switcher.Switcher{}

	assert.Nil(t, s.Prev())
}

func TestSwitcher_Add(t *testing.T) {
	t.Parallel()

	s := switcher.NewSwitcher(3)

	s.Add(wm("1"))

	assert.Equal(t, "1", s.Current().GetName())
}

func TestSwitcher_Add_Current(t *testing.T) {
	t.Parallel()

	const size = 3

	s := switcher.NewSwitcher(size)

	s.Add(wm("1"))
	s.Add(wm("2"))
	s.Add(wm("3"))
	s.Add(wm("4"))
	s.Add(wm("5"))

	assert.Equal(t, "5", s.Current().GetName())
}

func TestSwitcher_Prev(t *testing.T) {
	t.Parallel()

	const size = 10

	s := switcher.NewSwitcher(size)

	s.Add(wm("1"))
	s.Add(wm("2"))
	s.Add(wm("3"))
	s.Add(wm("4"))
	s.Add(wm("5"))

	assert.Equal(t, "4", s.Prev().GetName())
	assert.Equal(t, "3", s.Prev().GetName())
	assert.Equal(t, "2", s.Prev().GetName())
	assert.Equal(t, "1", s.Prev().GetName())
	assert.Equal(t, "5", s.Prev().GetName())
}

func TestSwitcher_Prev_SmallSize(t *testing.T) {
	t.Parallel()

	const size = 3

	s := switcher.NewSwitcher(size)

	s.Add(wm("1"))
	s.Add(wm("2"))
	s.Add(wm("3"))
	s.Add(wm("4"))
	s.Add(wm("5"))

	assert.Equal(t, "4", s.Prev().GetName())
	assert.Equal(t, "3", s.Prev().GetName())
	assert.Equal(t, "5", s.Prev().GetName())
	assert.Equal(t, "4", s.Prev().GetName())
	assert.Equal(t, "3", s.Prev().GetName())
}

func TestSwitcher_Add_SkipIfCurrent(t *testing.T) {
	t.Parallel()

	const size = 3

	s := switcher.NewSwitcher(size)

	s.Add(wm("1"))
	s.Add(wm("2"))
	s.Add(wm("2"))
	s.Add(wm("2"))

	assert.Equal(t, "2", s.Current().GetName())
	assert.Equal(t, "1", s.Prev().GetName())
}
