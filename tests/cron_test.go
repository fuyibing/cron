// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"context"
	"testing"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cron/v2"
)

func TestCrontab(t *testing.T) {
	x := cron.NewCrontab()
	x.AddTicker(
		cron.NewTicker("t1", "5s", ticker),
		cron.NewTicker("t2", "10s", ticker),
	)
	if err := x.Start(); err != nil {
		log.Errorf("error: %s.", err)
		return
	}

	for {

	}

}

func ticker(ctx context.Context, ticker cron.TickerInterface) error {
	log.Infofc(ctx, "[ticker=%s] ticker callback.", ticker.Name())
	return nil
}

