package fetcher

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testSlowCounter struct {
	i   int32
	max int32

	t *testing.T
}

func (c *testSlowCounter) Get() (string, error) {
	c.slowCall()
	return "", nil
}

func (c *testSlowCounter) List() ([]string, error) {
	c.slowCall()
	return nil, nil
}

func (c *testSlowCounter) slowCall() {
	i := atomic.AddInt32(&c.i, 1)
	if i > c.max {
		c.t.Errorf("Wrong count of fetcher: %d", i)
	} else {
		time.Sleep(20 * time.Millisecond)
	}
	atomic.AddInt32(&c.i, -1)
}

func Test_poolFetcher(t *testing.T) {
	pf := newPoolFetcher(&testSlowCounter{
		max: 2,
		t:   t,
	}, 2)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			_, err := pf.Get()
			if err != nil {
				t.Errorf("Error from fetcher: %v", err)
			}
			wg.Done()
		}(i)

		wg.Add(1)
		go func(i int) {
			_, err := pf.List()
			if err != nil {
				t.Errorf("Error from fetcher: %v", err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
