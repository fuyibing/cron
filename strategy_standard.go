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

// Standard mode.
type strategyStandard struct {
	err     error         // error
	format  string        // format
	mu      *sync.RWMutex // lock
	seconds map[int]bool  // column=1
	minutes map[int]bool  // column=2
	hours   map[int]bool  // column=3
	days    map[int]bool  // column=4
	months  map[int]bool  // column=5
	weeks   map[int]bool  // column=6
	years   map[int]bool  // column=7
}

// New strategy instance.
func newStandardStrategy(format string) StrategyInterface {
	log.Debugf("[strategy=standard] register %s.", format)
	// Recycle.
	filter := format
	if m := StrategyRecycle.FindStringSubmatch(format); len(m) == 3 {
		switch m[2] {
		case "", "s", "S":
			filter = fmt.Sprintf("0/%s * * * * * *", m[1])
		case "m", "M":
			filter = fmt.Sprintf("* 0/%s * * * * *", m[1])
		case "h", "H":
			filter = fmt.Sprintf("* * 0/%s * * * *", m[1])
		case "d", "D":
			filter = fmt.Sprintf("* * * 0/%s * * *", m[1])
		default:
			filter = ""
		}
	}
	// Standard.
	o := &strategyStandard{format: format, mu: new(sync.RWMutex)}
	s := strings.Split(filter, " ")
	l := len(s)
	if l < 5 {
		o.err = ErrStrategyFormat
	} else {
		if l >= 6 {
			// 6 & 7 columns
			o.parseSeconds(s[0])
			o.parseMinutes(s[1])
			o.parseHours(s[2])
			o.parseDays(s[3])
			o.parseMonths(s[4])
			o.parseWeeks(s[5])
			if l == 7 {
				o.parseYears(s[6])
			} else {
				o.parseYears("*")
			}
		} else {
			// 5 columns
			o.parseSeconds("0")
			o.parseMinutes(s[0])
			o.parseHours(s[1])
			o.parseDays(s[2])
			o.parseMonths(s[3])
			o.parseWeeks(s[4])
			o.parseYears("*")
		}
	}
	return o
}

// Return strategy error.
func (o *strategyStandard) Err() error { return o.err }

// Return strategy format.
func (o *strategyStandard) Format() string {
	return o.format
}

// Validate time is accessed.
func (o *strategyStandard) Validate(t time.Time) error {
	// Return if error had.
	if o.err != nil {
		return o.err
	}
	// Iterate validates.
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, f := range []func(time.Time) error{
		o.validateSeconds,
		o.validateMinutes,
		o.validateHours,
		o.validateDays,
		o.validateMonths,
		o.validateWeeks,
		o.validateYears,
	} {
		if err := f(t); err != nil {
			return err
		}
	}
	return nil
}

// Parse columns.
func (o *strategyStandard) parse(col string, min, max int) (error, map[int]bool) {
	// return if has any error.
	if o.err != nil {
		return o.err, nil
	}
	// empty column value.
	if col = strings.TrimSpace(col); col == "" {
		o.err = ErrStrategyFormatSpace
		return o.err, nil
	}
	// any integer.
	if col == "*" || col == "?" {
		return nil, nil
	}
	// specified value.
	res := map[int]bool{}
	for _, str := range strings.Split(col, ",") {
		if str = strings.TrimSpace(str); str == "" {
			continue
		}
		for _, i := range o.parser(str, min, max) {
			res[i] = true
		}
	}
	return nil, res
}
func (o *strategyStandard) parser(str string, min, max int) (res []int) {
	res = make([]int, 0)
	// regexp: integer.
	if StrategyFormatInteger.MatchString(str) {
		a, _ := strconv.ParseInt(str, 0, 32)
		res = append(res, int(a))
		return
	}
	// regexp: range.
	if m := StrategyFormatRange.FindStringSubmatch(str); len(m) == 3 {
		a, _ := strconv.ParseInt(m[1], 0, 64)
		b, _ := strconv.ParseInt(m[2], 0, 64)
		for i := a; i <= b; i++ {
			res = append(res, int(i))
		}
		return
	}
	// regexp: div.
	if m := StrategyFormatDiv.FindStringSubmatch(str); len(m) == 3 {
		start := 0
		if StrategyFormatInteger.MatchString(m[1]) {
			a, _ := strconv.ParseInt(m[1], 0, 32)
			start = int(a)
		}
		if start < min {
			start = min
		}
		b, _ := strconv.ParseInt(m[2], 0, 32)
		offset := int(b)
		// append integer.
		for i := start; i <= max; i += offset {
			res = append(res, i)
		}
		return
	}
	// not support.
	return nil
}
func (o *strategyStandard) parseSeconds(col string) { _, o.seconds = o.parse(col, 0, 59) }
func (o *strategyStandard) parseMinutes(col string) { _, o.minutes = o.parse(col, 0, 59) }
func (o *strategyStandard) parseHours(col string)   { _, o.hours = o.parse(col, 0, 23) }
func (o *strategyStandard) parseDays(col string)    { _, o.days = o.parse(col, 1, 31) }
func (o *strategyStandard) parseMonths(col string)  { _, o.months = o.parse(col, 1, 12) }
func (o *strategyStandard) parseWeeks(col string)   { _, o.weeks = o.parse(col, 0, 6) }
func (o *strategyStandard) parseYears(col string)   { _, o.years = o.parse(col, 2020, 2048) }

// Validate columns.
func (o *strategyStandard) validate(r map[int]bool, k int, err error) error {
	if r == nil {
		return nil
	}
	if _, ok := r[k]; ok {
		return nil
	}
	return err
}
func (o *strategyStandard) validateSeconds(t time.Time) error {
	return o.validate(o.seconds, t.Second(), ErrAccessSecond)
}
func (o *strategyStandard) validateMinutes(t time.Time) error {
	return o.validate(o.minutes, t.Minute(), ErrAccessMinute)
}
func (o *strategyStandard) validateHours(t time.Time) error {
	return o.validate(o.hours, t.Hour(), ErrAccessHour)
}
func (o *strategyStandard) validateDays(t time.Time) error {
	return o.validate(o.days, t.Day(), ErrAccessDay)
}
func (o *strategyStandard) validateMonths(t time.Time) error {
	return o.validate(o.months, int(t.Month()), ErrAccessMonth)
}
func (o *strategyStandard) validateWeeks(t time.Time) error {
	return o.validate(o.weeks, int(t.Weekday()), ErrAccessWeek)
}
func (o *strategyStandard) validateYears(t time.Time) error {
	return o.validate(o.years, t.Year(), ErrAccessYear)
}
