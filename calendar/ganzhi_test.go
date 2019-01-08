package calendar

import (
	"testing"
	"fmt"
)

func TestGanZhi(t *testing.T) {
	solarNow := NewSolarNow()
	fmt.Println(solarNow)
	lunarNow := NewLunarNow()
	fmt.Println(lunarNow)
	fmt.Println(AnimalYear(lunarNow.Year()))
	fmt.Println(lunarNow.Year(),lunarNow.Month(),lunarNow.Day())
	ganzhiY,ganzhiM,ganzhiD:=GanZhiYMD(solarNow.Year(),int(solarNow.Month()),solarNow.Day())
	fmt.Println(ganzhiY,ganzhiM,ganzhiD)
	fmt.Println(JieQiDay(2018,11))
	JieQiDay(2012,1)
	JieQiDay(2012,2)
	JieQiDay(2012,3)
	JieQiDay(2012,4)
	JieQiDay(2012,5)
	JieQiDay(2012,6)
	JieQiDay(2012,7)
	JieQiDay(2012,8)
	JieQiDay(2013,4)
	JieQiDay(1993,3)
	JieQiDay(1980,1)
	JieQiDay(1991,5)
	//command := "docker ps -a"
	//cmd := exec.Command("/bin/bash", "-c", command)
	//bytes,err := cmd.Output()
	//if err != nil {
	//	log.Println(err)
	//}
	//resp := string(bytes)
	//log.Println(resp)
}
