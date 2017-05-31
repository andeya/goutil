// Chinese Lunar Calendar Package.
package lunar

import (
	"fmt"
	"strconv"
	"time"
)

const (
	MinYear = 1900
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

	base = time.Date(MinYear, 1, 31, 0, 0, 0, 0, time.UTC)
)

//Solar structure
type Solar struct {
	time.Time
}

func NewSolar(year, month, day, hour, min, sec int) *Solar {
	if !isYearValid(year) {
		return nil
	}
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	return &Solar{t}
}

func NewSolarNow() *Solar {
	return &Solar{time.Now()}
}

func (s *Solar) String() string {
	return fmt.Sprintf("%d年%02d月%02d日 %2d时%2d分%2d秒",
		s.Year(), s.Month(), s.Day(), s.Hour(), s.Minute(), s.Second())
}

func (s *Solar) Festival(fm FestivalMap) (string, error) {
	m := fmt.Sprintf("%2d", int(s.Month()))
	d := fmt.Sprintf("%2d", s.Day())

	return fm.Get(m + d)
}

//Luanr structure
type Lunar struct {
	year        int
	month       int
	day         int
	hour        int
	minute      int
	second      int
	leapMonth   int
	isLeapMonth bool
}

func NewLunar(year, month, day, hour, min, sec int, leapFirst bool) *Lunar {
	if !isYearValid(year) {
		return nil
	}
	l := Lunar{}.Add(year, month, day, hour, min, sec)
	if leapFirst && l.month == l.leapMonth {
		l.isLeapMonth = true
	}
	return l
}

var (
	hourSecond = 60 * 60
	daySecond  = 24 * hourSecond
)

func (l Lunar) Add(years, months, days, hours, mins, secs int) *Lunar {
	l.year += years
	l.month += months
	l.day += days
	l.hour += hours
	l.minute += mins
	l.second += secs
	x := hourSecond*l.hour + 60*l.minute + l.second
	l.day += x / daySecond
	x = x % daySecond
	l.hour = x / hourSecond
	x = x % hourSecond
	l.minute = x / 60
	l.second = x % 60
	l.leapMonth = LeapMonth(l.year)
	l.isLeapMonth = false
	if l.leapMonth == 0 {
		l.year += l.month / 12
		l.month = l.month % 12
		l.leapMonth = LeapMonth(l.year)
	} else {
		addyear := l.month / 13
		l.year += addyear
		l.month = l.month % 13
		l.leapMonth = LeapMonth(l.year)
		if addyear > 0 && l.month > l.leapMonth {
			l.month--
			if l.month == l.leapMonth {
				l.isLeapMonth = true
			}
		}
	}

	if l.isLeapMonth {
		monthDays := LeapDays(l.year)
		l.month += l.day / monthDays
		l.day %= monthDays

	} else {
		monthDays := MonthDays(l.year, l.month)
		l.month += l.day / monthDays
		l.day %= monthDays
	}

	l.isLeapMonth = false
	if l.leapMonth == 0 {
		l.year += l.month / 12
		l.month = l.month % 12
		l.leapMonth = LeapMonth(l.year)

	} else {
		addyear := l.month / 13
		l.year += addyear
		l.month = l.month % 13
		l.leapMonth = LeapMonth(l.year)
		if addyear > 0 && l.month > l.leapMonth {
			l.month--
			if l.month == l.leapMonth {
				l.isLeapMonth = true
			}
		}
	}

	if !isYearValid(l.year) {
		return nil
	}
	if l.isLeapMonth {
		monthDays := LeapDays(l.year)
		l.month += l.day / monthDays
		l.day %= monthDays

	} else {
		monthDays := MonthDays(l.year, l.month)
		l.month += l.day / monthDays
		l.day %= monthDays
	}
	return &l
}

func NewLunarNow() *Lunar {
	return NewSolarNow().Convert()
}

func (l *Lunar) String() string {
	return fmt.Sprintf("%s%s%s %2d时%2d分%2d秒", YearString(l.Year()), MonthString(l.Month()), DayString(l.Day()), l.Hour(), l.Minute(), l.Second())
}

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
		temp = YearDays(i)
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
			temp = MonthDays(year, i)
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

	return &Lunar{year, month, day, s.Hour(), s.Minute(), s.Second(), LeapMonth(year), isLeap}
}

func (l *Lunar) Convert() *Solar {
	lyear := l.year
	lmonth := l.month
	lday := l.day
	offset := 0

	// increment year
	for i := MinYear; i < lyear; i++ {
		offset += YearDays(i)
	}

	// increment month
	// add days in all months up to the current month
	var cur int
	for cur = 1; cur < lmonth; cur++ {
		offset += MonthDays(lyear, cur)
		if cur == l.leapMonth {
			// add extra days for leap month
			offset += LeapDays(lyear)
		}
	}
	if l.isLeapMonth && l.leapMonth == lmonth {
		offset += MonthDays(lyear, cur)
	}

	// increment
	offset += lday - 1

	//BUG: maybe overflow
	d := time.Duration(offset*24) * time.Hour
	solar := base.Add(d)

	year := solar.Year()
	month := int(solar.Month())
	day := solar.Day()
	return NewSolar(year, month, day, l.Hour(), l.Minute(), l.Second())
}

/*
 * Common Methods
 */

func IsLeap(year int) bool {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return true
	}
	return false
}

//the total days of this year
func YearDays(year int) int {
	sum := 348
	for i := 0x8000; i > 0x8; i >>= 1 {
		if (lunarTable[year-MinYear] & i) != 0 {
			sum += 1
		}
	}
	return sum + LeapDays(year)
}

//which month leaps in this year?
//return 1-12(if there is one) or 0(no leap month).
func LeapMonth(year int) int {
	return int(lunarTable[year-MinYear] & 0xf)
}

//the total lunar months of this year
func LunarMonths(year int) int {
	if LeapMonth(year) == 0 {
		return 12
	}
	return 13
}

//the days of this year's leap month
func LeapDays(year int) int {
	if LeapMonth(year) != 0 {
		if (lunarTable[year-MinYear] & 0x10000) != 0 {
			return 30
		}
		return 29
	}
	return 0
}

//the days of the m-th month of this year
func MonthDays(year, month int) int {
	if (lunarTable[year-MinYear] & (0x10000 >> uint(month))) != 0 {
		return 30
	}
	return 29
}

/*
 * Lunar Methods
 */

func (l *Lunar) Year() int {
	return l.year
}

func (l *Lunar) Month() (int, bool) {
	return l.month, l.isLeapMonth
}

func (l *Lunar) Day() int {
	return l.day
}

func (l *Lunar) Hour() int {
	return l.hour
}

func (l *Lunar) Minute() int {
	return l.minute
}

func (l *Lunar) Second() int {
	return l.second
}

func (l *Lunar) Festival(fm FestivalMap) (string, error) {
	m := fmt.Sprintf("%2d", l.month)
	d := fmt.Sprintf("%2d", l.day)

	return fm.Get(m + d)
}

// Used Only by Lunar Object
func YearString(year int) string {
	return strconv.Itoa(year) + "年"
}

// Used Only by Lunar Object
func MonthString(month int, leap bool) string {
	if leap {
		return "闰" + lunarMonthNameTable[(month-1)%12] + "月"
	}
	return lunarMonthNameTable[(month-1)%12] + "月"
}

// Used Only by Lunar Object
func DayString(day int) (s string) {
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

func isYearValid(year int) bool {
	if year > MaxYear || year < MinYear {
		fmt.Printf("Invalid Year: %d, Year Range[%d - %d].\n", year, MinYear, MaxYear)
		return false
	}
	return true
}
