package paraller

import (
	"context"
)

type Paraller struct {
	fn      func() error
	errChan chan error
}

func (p *Paraller) Await() error {
	err := <-p.errChan
	return err
}

func (p *Paraller) run() {
	go func() {
		p.errChan <- p.fn()
	}()
}

func newParaller(fn func() error) *Paraller {
	paraller := &Paraller{
		fn:      fn,
		errChan: make(chan error, 1),
	}

	paraller.run()
	return paraller
}

// 异步的调用，返回一个Paraller对象
func Async(fn func() error) *Paraller {
	return newParaller(fn)
}

// 阻塞等待所有异步调用
func Await(ctx context.Context, ps ...*Paraller) error {
	errChan := make(chan error)
	for _, p := range ps {
		go func(p *Paraller) {
			errChan <- p.Await()
		}(p)
	}

	var cnt int
	cancelFn := func() {
		if cnt == len(ps) {
			return
		}
		for range errChan {
			cnt++
			if cnt == len(ps) {
				break
			}
		}
	}
	for {
		select {
		case <-ctx.Done():
			go cancelFn()
			return ctx.Err()
		case err := <-errChan:
			cnt++
			if err != nil {
				go cancelFn()
				return err
			}

			if cnt == len(ps) {
				return nil
			}
		}
	}
}
