package main

import (
	"sort"
	"time"
	//"fmt"
)

type sortedIntMap struct{
	m map[int]float64
	s []int
}
func (sm *sortedIntMap) Len()int{
	return len(sm.m)
}

func (sm *sortedIntMap) Less(i,j int)bool{
	return sm.m[sm.s[i]]<sm.m[sm.s[j]]
}
func (sm *sortedIntMap) Swap(i, j int){
	sm.s[i],sm.s[j] = sm.s[j],sm.s[i]
}

func SortIntMap(m map[int]float64)[]int{
	sm:= new(sortedIntMap)
	sm.m = m
	sm.s = make([]int, len(m))
	i:=0
	for key:= range m{
		sm.s[i]=key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func Dateplus(dd string ,delta int) int64{
	t1,err := time.Parse("2006-01-02",dd)
	if err!=nil	{
		return -1
	}
	t2 := t1.AddDate(0,0,delta)
	return t2.Unix()
}
func Today() string{
	t1 := time.Now()
	return t1.Format("2006-01-02")
}

func Datebyoff(delta int) int64{
	return Dateplus(Today(),delta)
}

func Int2date(tt int64) string{
	return time.Unix(tt,0).Format("2006-01-02")
}

func getTime(tt string) int64 {
    the_time, err := time.Parse("2006-01-02 15:04:05", tt)
    if err != nil {
        unix_time := the_time.Unix()
        fmt.Println(unix_time)
    }
    return the_time.Unix()
}

//func main(){
//	fmt.Println(time.Parse("2006-01-02","2015-03-12"))
//	fmt.Println(dateplus("2015-03-12",8))
//	v2:= dateplus("2015-03-12",30)
//	fmt.Println(Int2date(v2))
//	fmt.Println(today())
//	fmt.Println(datebyoff(0))
//}
