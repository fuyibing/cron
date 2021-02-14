// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package cron

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	regexpFormatResetSpaces = regexp.MustCompile(`[\s]{2,}`)
	regexpFormatRecycles    = regexp.MustCompile(`^(\d+)\s*([a-zA-Z]?)$`)
	regexpFormatTimelines   = regexp.MustCompile(`^(\d+)\s*:\s*(\d+)\s*[:]?\s*(\d*)`)
	regexpFormatStandard    = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)$`)
	regexpFormatWithScope   = regexp.MustCompile(`^(\d+)-(\d+)$`)
	regexpFormatWithInteger = regexp.MustCompile(`^(\d+)$`)
)

// Strategy interface.
type Strategy interface {
	Err() error
	Format() string
	GetDays() map[int]int
	GetMonths() map[int]int
	GetWeeks() map[int]int
	GetTimes() map[string]int
	Update(time.Time)
	Validate(time.Time) (can bool, err error)
}

// Strategy struct.
type strategy struct {
	err             error
	format          string
	lastRuntime     time.Time
	listTimelines   map[string]int
	listMinutes     map[int]int
	listHours       map[int]int
	listDays        map[int]int
	listMonths      map[int]int
	listWeeks       map[int]int
	mu              *sync.RWMutex
	recycle         bool
	recycleDistance int64
	timeline        bool
}

// Create strategy instance.
func NewStrategy(format string) Strategy {
	o := new(strategy)
	o.mu = new(sync.RWMutex)
	o.lastRuntime = time.Unix(0, 0)
	o.listTimelines = make(map[string]int)
	o.generateFormat(format)
	o.generateTimeline()
	return o
}

// Return strategy error.
func (o *strategy) Err() error {
	return o.err
}

// Return strategy format.
func (o *strategy) Format() string {
	return o.format
}

// Return strategy days.
func (o *strategy) GetDays() map[int]int {
	return o.listDays
}

// Return strategy months.
func (o *strategy) GetMonths() map[int]int {
	return o.listMonths
}

// Return strategy weeks.
func (o *strategy) GetWeeks() map[int]int {
	return o.listWeeks
}

// Show strategy times.
func (o *strategy) GetTimes() map[string]int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.listTimelines
}

// Update last validate.
func (o *strategy) Update(t time.Time) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lastRuntime = t
}

// Validate strategy.
//
// return bool, can run handler or not.
// return err, parse strategy error.
func (o *strategy) Validate(t time.Time) (bool, error) {
	// 1. has error.
	if o.err != nil {
		return false, o.err
	}
	// 2. recycle mode.
	if o.recycle {
		o.mu.RLock()
		defer o.mu.RUnlock()
		diff := t.Sub(o.lastRuntime).Milliseconds()
		if diff >= o.recycleDistance {
			return true, nil
		}
		return false, nil
	}
	// 3. timeline mode.
	if o.timeline {
		o.mu.RLock()
		defer o.mu.RUnlock()
		_, ok := o.listTimelines[t.Format("15:04:05")]
		return ok, nil
	}
	// 4. standard mode.
	if o.matchWeek(t) && o.matchMonth(t) && o.matchDay(t) && o.matchTimeline(t) {
		return true, nil
	}
	// n. ended
	return false, nil
}

// Generate strategy format.
//   3s, 5m, 2h, 3d
//   0 0 * * *
//   * 0-3 * * *
//   * 0-3 1-5,9,12 * *
//   * * * */1 *
func (o *strategy) generateFormat(format string) {
	// 1. space manager.
	//    trim spaces and middle spaces.
	o.format = strings.TrimSpace(format)
	o.format = regexpFormatResetSpaces.ReplaceAllString(o.format, " ")
	// 2. recycle mode.
	if m := regexpFormatRecycles.FindStringSubmatch(format); len(m) == 3 {
		n, _ := strconv.ParseInt(m[1], 0, 64)
		switch m[2] {
		case "", "s", "S":
			o.recycleDistance = n
		case "m", "M":
			o.recycleDistance = n * 60
		case "h", "H":
			o.recycleDistance = n * 3600
		case "d", "D":
			o.recycleDistance = n * 86400
		default:
			o.err = errors.New(fmt.Sprintf("unknown recycle mode: %s", m[2]))
		}
		o.recycle = true
		o.recycleDistance *= int64(time.Microsecond)
		o.recycleDistance -= 500
		return
	}
	// 3. force timeline
	if m := regexpFormatTimelines.FindStringSubmatch(o.format); len(m) > 0 {
		o.timeline = true
		for _, s := range strings.Split(o.format, ",") {
			if m1 := regexpFormatTimelines.FindStringSubmatch(s); len(m1) == 4 {
				var n3 int64 = 0
				if m1[3] != "" {
					n3, _ = strconv.ParseInt(m1[3], 0, 32)
				}
				n1, _ := strconv.ParseInt(m1[1], 0, 32)
				n2, _ := strconv.ParseInt(m1[2], 0, 32)
				s1 := fmt.Sprintf("%02d:%02d:%02d", n1, n2, n3)
				o.listTimelines[s1] = int(n1)
			}
		}
		return
	}
	// 4. standard.
	if m := regexpFormatStandard.FindStringSubmatch(o.format); len(m) == 6 {
		// 4.1 column parse.
		for i := 1; i <= 5; i++ {
			s := m[i]
			// 4.2 split with comma.
			for _, x := range strings.Split(s, ",") {
				x = strings.TrimSpace(x)
				// 4.2.1 empty
				if x == "" {
					continue
				}
				// 4.2.2 scope
				if m1 := regexpFormatWithScope.FindStringSubmatch(x); len(m1) == 3 {
					n1, _ := strconv.ParseInt(m1[1], 0, 32)
					n2, _ := strconv.ParseInt(m1[2], 0, 32)
					for i1 := int(n1); i1 <= int(n2); i1++ {
						o.generateColumn(i, i1)
					}
					continue
				}
				// 4.2.3 integer
				if m2 := regexpFormatWithInteger.FindStringSubmatch(x); len(m2) == 2 {
					n2, _ := strconv.ParseInt(m2[1], 0, 32)
					o.generateColumn(i, int(n2))
					continue
				}
			}
		}
	}
}

func (o *strategy) generateColumn(col int, val int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// Week.
	if col == 5 {
		if o.listWeeks == nil {
			o.listWeeks = make(map[int]int)
		}
		val = val % 7
		o.listWeeks[val] = val
		return
	}
	// Month.
	if col == 4 {
		if o.listMonths == nil {
			o.listMonths = make(map[int]int)
		}
		if val >= 1 && val <= 12 {
			o.listMonths[val] = val
		}
		return
	}
	// Day.
	if col == 3 {
		if o.listDays == nil {
			o.listDays = make(map[int]int)
		}
		if val >= 1 && val <= 31 {
			o.listDays[val] = val
		}
		return
	}
	// Hour.
	if col == 2 {
		if o.listHours == nil {
			o.listHours = make(map[int]int)
		}
		if val >= 0 && val < 23 {
			o.listHours[val] = val
		}
		return
	}
	// Minute.
	if col == 1 {
		if o.listMinutes == nil {
			o.listMinutes = make(map[int]int)
		}
		if val >= 0 && val <= 59 {
			o.listMinutes[val] = val
		}
		return
	}
}

func (o *strategy) generateTimeline() {
	if o.timeline || o.recycle {
		return
	}
	t := time.Now()
	// no time specified.
	if o.listMinutes == nil && o.listHours == nil {
		s := t.Format("15:04:05")
		o.listTimelines[s] = t.Hour()
		return
	}
	// has hour.
	if o.listHours != nil {
		for _, hour := range o.listHours {
			if o.listMinutes == nil {
				s := fmt.Sprintf("%02d:%02d:00", hour, t.Minute())
				o.listTimelines[s] = hour
			} else {
				for _, minute := range o.listMinutes {
					s := fmt.Sprintf("%02d:%02d:00", hour, minute)
					o.listTimelines[s] = hour
				}
			}
		}
		return
	}
	// has minute.
	if o.listMinutes != nil {
		hour := t.Hour()
		for _, minute := range o.listMinutes {
			s := fmt.Sprintf("%02d:%02d:00", hour, minute)
			o.listTimelines[s] = hour
		}
	}
}

func (o *strategy) matchTimeline(t time.Time) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	s := t.Format("15:04:05")
	_, ok := o.listTimelines[s]
	return ok
}

func (o *strategy) matchDay(t time.Time) bool {
	if o.listDays == nil {
		return true
	}
	o.mu.RLock()
	defer o.mu.RUnlock()
	_, ok := o.listDays[t.Day()]
	return ok
}

func (o *strategy) matchMonth(t time.Time) bool {
	if o.listMonths == nil {
		return true
	}
	w := int(t.Month())
	o.mu.RLock()
	defer o.mu.RUnlock()
	_, ok := o.listMonths[w]
	return ok
}

func (o *strategy) matchWeek(t time.Time) bool {
	if o.listWeeks == nil {
		return true
	}
	w := int(t.Weekday())
	o.mu.RLock()
	defer o.mu.RUnlock()
	_, ok := o.listWeeks[w]
	return ok
}
