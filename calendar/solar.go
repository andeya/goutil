package calendar

import (
	"fmt"
	"time"
)

//Solar structure
type Solar struct {
	time.Time
}

// NewSolarNow creates current solar time.
func NewSolarNow() *Solar {
	return &Solar{time.Now().In(CST)}
}

// NewSolarTime creates a solar time from time.Time.
func NewSolarTime(t time.Time) *Solar {
	return &Solar{t.In(CST)}
}

// NewSolar creates a solar time.
func NewSolar(year, month, day, hour, min, sec int, nsec int, loc *time.Location) *Solar {
	if !isYearValid(year) {
		return nil
	}
	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc)
	return &Solar{t.In(CST)}
}

// String formats time.
func (s *Solar) String() string {
	return fmt.Sprintf("%d年%02d月%02d日 %2d时%2d分%2d秒",
		s.Year(), s.Month(), s.Day(), s.Hour(), s.Minute(), s.Second())
}

// Festival returns festival.
func (s *Solar) Festival(fm FestivalMap) (string, error) {
	m := fmt.Sprintf("%2d", int(s.Month()))
	d := fmt.Sprintf("%2d", s.Day())

	return fm.Get(m + d)
}

// Convert converts to a lunar calendar time.
func (s *Solar) Convert() *Lunar {
	var i int
	var leap int
	var isLeap bool
	var temp int

	var day int
	var month int
	var year int

	//offset days
	offset := int(s.Sub(base).Seconds() / 86400)

	for i = MinYear; i < MaxYear && offset > 0; i++ {
		temp = LunarYearDays(i)
		offset -= temp
	}

	if offset < 0 {
		offset += temp
		i--
	}

	year = i
	leap = LeapMonth(i)
	isLeap = false

	for i = 1; i < 13 && offset > 0; i++ {
		//leap month
		if leap > 0 && i == (leap+1) && isLeap == false {
			i--
			isLeap = true
			temp = LeapDays(year)
		} else {
			temp = LunarMonthDays(year, i)
		}
		//reset leap month
		if isLeap == true && i == (leap+1) {
			isLeap = false
		}
		offset -= temp
	}

	if offset == 0 && leap > 0 && i == (leap+1) {
		if isLeap {
			isLeap = false
		} else {
			isLeap = true
			i--
		}
	}

	if offset < 0 {
		offset += temp
		i--
	}

	month = i
	day = offset + 1

	return &Lunar{year, month, day, s.Hour(), s.Minute(), s.Second(), s.Nanosecond(), LeapMonth(year), isLeap}
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (s Solar) Truncate(d time.Duration) *Solar {
	return &Solar{s.Time.Truncate(d)}
}

// IsLeapYear determines whether it is a leap year.
func IsLeapYear(year int) bool {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return true
	}
	return false
}

// SolarMonthDays the days of the m-th month of this year.
func SolarMonthDays(year int, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	}
	return 0
}
