// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"errors"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
)

// Crontab struct.
type crontab struct {
	mu      *sync.RWMutex
	tickers map[string]TickerInterface
}

// Add ticker.
func (o *crontab) AddTicker(ts ...TickerInterface) CrontabInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, ti := range ts {
		o.tickers[ti.Name()] = ti
	}
	return o
}

// Start crontab.
func (o *crontab) Start() error {
	if len(o.tickers) == 0 {
		return errors.New("no ticker found")
	}
	o.listen()
	return nil
}

// Dispatch ticker.
func (o *crontab) dispatcher() {
	o.mu.RLock()
	defer o.mu.RUnlock()
	t := time.Now()
	for _, x := range o.tickers {
		go x.Run(t)
	}
}

// Listen second.
func (o *crontab) listen() {
	go func() {
		defer o.listen()
		log.Debugf("[crontab] crontab listening.")
		for range time.NewTicker(time.Second).C {
			go o.dispatcher()
		}
	}()
}

// New crontab instance.
func NewCrontab() CrontabInterface {
	o := &crontab{}
	o.mu = new(sync.RWMutex)
	o.tickers = make(map[string]TickerInterface)
	return o
}
