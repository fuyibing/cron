// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cron/v2"
)

func TestStrategy(t *testing.T) {
	// s := cron.NewStrategy("* * * * *")
	s := cron.NewStrategy("5s")
	log.Infof("strategy: %v.", s)
	log.Infof("format: %s.", s.Format())
	log.Infof(" error: %v.", s.Err())
	for i := 0; i <= 12; i++ {
		log.Infof(" valid: %v.", s.Validate(time.Now()))
		time.Sleep(time.Second)
	}
}

func TestStrategyStandard(t *testing.T) {
	s := cron.NewStrategy("*/5 * * * * * *")
	// s := cron.NewStrategy("5x")
	log.Infof("  type: %s.", reflect.TypeOf(s).Elem().String())
	log.Infof("format: %s.", s.Format())
	log.Infof(" error: %v.", s.Err())
	for i := 0; i <= 12; i++ {
		log.Infof(" valid: %v.", s.Validate(time.Now()))
		time.Sleep(time.Second)
	}
}

func TestStrategyTimeline(t *testing.T) {
	// s := cron.NewStrategy("18:00,19:20:35")
	// s := cron.NewStrategy("5x")
	// log.Infof("  type: %s.", reflect.TypeOf(s).Elem().String())
	// log.Infof("format: %s.", s.Format())
	// log.Infof(" error: %v.", s.Err())
	// log.Infof(" valid: %v.", s.Validate(time.Now()))
}
