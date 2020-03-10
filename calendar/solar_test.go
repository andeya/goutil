package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddMonth(t *testing.T) {
	solar := NewSolarTime(time.Date(2020, 10, 31, 17, 40, 21, 100, CST))
	assert.Equal(t, "2021年02月28日 17时40分21秒", solar.AddMonth(4).String())
	assert.Equal(t, "2019年09月30日 17时40分21秒", solar.AddMonth(-13).String())
}
