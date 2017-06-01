# calendar

Chinese Lunar Calendar, Solar Calendar and cron time rules.

## About

	//Solar structure
	type Solar struct {
		time.Time
	}
	
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

```go
func NewSolar(year, month, day, hour, min, sec, nsec int) *Solar
```

```go
func NewSolarNow() *Solar
```

```go
func NewLunar(year, month, day, hour, min, sec, nsec int, leapFirst bool) *Lunar
```

```go
func NewLunarNow() *Lunar
```

Lunar or Solar has a method `Convert` to convert itself to the *opposite* one.

```go
func (s *Solar) Convert() *Lunar
```

```go
func (l *Lunar) Convert() *Solar
```

```go
func (l Lunar) Add(years, months, days, hours, mins, secs, nsecs int) *Lunar
```

## NOTICE

This package's year range is `[1900,2050]` and month range is `[1,12]`.

## Example

```go
package main

import "github.com/henrylee2cn/goutil/calendar"
import "fmt"

func main() {
	l0 := lunar.NewLunar(1988, 2, 11, 9, 9, 9, 0, false)
	fmt.Println(l0)
	s0 := l0.Convert()
	fmt.Println(s0)

	s1 := lunar.NewSolar(1988, 3, 28, 9, 9, 9, 0)
	fmt.Println(s1)
	l1 := s0.Convert()
	fmt.Println(l1)

	s := lunar.NewSolarNow()
	fmt.Println(s)

	l := lunar.NewLunarNow()
	fmt.Println(l)

	s2l := s.Convert()
	fmt.Println(s2l)

	l2s := l.Convert()
	fmt.Println(l2s)
}
```

## Cron

[cron](cron/README.md)