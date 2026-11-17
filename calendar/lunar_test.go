package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Month starts from the Hong Kong Observatory conversion tables
// (https://www.hko.gov.hk/en/gts/time/calendar/text/files/T1954e.txt etc.).
func TestSolarToLunarMatchesHKO(t *testing.T) {
	cases := []struct {
		solar            string
		month, day, leap int
	}{
		{"1933-07-22", 5, 30, 1},
		{"1933-07-23", 6, 1, 0},
		{"1954-11-25", 11, 1, 0},
		{"1956-12-02", 11, 1, 0},
		{"1978-09-02", 7, 30, 0},
		{"1978-09-03", 8, 1, 0},
		{"1996-07-15", 5, 30, 0},
		{"1996-07-16", 6, 1, 0},
		{"1996-10-12", 9, 1, 0},
		{"2033-08-25", 8, 1, 0},
		{"2033-12-22", 11, 1, 1},
	}
	for _, c := range cases {
		d, err := time.ParseInLocation("2006-01-02", c.solar, CST)
		assert.NoError(t, err)
		l := NewSolarTime(d).Convert()
		assert.Equal(t, []int{c.month, c.day, c.leap}, []int{l.Month(), l.Day(), boolToInt(l.IsLeapMonth())}, c.solar)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
