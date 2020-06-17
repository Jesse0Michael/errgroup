package errgroup

import (
	"context"
	"fmt"
	"sync"
)

// SizedErrGroup is used to work through a queue of work using an errgroup
// with limits on the amount of work that will be done concurrently
type SizedErrGroup struct {
	current chan struct{}
	ctx     context.Context
	cancel  func()
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// NewSizedErrGroup creates a SizedErrGroup with the set size of concurrent work
func NewSizedErrGroup(size int) *SizedErrGroup {
	seg, _ := WithContext(context.Background(), size)
	return seg
}

// WithContext returns a new Group and an associated Context derived from ctx.
// The derived Context is canceled the first time a function passed to Go returns an error
func WithContext(ctx context.Context, size int) (*SizedErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	if size < 1 {
		size = 1
	}
	group := &SizedErrGroup{
		current: make(chan struct{}, size),
		ctx:     ctx,
		cancel:  cancel,
	}
	return group, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *SizedErrGroup) Wait() error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		g.wg.Wait()
	}()

	select {
	case <-c:
		return g.err
	case <-g.ctx.Done():
		if g.err != nil {
			return g.err
		}
		return fmt.Errorf("wait group context cancelled")
	}
}

// done decrements the SizedWaitGroup counter.
// See sync.WaitGroup documentation for more information.
func (g *SizedErrGroup) done() {
	<-g.current
	g.wg.Done()
}

// Go calls the given function in a new goroutine.
// The first call to return a non-nil error cancels the group; its error will be returned by Wait.
func (g *SizedErrGroup) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		select {
		case g.current <- struct{}{}:

			go func() {
				defer g.done()

				if err := f(); err != nil {
					g.errOnce.Do(func() {
						g.err = err
						if g.cancel != nil {
							g.cancel()
						}
					})
				}
			}()
			return
		case <-g.ctx.Done():
			return
		}
	}()
}
