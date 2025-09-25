package done

import (
	"testing"
	"time"
)

func delayedClose(d time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		time.Sleep(d)
		close(c)
	}()
	return c
}

func TestOrZeroChannels(t *testing.T) {
	merged := Or()
	select {
	case <-merged:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("should close immediately")
	}
}

func TestOrOneChannelClosed(t *testing.T) {
	c := make(chan interface{})
	close(c)
	merged := Or(c)
	select {
	case <-merged:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("should close immediately")
	}
}

func TestOrOneChannelDelayed(t *testing.T) {
	c := delayedClose(50 * time.Millisecond)
	merged := Or(c)
	select {
	case <-merged:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("should close after delay")
	}
}

func TestOrTwoChannelsFirstCloses(t *testing.T) {
	c1 := make(chan interface{})
	close(c1)
	c2 := delayedClose(1 * time.Second)
	merged := Or(c1, c2)
	select {
	case <-merged:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("should close immediately")
	}
}

func TestOrTwoChannelsSecondCloses(t *testing.T) {
	c1 := delayedClose(1 * time.Second)
	c2 := delayedClose(50 * time.Millisecond)
	merged := Or(c1, c2)
	select {
	case <-merged:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("should close when second does")
	}
}

func TestOrTwoChannelsBothClose(t *testing.T) {
	c1 := delayedClose(50 * time.Millisecond)
	c2 := delayedClose(50 * time.Millisecond)
	merged := Or(c1, c2)
	select {
	case <-merged:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("should close")
	}
}

func TestOrThreeChannelsFirstCloses(t *testing.T) {
	c1 := make(chan interface{})
	close(c1)
	c2 := delayedClose(1 * time.Second)
	c3 := delayedClose(1 * time.Second)
	merged := Or(c1, c2, c3)
	select {
	case <-merged:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("should close immediately")
	}
}

func TestOrThreeChannelsLastCloses(t *testing.T) {
	c1 := delayedClose(1 * time.Second)
	c2 := delayedClose(1 * time.Second)
	c3 := delayedClose(50 * time.Millisecond)
	merged := Or(c1, c2, c3)
	select {
	case <-merged:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("should close when last does")
	}
}

func TestOrFiveChannelsRandomCloses(t *testing.T) {
	c1 := delayedClose(1 * time.Second)
	c2 := delayedClose(1 * time.Second)
	c3 := make(chan interface{})
	close(c3)
	c4 := delayedClose(1 * time.Second)
	c5 := delayedClose(1 * time.Second)
	merged := Or(c1, c2, c3, c4, c5)
	select {
	case <-merged:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("should close immediately")
	}
}

func TestOrNilChannel(t *testing.T) {
	merged := Or((<-chan interface{})(nil))
	select {
	case <-merged:
		t.Fatal("should not close on nil channel")
	case <-time.After(10 * time.Millisecond):
	}
}
