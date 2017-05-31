package main

import (
	"fmt"
	"github.com/henrylee2cn/goutil/lunar"
)

func main() {
	if lunar.MonthDays(1987, 1) != 30 {
		fmt.Println("1!=30: ", lunar.MonthDays(1987, 1))
		return
	}
	if lunar.MonthDays(1987, 2) != 29 {
		fmt.Println("2!=29: ", lunar.MonthDays(1987, 2))
		return
	}
	if lunar.MonthDays(1987, 3) != 30 {
		fmt.Println("3!=30: ", lunar.MonthDays(1987, 3))
		return
	}
	if lunar.MonthDays(1987, 4) != 29 {
		fmt.Println("4!=29: ", lunar.MonthDays(1987, 4))
		return
	}
	if lunar.MonthDays(1987, 5) != 30 {
		fmt.Println("5!=30: ", lunar.MonthDays(1987, 5))
		return
	}
	if lunar.MonthDays(1987, 6) != 30 {
		fmt.Println("6!=30: ", lunar.MonthDays(1987, 6))
		return
	}
	if lunar.MonthDays(1987, 13) != 29 {
		fmt.Println("闰6!=29: ", lunar.MonthDays(1987, 13))
		return
	}
	if lunar.MonthDays(1987, 7) != 30 {
		fmt.Println("7!=30: ", lunar.MonthDays(1987, 7))
		return
	}
	if lunar.MonthDays(1987, 8) != 30 {
		fmt.Println("8!=30: ", lunar.MonthDays(1987, 8))
		return
	}
	if lunar.MonthDays(1987, 9) != 29 {
		fmt.Println("9!=29: ", lunar.MonthDays(1987, 9))
		return
	}
	if lunar.MonthDays(1987, 10) != 30 {
		fmt.Println("10!=30: ", lunar.MonthDays(1987, 10))
		return
	}
	if lunar.MonthDays(1987, 11) != 29 {
		fmt.Println("11!=29: ", lunar.MonthDays(1987, 11))
		return
	}
	if lunar.MonthDays(1987, 12) != 29 {
		fmt.Println("12!=29: ", lunar.MonthDays(1987, 12))
		return
	}
	if lunar.MonthDays(1988, 1) != 30 {
		fmt.Println("——1!=30: ", lunar.MonthDays(1988, 1))
		return
	}

	l0 := lunar.NewLunar(1987, 7, 31, 9, 9, 9, true)
	fmt.Println(l0)
	s0 := l0.Convert()
	fmt.Println(s0)

	s1 := lunar.NewSolar(1987, 9, 7, 9, 9, 9)
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

	y, m, d := lunar.GanZhiYMD(2014, 5, 1) //should be:甲午 戊辰 壬申
	fmt.Println("2014.5.1", y, m, d)

	y, m, d = lunar.GanZhiYMD(2014, 5, 5) //should be:甲午 己巳 丙子
	fmt.Println("2014.5.5", y, m, d)

	a := lunar.AnimalYear(1900)  //should be 鼠
	a2 := lunar.AnimalYear(1988) //should be 龙
	fmt.Println(a, a2)

	z := lunar.ZhiHour(0)
	z1 := lunar.ZhiHour(23)
	z2 := lunar.ZhiHour(1)
	z3 := lunar.ZhiHour(12)
	z4 := lunar.ZhiHour(22)
	fmt.Println(z, z1, z2, z3, z4)
}
