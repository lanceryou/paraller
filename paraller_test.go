package paraller

import (
	"context"
	"errors"
	"testing"
	"time"
)

// 测试正常返回
// 延时返回
// 错误返回
func TestAsync(t *testing.T) {
	var m1, m2 string
	sleepErr := errors.New("sleep")
	st := []struct {
		f1  func() error
		f2  func() error
		ctx context.Context
		m1  string
		m2  string
		err error
	}{
		{
			f1:  func() error { m1 = "normal f1"; return nil },
			f2:  func() error { m2 = "normal f2"; return nil },
			ctx: context.TODO(),
			m1:  "normal f1",
			m2:  "normal f2",
		},
		{
			f1:  func() error { m1 = "latency f1"; time.Sleep(time.Second); return nil },
			f2:  func() error { m2 = "normal f2"; return nil },
			ctx: context.TODO(),
			m1:  "latency f1",
			m2:  "normal f2",
		},
		{
			f1:  func() error { m1 = "latency f1"; time.Sleep(time.Second); return sleepErr },
			f2:  func() error { m2 = "normal f2"; return nil },
			ctx: context.TODO(),
			m1:  "latency f1",
			m2:  "normal f2",
			err: sleepErr,
		},
		{
			f1:  func() error { m1 = "latency f1"; time.Sleep(time.Second); return sleepErr },
			f2:  func() error { m2 = "normal f2"; time.Sleep(3 * time.Second); return errors.New("long sleep err") },
			ctx: context.TODO(),
			m1:  "latency f1",
			m2:  "normal f2",
			err: sleepErr,
		},
		{
			f1:  func() error { m1 = "latency f1"; time.Sleep(3 * time.Second); return errors.New("long sleep err") },
			f2:  func() error { m2 = "normal f2"; time.Sleep(time.Second); return sleepErr },
			ctx: context.TODO(),
			m1:  "latency f1",
			m2:  "normal f2",
			err: sleepErr,
		},
	}

	for _, s := range st {
		err := Await(s.ctx, Async(s.f1), Async(s.f2))
		if err != s.err {
			t.Errorf("err :%v expect %v", err, s.err)
		}

		if m1 != s.m1 {
			t.Errorf("m1 %v: expect %v", m1, s.m1)
		}

		if m2 != s.m2 {
			t.Errorf("m2 %v: expect %v", m2, s.m2)
		}
	}

	time.Sleep(time.Second * 10)
}
