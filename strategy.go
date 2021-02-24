// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"regexp"
	"strings"
	"time"
)

var (
	StrategyResetSpace    = regexp.MustCompile(`[\s]{2,}`)
	StrategyRecycle       = regexp.MustCompile(`^(\d+)\s*([a-zA-Z]?)$`)
	StrategyTimeline      = regexp.MustCompile(`^(\d+)\s*:\s*(\d+)\s*[:]?\s*(\d*)`)
	StrategyFormatRange   = regexp.MustCompile(`^(\d+)-(\d+)$`)
	StrategyFormatInteger = regexp.MustCompile(`^(\d+)$`)
	StrategyFormatDiv     = regexp.MustCompile(`^([0-9\*]+)/(\d+)$`)
)

// Strategy interface.
type StrategyInterface interface {
	// Return strategy error.
	Err() error

	// Return strategy format.
	Format() string

	// Validate time is accessed.
	Validate(t time.Time) error
}

// New strategy instance.
// Return strategy interface of recycle or standard or timeline.
func NewStrategy(format string) StrategyInterface {
	// space operate.
	format = strings.TrimSpace(format)
	format = StrategyResetSpace.ReplaceAllString(format, " ")
	// use timeline.
	// eg: 00:10,01:02:03
	if StrategyTimeline.MatchString(format) {
		return newTimelineStrategy(format)
	}
	// standard mode.
	return newStandardStrategy(format)
}
