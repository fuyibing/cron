// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"testing"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cron/v2"
)

func TestStrategy(t *testing.T) {
	s := cron.NewStrategy("* * * * *")
	// s := cron.NewStrategy("5x")
	log.Infof("format: %s.", s.Format())
	log.Infof(" error: %v.", s.Err())
}
