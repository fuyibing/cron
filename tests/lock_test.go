// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/fuyibing/cache"
)

func TestLocks(t *testing.T) {

	for i := 0; i < 10; i++ {
		go func(n int) {
			l := cache.NewLock(fmt.Sprintf("test:%d", n))
			l.NotRenewal(nil)
			_, err := l.Set(nil);
			if err != nil {
				return
			}

			// defer l.Unset(nil, s)

			time.Sleep(time.Second)

		}(i)
	}

	time.Sleep(time.Second * 15)

	// x := cron.NewCrontab()
	// x.AddTicker(cron.NewTicker("t1", "5s", nil), cron.NewTicker("t2", "1m", nil))
	// if err := x.Start(); err != nil {
	// 	log.Errorf("error: %s.", err)
	// }
}
