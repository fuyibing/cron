// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
)

type StrategyMode int

const (
	StrategyModeRecycle StrategyMode = iota
)

var (
	StrategyRegexpFormat     = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)$`)
	StrategyRegexpRecycle    = regexp.MustCompile(`^(\d+)\s*([a-zA-Z]?)$`)
	StrategyRegexpResetSpace = regexp.MustCompile(`[\s]{2,}`)
)

// Strategy struct.
type strategy struct {
	err    error
	format string
	mu     *sync.RWMutex
	time   time.Time
}

// Return strategy error.
func (o *strategy) Err() error { return o.err }

// Return strategy format.
func (o *strategy) Format() string { return o.format }

// Refresh last used time.
func (o *strategy) Refresh(t time.Time) { o.time = t }

// Validate time is accessed.
func (o *strategy) Validate(t time.Time) error { return nil }

// Parse strategy format.
func (o *strategy) parse(format string) {
	// recycle mode.
	// eg: 10s、5m、2h、1d
	if m := StrategyRegexpRecycle.FindStringSubmatch(format); len(m) == 3 {
		switch m[2] {
		case "", "s", "S":
			format = fmt.Sprintf("*/%s * * * * *", m[1])
		case "m", "M":
			format = fmt.Sprintf("* */%s * * * *", m[1])
		case "h", "H":
			format = fmt.Sprintf("* * */%s * * *", m[1])
		case "d", "D":
			format = fmt.Sprintf("* * * */%s * *", m[1])
		default:
			o.err = errors.New("invalid strategy format")
			o.format = format
			return
		}
	}
	// assign format.
	o.format = format
	// parse format.
	m := StrategyRegexpFormat.FindStringSubmatch(format)
	if l := len(m); l != 7 {
		o.err = errors.New("invalid strategy format length")
		return
	}
	// parse[1]: second
	// parse[2]: minute
	// parse[3]: hour
	// parse[4]: day
	// parse[5]: month
	// parse[6]: week
}

// New strategy instance.
func NewStrategy(format string) StrategyInterface {
	log.Debugf("[strategy] new strategy: %s.", format)
	// 1. reset format space.
	format = strings.TrimSpace(format)
	format = StrategyRegexpResetSpace.ReplaceAllString(format, " ")
	// 2. create format.
	o := &strategy{mu: new(sync.RWMutex), time: time.Unix(0, 0)}
	o.parse(format)
	return o
}
