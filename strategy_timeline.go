// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package cron

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
)

// Timeline mode.
type strategyTimeline struct {
	err    error
	format string
	mu     *sync.RWMutex
	times  map[string]bool
}

// New strategy instance.
func newTimelineStrategy(format string) StrategyInterface {
	log.Debugf("[strategy=timeline] register %s.", format)
	o := &strategyTimeline{format: format, mu: new(sync.RWMutex), times: make(map[string]bool)}
	for _, s := range strings.Split(format, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		if m := StrategyTimeline.FindStringSubmatch(s); len(m) == 4 {
			n := []int64{0, 0, 0, 0}
			if m[3] == "" {
				m[3] = "0"
			}
			for i := 1; i <= 3; i++ {
				n[i], _ = strconv.ParseInt(m[i], 0, 64)
			}
			x := fmt.Sprintf("%02d:%02d:%02d", n[1], n[2], n[3])
			o.times[x] = true
			log.Debugf("[strategy=timeline] config %s.", x)
		}
	}
	if len(o.times) == 0 {
		o.err = ErrStrategyFormat
	}
	return o
}

// Return strategy error.
func (o *strategyTimeline) Err() error { return o.err }

// Return strategy format.
func (o *strategyTimeline) Format() string { return o.format }

// Validate time is accessed.
func (o *strategyTimeline) Validate(t time.Time) error {
	o.mu.RLock()
	defer o.mu.RUnlock()
	k := t.Format("15:04:05")
	if _, ok := o.times[k]; ok {
		return nil
	}
	return ErrAccessTimeline
}
