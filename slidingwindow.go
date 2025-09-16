package coding_exercises

import (
	"container/list"
	"fmt"
	"time"
)

type Counter interface {
	Increment()
	GetValue() int64
}

type Clock interface {
	NowInSeconds() int64
}

type SystemClock struct{}

func (SystemClock) NowInSeconds() int64 {
	return time.Now().Unix() // Current time in seconds
}

// Sliding Counter

type ListImpl struct {
	clock                 Clock      // Time in seconds
	windowLengthInSeconds int64      // Set window length e.g. 300sec
	buckets               *list.List // linked list of buckets, each one per second
	count                 int64      // total number of clicks in a window
}

type Bucket struct {
	timeSeconds int64 // time label in seconds
	value       int64 // clicks happened in that second
}

func NewListImpl(clock Clock, windowLengthInSeconds int64) *ListImpl {
	return &ListImpl{
		clock:                 clock,
		windowLengthInSeconds: windowLengthInSeconds,
		buckets:               list.New(),
	}
}

func (c *ListImpl) Increment() {
	t := c.clock.NowInSeconds()
	if c.buckets.Len() > 0 {
		last := c.buckets.Back().Value.(*Bucket)
		if last.timeSeconds == t {
			last.value++
			c.count++
			return
		}
	}

	// new bucket if no buckets present
	b := &Bucket{timeSeconds: t, value: 1}
	c.buckets.PushBack(b)
	c.count++
	c.gc(t) // Garbage Collector below
}

func (c *ListImpl) GetValue() int64 {
	c.gc(c.clock.NowInSeconds())
	return c.count
}

func (c *ListImpl) gc(t int64) {
	cutOff := t - c.windowLengthInSeconds
	for c.buckets.Len() > 0 {
		first := c.buckets.Front().Value.(*Bucket)
		if first.timeSeconds <= cutOff {
			c.count -= first.value
			c.buckets.Remove(c.buckets.Front())
		} else {
			break
		}
	}
}

func main() {
	counter := NewListImpl(SystemClock{}, 300) // 5 minutes = 300 secs

	// Call Increment twice
	counter.Increment()
	counter.Increment()
	fmt.Println("Count:", counter.GetValue())

	time.Sleep(2 * time.Second) // new second
	counter.Increment()         // another click creates new bucket, 2nd second
	fmt.Println("Count:", counter.GetValue())

}
