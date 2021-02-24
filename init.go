// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"context"
	"errors"
)

var (
	ErrAccessSecond        = errors.New("second no access")
	ErrAccessMinute        = errors.New("minute no access")
	ErrAccessHour          = errors.New("hour no access")
	ErrAccessDay           = errors.New("day no access")
	ErrAccessMonth         = errors.New("month no access")
	ErrAccessWeek          = errors.New("week no access")
	ErrAccessYear          = errors.New("year no access")
	ErrAccessTimeline      = errors.New("timeline no access")
	ErrStrategyFormat      = errors.New("invalid format")
	ErrStrategyFormatSpace = errors.New("invalid format space")
	ErrLockFailed          = errors.New("locked by other process")
	ErrNoTicker            = errors.New("no ticker added")
)

type Handler func(ctx context.Context, ticker TickerInterface) error
