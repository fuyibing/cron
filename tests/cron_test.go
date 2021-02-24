// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"testing"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cron/v2"
)

func TestCrontab(t *testing.T) {
	x := cron.NewCrontab()
	x.AddTicker(cron.NewTicker("t1", "5s", nil), cron.NewTicker("t2", "1m", nil))
	if err := x.Start(); err != nil {
		log.Errorf("error: %s.", err)
	}
}
