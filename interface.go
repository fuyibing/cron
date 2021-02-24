// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package cron

import (
	"context"
	"time"
)

type Handler func(ctx context.Context, ticker TickerInterface) error

// Crontab interface.
type CrontabInterface interface {
	// Add ticker.
	AddTicker(...TickerInterface) CrontabInterface

	// Start crontab.
	Start() error
}

// Strategy interface.
type StrategyInterface interface {
	// Return strategy error.
	Err() error

	// Return strategy format.
	Format() string

	// Refresh last used time.
	Refresh(t time.Time)

	// Validate time is accessed.
	Validate(t time.Time) error
}

// Ticker interface.
type TickerInterface interface {
	// Return ticker name.
	Name() string

	// Run ticker.
	Run(time.Time)

	// Return ticker strategy.
	Strategy() StrategyInterface

	// Single node only.
	SingleNode(bool) TickerInterface
}
