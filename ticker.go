// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fuyibing/lock"

	"github.com/fuyibing/log/v2"
)

// Ticker interface.
type Ticker interface {
	Name() string
	Run(time.Time) error
	SingleNode(bool) Ticker
	Validate(time.Time) (bool, error)
}

// Ticker handler.
type TickerHandler func(context.Context, Ticker) error

// Ticker struct.
type ticker struct {
	handler    TickerHandler
	name       string
	strategy   Strategy
	singleNode bool
}

// Create ticker instance.
//   t1 := cron.NewTicker("t1", "3s", f1)
//   t2 := cron.NewTicker("t2", "18:00", f2)
//   t3 := cron.NewTicker("t3", "00:00:00,03:30:30", f3)
//   t4 := cron.NewTicker("t4", "* * 1-3,11,16,21 * *", f4)
func NewTicker(name, format string, handler TickerHandler) Ticker {
	o := &ticker{name: name, handler: handler}
	o.strategy = NewStrategy(format)
	return o
}

// Return ticker name.
func (o *ticker) Name() string {
	return o.name
}

// Run ticker.
func (o *ticker) Run(t time.Time) (err error) {
	var ctx = log.NewContext()
	// begin
	log.Infofc(ctx, "[cron][ticker=%s] begin run.", o.name)
	t1 := time.Now()
	defer func() {
		d1 := time.Now().Sub(t1).Seconds()
		// catch panic.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			log.Errorfc(ctx, "[cron][ticker=%s][d=%f] %s.", o.name, d1, err)
		} else {
			if err != nil {
				log.Warnfc(ctx, "[cron][ticker=%s][d=%f] %s.", o.name, d1, err)
			} else {
				log.Infofc(ctx, "[cron][ticker=%s][d=%f] run completed.", o.name, d1)
			}
		}
	}()
	// Strategy error.
	if err = o.strategy.Err(); err != nil {
		return
	}
	// No handler.
	if o.handler == nil {
		err = errors.New("handler not specified")
		return
	}
	// Signal node.
	if o.singleNode {
		l := lock.New(o.name)
		s1, e1 := l.Set(ctx)
		if e1 != nil {
			err = e1
			return
		}
		if !s1 {
			err = errors.New("can not get redis lock")
			return
		}
		defer func() {
			_, _ = l.Unset(ctx)
		}()
	}
	// Call handler.
	o.strategy.Update(t)
	err = o.handler(ctx, o)
	return
}

// Use single node.
func (o *ticker) SingleNode(singleNode bool) Ticker {
	o.singleNode = singleNode
	return o
}

// Validate strategy.
func (o *ticker) Validate(t time.Time) (bool, error) {
	return o.strategy.Validate(t)
}
