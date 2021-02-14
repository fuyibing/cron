// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/fuyibing/log"

	"github.com/fuyibing/cron"
)

func init() {
	log.Config.TimeFormat = "15:04:05.999999"
	log.Logger.SetAdapter(log.AdapterTerm)
	log.Logger.SetLevel(log.LevelDebug)
}

func TestTicker(t *testing.T) {

	ctx := log.NewContext()
	name := "test-1"
	log.Infofc(ctx, "[ticker=%s] test 1 ticker begin.", name)

	t1 := cron.NewTicker(name, "3s", ticker1)
	t1.SingleNode(true)

	if err := t1.Run(time.Now()); err != nil {
		t.Errorf("Ticker run error: %v.", err)
		return
	}

	t.Logf("Ticker run end.")

	time.Sleep(time.Second)
}

func ticker1(ctx context.Context, ticker cron.Ticker) error {
	// return errors.New("test 1 ticker error")
	log.Debugfc(ctx, "[ticker=%s] run use handler.", ticker.Name())

	return nil
}
