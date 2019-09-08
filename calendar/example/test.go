package main

import (
	"fmt"
	"time"

	"github.com/henrylee2cn/goutil/calendar"
)

func main() {
	fmt.Println("lunar zero:", calendar.LunarZero, calendar.IsLunarZero(calendar.LunarZero))
	fmt.Println("lunar zero -> solar:", calendar.LunarZero.Convert())

	if calendar.LunarMonthDays(1987, 1) != 30 {
		fmt.Println("1!=30: ", calendar.LunarMonthDays(1987, 1))
		return
	}
	if calendar.LunarMonthDays(1987, 2) != 29 {
		fmt.Println("2!=29: ", calendar.LunarMonthDays(1987, 2))
		return
	}
	if calendar.LunarMonthDays(1987, 3) != 30 {
		fmt.Println("3!=30: ", calendar.LunarMonthDays(1987, 3))
		return
	}
	if calendar.LunarMonthDays(1987, 4) != 29 {
		fmt.Println("4!=29: ", calendar.LunarMonthDays(1987, 4))
		return
	}
	if calendar.LunarMonthDays(1987, 5) != 30 {
		fmt.Println("5!=30: ", calendar.LunarMonthDays(1987, 5))
		return
	}
	if calendar.LunarMonthDays(1987, 6) != 30 {
		fmt.Println("6!=30: ", calendar.LunarMonthDays(1987, 6))
		return
	}
	if calendar.LunarMonthDays(1987, 13) != 29 {
		fmt.Println("闰6!=29: ", calendar.LunarMonthDays(1987, 13))
		return
	}
	if calendar.LunarMonthDays(1987, 7) != 30 {
		fmt.Println("7!=30: ", calendar.LunarMonthDays(1987, 7))
		return
	}
	if calendar.LunarMonthDays(1987, 8) != 30 {
		fmt.Println("8!=30: ", calendar.LunarMonthDays(1987, 8))
		return
	}
	if calendar.LunarMonthDays(1987, 9) != 29 {
		fmt.Println("9!=29: ", calendar.LunarMonthDays(1987, 9))
		return
	}
	if calendar.LunarMonthDays(1987, 10) != 30 {
		fmt.Println("10!=30: ", calendar.LunarMonthDays(1987, 10))
		return
	}
	if calendar.LunarMonthDays(1987, 11) != 29 {
		fmt.Println("11!=29: ", calendar.LunarMonthDays(1987, 11))
		return
	}
	if calendar.LunarMonthDays(1987, 12) != 29 {
		fmt.Println("12!=29: ", calendar.LunarMonthDays(1987, 12))
		return
	}
	if calendar.LunarMonthDays(1988, 1) != 30 {
		fmt.Println("——1!=30: ", calendar.LunarMonthDays(1988, 1))
		return
	}

	l0 := calendar.NewLunar(1987, 6, 15, 9, 9, 9, 0, false)
	fmt.Println(l0, l0.Weekday())
	s0 := l0.Convert()
	fmt.Println(s0, s0.Weekday())

	l0 = l0.Add(0, 1, 0, 0, 0, 0, 0)
	fmt.Println(l0)
	s0 = l0.Convert()
	fmt.Println(s0)

	l0 = l0.Truncate(time.Hour)
	fmt.Println(l0)
	s0 = l0.Convert()
	fmt.Println(s0)

	s1 := calendar.NewSolar(1987, 8, 9, 9, 0, 0, 0, time.UTC)
	fmt.Println(s1)
	l1 := s0.Convert()
	fmt.Println(l1)

	s := calendar.NewSolarNow()
	fmt.Println(s)

	l := calendar.NewLunarNow()
	fmt.Println(l)

	s2l := s.Convert()
	fmt.Println(s2l)

	l2s := l.Convert()
	fmt.Println(l2s)

	y, m, d := calendar.GanZhiYMD(2014, 5, 1) //should be:甲午 戊辰 壬申
	fmt.Println("2014.5.1", y, m, d)

	y, m, d = calendar.GanZhiYMD(2014, 5, 5) //should be:甲午 己巳 丙子
	fmt.Println("2014.5.5", y, m, d)

	a := calendar.AnimalYear(1900)  //should be 鼠
	a2 := calendar.AnimalYear(1988) //should be 龙
	fmt.Println(a, a2)

	z := calendar.ZhiHour(0)
	z1 := calendar.ZhiHour(23)
	z2 := calendar.ZhiHour(1)
	z3 := calendar.ZhiHour(12)
	z4 := calendar.ZhiHour(22)
	fmt.Println(z, z1, z2, z3, z4)

	l = calendar.NewSolar(2012, 7, 9, 23, 35, 51, 0, calendar.CST).Convert()
	fmt.Println("Mon Jul 9 23:35:51 2012 -> ", l)
	fmt.Println("2012年五月廿二 10时20分15秒 -> ", calendar.NewLunar(2012, 5, 22, 10, 20, 15, 0, false))

	fmt.Println(calendar.NewSolar(2018, 11, 28, 0, 0, 0, 0, calendar.CST).DiffWithYMD(2018, 12, 1))
}
