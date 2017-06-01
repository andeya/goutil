package calendar

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

type FestivalMap map[string]string

func NewFestivalMap() FestivalMap {
	return make(FestivalMap)
}

func (fm FestivalMap) Add(key, val string) {
	fm[key] = val
}

func (fm FestivalMap) Del(key string) {
	delete(fm, key)
}

func (fm FestivalMap) Get(key string) (string, error) {
	desc, ok := fm[key]
	if ok {
		return desc, nil
	}
	return "", errors.New("FestivalMap KEY NotFound")
}

func (fm FestivalMap) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	file.Close()
	for k, v := range fm {
		file.WriteString(k + " " + v + "\n")
	}
	return nil
}

func NewFestivalsFromFile(filename string) FestivalMap {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fest := NewFestivalMap()
	r := bufio.NewReader(file)
	for {
		buf, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		line := strings.Trim(string(buf), " ")
		items := strings.Split(line, " ")
		date := items[0]
		desc := items[1]
		fest.Add(date, desc)
	}
	return fest
}

var (
	SolarFestivals = FestivalMap{
		"0101": "元旦",
		"0214": "情人节",
		"0308": "妇女节",
		"0312": "植树节",
		"0401": "愚人节",
		"0422": "地球日",
		"0501": "劳动节",
		"0504": "青年节",
		"0531": "无烟日",
		"0601": "儿童节",
		"0606": "爱眼日",
		"0701": "建党日",
		"0707": "抗战纪念日",
		"0801": "建军节",
		"0910": "教师节",
		"0918": "九·一八事变纪念日",
		"1001": "国庆节",
		"1031": "万圣节",
		"1111": "光棍节",
		"1201": "艾滋病日",
		"1213": "南京大屠杀纪念日",
		"1224": "平安夜",
		"1225": "圣诞节",
	}
	LunarFestivals = FestivalMap{
		"0101": "春节",
		"0115": "元宵节",
		"0202": "龙抬头",
		"0505": "端午节",
		"0707": "七夕",
		"0715": "中元节",
		"0815": "中秋节",
		"0909": "重阳节",
		"1208": "腊八节",
		"1223": "小年",
		"0100": "除夕",
	}
)
