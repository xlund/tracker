package domain

import (
	"testing"
)

func TestGameClockString(t *testing.T) {
	gc := GameClock{Initial: 3800, Increment: 30}
	s := gc.Readable()
	println(s)
}
