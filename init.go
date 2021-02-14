// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"sync"
)

var (
	Cron Crontab
)

func init() {
	new(sync.Once).Do(func() {
		Cron = new(crontab)
		Cron.initialize()
	})
}
