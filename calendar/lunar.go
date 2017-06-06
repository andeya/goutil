package calendar

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// MinYear smallest year of support
	MinYear = 1900
	// MaxYear largest year of support
	MaxYear = 2050
)

var (
	// 具体的方法是:用一位来表示一个月的大小，大月记为1，小月记为0，
	// 这样就用掉12位(无闰月)或13位(有闰月)，再用高四位来表示闰月的月份，没有闰月记为0。
	// 例如：2000年的信息数据是0xc96，化成二进制就是110010010110B，表示的
	// 含义是:1、2、5、8、10、11月大，其余月份小。
	// Since 1900~2050
	lunarTable = []int{
		0x04bd8, 0x04ae0, 0x0a570, 0x054d5, 0x0d260,
		0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2,
		0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255,
		0x0b540, 0x0d6a0, 0x0ada2, 0x095b0, 0x14977,
		0x04970, 0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40,
		0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970,
		0x06566, 0x0d4a0, 0x0ea50, 0x06e95, 0x05ad0,
		0x02b60, 0x186e3, 0x092e0, 0x1c8d7, 0x0c950,
		0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4,
		0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557,
		0x06ca0, 0x0b550, 0x15355, 0x04da0, 0x0a5d0,
		0x14573, 0x052d0, 0x0a9a8, 0x0e950, 0x06aa0,
		0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570,
		0x05260, 0x0f263, 0x0d950, 0x05b57, 0x056a0,
		0x096d0, 0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4,
		0x0d250, 0x0d558, 0x0b540, 0x0b5a0, 0x195a6,
		0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a,
		0x06a50, 0x06d40, 0x0af46, 0x0ab60, 0x09570,
		0x04af5, 0x04970, 0x064b0, 0x074a3, 0x0ea50,
		0x06b58, 0x055c0, 0x0ab60, 0x096d5, 0x092e0,
		0x0c960, 0x0d954, 0x0d4a0, 0x0da50, 0x07552,
		0x056a0, 0x0abb7, 0x025d0, 0x092d0, 0x0cab5,
		0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9,
		0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930,
		0x07954, 0x06aa0, 0x0ad50, 0x05b52, 0x04b60,
		0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65, 0x0d530,
		0x05aa0, 0x076a3, 0x096d0, 0x04bd7, 0x04ad0,
		0x0a4d0, 0x1d0b6, 0x0d250, 0x0d520, 0x0dd45,
		0x0b5a0, 0x056d0, 0x055b2, 0x049b0, 0x0a577,
		0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0,
	}
	lunarMonthNameTable = []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "腊"}
	monthStr1           = []string{"初", "十", "廿", "卅"}
	monthStr2           = []string{"日", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

	base = time.Date(MinYear, 1, 31, 0, 0, 0, 0, CST)

	// CST  CST China Standard Time UT 8:00
	CST = time.FixedZone("CST", 3600*8)
)

//Luanr structure
type Lunar struct {
	year        int
	month       int
	day         int
	hour        int
	minute      int
	second      int
	nanosecond  int
	leapMonth   int
	isLeapMonth bool
}

// NewLunarNow creates current lunar time.
func NewLunarNow() *Lunar {
	return NewSolarNow().Convert()
}

// NewLunar creates a lunar time.
func NewLunar(year, month, day, hour, min, sec int, nsec int, leapFirst bool) *Lunar {
	if !isYearValid(year) {
		return nil
	}
	if year < 0 || month < 0 || day < 0 || hour < 0 || min < 0 || sec < 0 {
		return nil
	}
	leapMonth := LeapMonth(year)
	if leapMonth > 0 {
		if month > leapMonth || leapFirst && month == leapMonth {
			month++
		}
	}
	l := Lunar{}.Add(year, month, day, hour, min, sec, nsec)
	return l
}

// LunarZero lunar zero time.
var LunarZero = NewLunar(MinYear, 1, 1, 0, 0, 0, 0, false)

// IsLunarZero judge whether the lunar time is equal to lunar zero time.
func IsLunarZero(l *Lunar) bool {
	return l.Equal(LunarZero)
}

const (
	hourSecond = 60 * 60
	daySecond  = 24 * hourSecond
)

// Add adds lunar time, and returns the result of lunar time.
func (l Lunar) Add(years, months, days, hours, mins, secs, nsecs int) *Lunar {
	// add years
	l.addYear(years)

	// add months or years
	l.addMonth(months)

	// add day, months or years
	l.addDay(days)

	// add hours, mins, secs, nsecs, day, months or years
	l.addSecond(hours*3600+mins*60+secs, nsecs)

	return &l
}

// add years
func (l *Lunar) addYear(years int) {
	if years == 0 {
		return
	}
	l.year += years
	_leapMonth := l.leapMonth
	l.leapMonth = LeapMonth(l.year)
	if l.leapMonth != _leapMonth {
		l.isLeapMonth = false
	}
}

// add months or years
func (l *Lunar) addMonth(months int) {
	if months == 0 {
		return
	}
	var op int
	if months > 0 {
		op = 1
	} else {
		op = -1
		months *= -1
	}
	var r int
	for i := 0; i < months; i++ {
		if !l.isLeapMonth && l.leapMonth > 0 && l.month == l.leapMonth {
			l.isLeapMonth = true
			continue
		}
		r = l.month + op
		if r <= 0 {
			l.addYear(op)
			l.month = 12
			if l.month == l.leapMonth {
				l.isLeapMonth = true
			} else {
				l.isLeapMonth = false
			}
			continue
		}
		if r > 12 {
			l.addYear(op)
			l.month = 1
			if l.month == l.leapMonth {
				l.isLeapMonth = true
			} else {
				l.isLeapMonth = false
			}
			continue
		}

		l.month = r
		if l.leapMonth == 0 || l.month != l.leapMonth {
			l.isLeapMonth = false
		}
	}
}

// add day, months or years
func (l *Lunar) addDay(days int) {
	if days == 0 {
		return
	}
	var op int
	if days > 0 {
		op = 1
	} else {
		op = -1
		days *= -1
	}
	var r int
	for i := 0; i < days; i++ {
		r = l.day + op
		if r <= 0 {
			l.addMonth(op)
			l.day = LunarMonthDays(l.year, l.month)
			continue
		}
		if r > LunarMonthDays(l.year, l.month) {
			l.addMonth(op)
			l.day = 1
			continue
		}
		l.day = r
	}
}

// add hours, mins, secs, nsecs, day, months or years
func (l *Lunar) addSecond(secs, nsecs int) {
	if secs == 0 && nsecs == 0 {
		return
	}
	secs += l.hour*3600 + l.minute*60 + l.second
	var days int
	days, l.hour, l.minute, l.second, l.nanosecond = SplitDuration(secs, l.nanosecond+nsecs)
	l.addDay(days)
}

// Convert converts to a solar calendar time.
func (l *Lunar) Convert() *Solar {
	lyear := l.year
	lmonth := l.month
	lday := l.day
	offset := 0

	// increment year
	for i := MinYear; i < lyear; i++ {
		offset += LunarYearDays(i)
	}

	// increment month
	// add days in all months up to the current month
	var cur int
	for cur = 1; cur < lmonth; cur++ {
		offset += LunarMonthDays(lyear, cur)
		if cur == l.leapMonth {
			// add extra days for leap month
			offset += LeapDays(lyear)
		}
	}
	if l.isLeapMonth && l.leapMonth == lmonth {
		offset += LunarMonthDays(lyear, cur)
	}

	// increment
	offset += lday - 1

	//BUG: maybe overflow
	d := time.Duration(offset*24) * time.Hour
	solar := base.Add(d)

	year := solar.Year()
	month := int(solar.Month())
	day := solar.Day()
	return NewSolar(year, month, day, l.Hour(), l.Minute(), l.Second(), l.Nanosecond(), CST)
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (l Lunar) Truncate(d time.Duration) *Lunar {
	return l.Convert().Truncate(d).Convert()
}

// MonthFirst returns lunar time of the month 1 day.
func (l Lunar) MonthFirst() *Lunar {
	l.day = 1
	return &l
}

// MonthLast returns to the last day of the month.
func (l Lunar) MonthLast() *Lunar {
	if l.isLeapMonth {
		l.day = LunarMonthDays(l.year, 13)
		return &l
	}
	l.day = LunarMonthDays(l.year, l.month)
	return &l
}

// Equal returns whether it is equal to the lunar time.
func (l *Lunar) Equal(lunar *Lunar) bool {
	return l.year == lunar.year &&
		l.month == lunar.month &&
		l.day == lunar.day &&
		l.hour == lunar.hour &&
		l.minute == lunar.minute &&
		l.second == lunar.second &&
		l.nanosecond == lunar.nanosecond &&
		l.isLeapMonth == lunar.isLeapMonth
}

// String formats time.
func (l *Lunar) String() string {
	return fmt.Sprintf("%s%s%s %2d时%2d分%2d秒", LunarYearString(l.Year()), LunarMonthString(l.Month(), l.IsLeapMonth()), LunarDayString(l.Day()), l.Hour(), l.Minute(), l.Second())
}

// Year returns year.
func (l *Lunar) Year() int {
	return l.year
}

// Month returns month.
func (l *Lunar) Month() int {
	return l.month
}

// LeapMonth returns leap month number.
func (l *Lunar) LeapMonth() int {
	return l.leapMonth
}

// IsLeapMonth returns whether it is leap month.
func (l *Lunar) IsLeapMonth() bool {
	return l.isLeapMonth
}

// Day returns day.
func (l *Lunar) Day() int {
	return l.day
}

// Weekday returns weekday.
func (l *Lunar) Weekday() time.Weekday {
	return l.Convert().Weekday()
}

// Hour returns hour.
func (l *Lunar) Hour() int {
	return l.hour
}

// Minute returns minute.
func (l *Lunar) Minute() int {
	return l.minute
}

// Second returns second.
func (l *Lunar) Second() int {
	return l.second
}

// Nanosecond returns nanosecond.
func (l *Lunar) Nanosecond() int {
	return l.nanosecond
}

// SetHour sets hour.
// NOTE: hour range [0:23].
func (l *Lunar) SetHour(hour int) *Lunar {
	if hour < 0 || hour > 23 {
		panic("hour range [0:23]")
	}
	l.hour = hour
	return l
}

// SetMinute sets minute.
// NOTE: minute range [0:59].
func (l *Lunar) SetMinute(minute int) *Lunar {
	if minute < 0 || minute > 59 {
		panic("minute range [0:59]")
	}
	l.minute = minute
	return l
}

// SetSecond sets second.
// NOTE: second range [0:59].
func (l *Lunar) SetSecond(second int) *Lunar {
	if second < 0 || second > 59 {
		panic("second range [0:59]")
	}
	l.second = second
	return l
}

// SetNanosecond sets nanosecond.
// NOTE: nanosecond range [0:999999999].
func (l *Lunar) SetNanosecond(nanosecond int) *Lunar {
	if nanosecond < 0 || nanosecond > 999999999 {
		panic("nanosecond range [0:999999999]")
	}
	l.nanosecond = nanosecond
	return l
}

// Copy returns a copy .
func (l Lunar) Copy() *Lunar {
	return &l
}

// Festival returns festival.
func (l *Lunar) Festival(fm FestivalMap) (string, error) {
	m := fmt.Sprintf("%2d", l.month)
	d := fmt.Sprintf("%2d", l.day)

	return fm.Get(m + d)
}

// LunarYearDays the total days of this year
func LunarYearDays(year int) int {
	sum := 348
	for i := 0x8000; i > 0x8; i >>= 1 {
		if (lunarTable[year-MinYear] & i) != 0 {
			sum += 1
		}
	}
	return sum + LeapDays(year)
}

// LeapMonth which month leaps in this year?
// return 1-12(if there is one) or 0(no leap month).
func LeapMonth(year int) int {
	return int(lunarTable[year-MinYear] & 0xf)
}

// LunarMonths the total lunar months of this year
func LunarMonths(year int) int {
	if LeapMonth(year) == 0 {
		return 12
	}
	return 13
}

// LeapDays the days of this year's leap month
func LeapDays(year int) int {
	if LeapMonth(year) != 0 {
		if (lunarTable[year-MinYear] & 0x10000) != 0 {
			return 30
		}
		return 29
	}
	return 0
}

// LunarMonthDays the days of the m-th month of this year.
func LunarMonthDays(year, month int) int {
	if (lunarTable[year-MinYear] & (0x10000 >> uint(month))) != 0 {
		return 30
	}
	return 29
}

// LunarYearString used Only by Lunar Object
func LunarYearString(year int) string {
	return strconv.Itoa(year) + "年"
}

// LunarMonthString used Only by Lunar Object
func LunarMonthString(month int, leap bool) string {
	if leap {
		return "闰" + lunarMonthNameTable[(month-1)%12] + "月"
	}
	return lunarMonthNameTable[(month-1)%12] + "月"
}

// LunarDayString used Only by Lunar Object
func LunarDayString(day int) (s string) {
	switch day {
	case 10:
		s = "初十"
	case 20:
		s = "二十"
	case 30:
		s = "三十"
	default:
		s = monthStr1[int(day/10)]
		s += monthStr2[day%10]
	}
	return
}

// SplitDuration split duration, accurate to second.
func SplitDuration(second, nanosecond int) (days, hours, minutes, seconds, nsecs int) {
	seconds, nsecs = norm(second, nanosecond, 1e9)
	minutes, seconds = norm(0, seconds, 60)
	hours, minutes = norm(0, minutes, 60)
	days, hours = norm(0, hours, 24)
	return
}

// norm returns nhi, nlo such that
//	hi * base + lo == nhi * base + nlo
//	0 <= nlo < base
func norm(hi, lo, base int) (nhi, nlo int) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	}
	if lo >= base {
		n := lo / base
		hi += n
		lo -= n * base
	}
	return hi, lo
}

func isYearValid(year int) bool {
	if year > MaxYear || year < MinYear {
		fmt.Printf("Invalid Year: %d, Year Range[%d - %d].\n", year, MinYear, MaxYear)
		return false
	}
	return true
}
