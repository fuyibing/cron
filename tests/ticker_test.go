// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cron/v2"
)

func TestTicker(t *testing.T) {
	x := cron.NewTicker("t1", "5s", handler)
	x.SingleNode(true)
	log.Infof("name: %s.", x.Name())
	log.Infof("strategy: %s.", x.Strategy().Format())

	x.Run(time.Now())
}

func handler(ctx context.Context, ticker cron.TickerInterface) error {

	for i := 0; i < 10; i++ {
		log.Infofc(ctx, "[ticker=%s] ticker callback.", ticker.Name())
		time.Sleep(time.Second)
	}
	return nil
}
