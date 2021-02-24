// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
)

// Crontab interface.
type CrontabInterface interface {
	// Add ticker.
	AddTicker(...TickerInterface) CrontabInterface

	// Delete specified ticker.
	DelTicker(...TickerInterface) CrontabInterface

	// Start crontab.
	Start() error
}

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

// Delete tickers.
func (o *crontab) DelTicker(ts ...TickerInterface) CrontabInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, ti := range ts {
		name := ti.Name()
		if _, ok := o.tickers[name]; ok {
			delete(o.tickers, name)
		}
	}
	return o
}

// Start crontab.
func (o *crontab) Start() error {
	if len(o.tickers) == 0 {
		return ErrNoTicker
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
		log.Debugf("[crontab] listening.")
		for range time.NewTicker(time.Second).C {
			go o.dispatcher()
		}
	}()
}

// New crontab instance.
func NewCrontab() CrontabInterface {
	log.Debugf("[crontab] new crontab.")
	o := &crontab{}
	o.mu = new(sync.RWMutex)
	o.tickers = make(map[string]TickerInterface)
	return o
}
