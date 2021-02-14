// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"errors"
	"sync"
	"time"
)

// Crontab interface.
type Crontab interface {
	Add(...Ticker) Crontab
	Start() error
	initialize()
}

// Crontab struct.
type crontab struct {
	mu      *sync.RWMutex
	tickers []Ticker
}

// Add ticker.
func (o *crontab) Add(tickers ...Ticker) Crontab {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.tickers = append(o.tickers, tickers...)
	return o
}

// Start crontab.
func (o *crontab) Start() error {
	// empty ticker.
	if len(o.tickers) == 0 {
		return errors.New("no ticker found")
	}
	// listen channel.
	o.listen()
	return nil
}

// Initialize crontab.
func (o *crontab) initialize() {
	o.mu = new(sync.RWMutex)
	o.tickers = make([]Ticker, 0)
}

// Listen channel.
func (o *crontab) listen() {
	go func() {
		defer o.listen()
		for range time.NewTicker(time.Second).C {
			go o.loop()
		}
	}()
}

// Loop tickers.
func (o *crontab) loop() {
	tm := time.Now()
	for _, tick := range o.tickers {
		go func(ticker Ticker) {
			if ok, _ := ticker.Validate(tm); ok {
				_ = ticker.Run(tm)
			}
		}(tick)
	}
}
