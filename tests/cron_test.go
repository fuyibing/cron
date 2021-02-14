// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package tests

import (
	"context"
	"testing"

	"github.com/fuyibing/cron"
)

func TestCron(t *testing.T) {
	c := cron.Cron
	c.Add(
		cron.NewTicker("t1", "3s", f1),
		cron.NewTicker("t2", "19:00:00", f2),
	)

	if err := c.Start(); err != nil {
		t.Errorf("start error: %v.", err)
		return
	}

	for{}

}

func f1(ctx context.Context, ticker cron.Ticker) error {
	return nil
}

func f2(ctx context.Context, ticker cron.Ticker) error {
	return nil
}
