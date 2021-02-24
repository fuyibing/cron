// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"errors"
	"sync"
	"time"

	"github.com/fuyibing/cache"

	"github.com/fuyibing/log/v2"
)

// Ticker struct.
type ticker struct {
	handler    Handler
	mu         *sync.RWMutex
	name       string
	running    bool
	singleNode bool
	strategy   StrategyInterface
}

// Return ticker name.
func (o *ticker) Name() string {
	return o.name
}

// Run ticker.
func (o *ticker) Run(t time.Time) {
	// 1. return if status is running.
	o.mu.Lock()
	if o.running {
		o.mu.Unlock()
		return
	}
	o.running = true
	o.mu.Unlock()
	// 2. run response.
	ctx := log.NewContext()
	log.Debugfc(ctx, "[ticker=%s] begin ticker.", o.name)

	// 3. strategy check
	var err error
	if err = o.strategy.Err(); err != nil {
		log.Errorfc(ctx, "[ticker=%s] strategy error: %v.", o.name, err)
		return
	}
	if err = o.strategy.Validate(t); err != nil {
		log.Debugfc(ctx, "[ticker=%s] ticker ignored: %v.", o.name, err)
		return
	}
	// 4.
	defer func() {
		// duration and status.
		o.mu.Lock()
		o.running = false
		o.mu.Unlock()
		d := time.Now().Sub(t).Seconds()
		// result check.
		if r := recover(); r != nil {
			log.Errorfc(ctx, "[ticker=%s][d=%f] fatal error: %v.", o.name, d, r)
		} else {
			if err != nil {
				log.Errorfc(ctx, "[ticker=%s][d=%f] ticker fail: %v.", o.name, d, err)
			} else {
				log.Infofc(ctx, "[ticker=%s][d=%f] ticker completed.", o.name, d)
			}
		}
	}()
	// 5. single node.
	if o.singleNode {
		receipt := ""
		lock := cache.NewLock(o.name)
		if receipt, err = lock.Set(ctx); err != nil {
			return
		}
		if receipt == "" {
			err = errors.New("locked by other process")
			return
		}
		defer func() {
			_ = lock.Unset(ctx, receipt)
		}()
	}
	// 6. run processing.
	o.strategy.Refresh(t)
	err = o.handler(ctx, o)
}

// Set ticker handler.
func (o *ticker) SetHandler(handler Handler) TickerInterface {
	o.handler = handler
	return o
}

// Set single node running.
func (o *ticker) SingleNode(singleNode bool) TickerInterface {
	o.singleNode = singleNode
	return o
}

// Return ticker strategy.
func (o *ticker) Strategy() StrategyInterface {
	return o.strategy
}

// New ticker instance.
func NewTicker(name, format string, handler Handler) TickerInterface {
	log.Debugf("[ticker=%s] new ticker.", name)
	o := &ticker{mu: new(sync.RWMutex), name: name, running: false, singleNode: false}
	o.strategy = NewStrategy(format)
	o.SetHandler(handler)
	return o
}
