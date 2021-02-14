// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package tests

import (
	"testing"
	"time"

	"github.com/fuyibing/cron"
)

func TestStrategy(t *testing.T) {
	t.Logf("---- strategy standard ----")

	s := cron.NewStrategy("35 17 * * *")
	t.Logf("format: %s.", s.Format())
	t.Logf(" weeks: %v.", s.GetWeeks())
	t.Logf("months: %v.", s.GetMonths())
	t.Logf("  days: %v.", s.GetDays())
	t.Logf(" times: %v.", s.GetTimes())

	for i := 0; i < 90; i++ {
		m := time.Now()
		b, _ := s.Validate(m)
		if b {
			s.Update(m)
		}
		t.Logf("[%s]     validate: %v.", m.Format("15:04:05"), b)
		time.Sleep(time.Second)
	}

}

func TestStrategyTimeline(t *testing.T) {
	t.Logf("---- strategy timeline ----")
	s := cron.NewStrategy("20:10,17:52:55")
	t.Logf("strategy:format: %s.", s.Format())
	t.Logf("strategy: times: %v.", s.GetTimes())
	for i := 0; i < 30; i++ {
		m := time.Now()
		b, _ := s.Validate(m)
		if b {
			s.Update(m)
		}
		t.Logf("[%s]     validate: %v.", m.Format("15:04:05"), b)
		time.Sleep(time.Second)
	}
}

func TestStrategyRecycle(t *testing.T) {
	t.Logf("---- strategy recycle ----")
	s := cron.NewStrategy("2s")
	for i := 0; i < 10; i++ {
		m := time.Now()
		b, _ := s.Validate(m)
		if b {
			s.Update(m)
		}
		t.Logf("[%s]     validate: %v.", m.Format("15:04:05.999999"), b)
		time.Sleep(time.Second)
	}
	t.Logf("strategy:format: %s.", s.Format())
}
