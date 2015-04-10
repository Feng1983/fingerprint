package main

import (
	"sort"
	"time"
	"net"
	
        //"bytes"
        //"encoding/binary"
	"fmt"
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
	t1 := time.Unix(getTime(dd),0)
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
	return time.Unix(tt,0).Format("2006-01-02 15:04:05")
}

func getTime(tt string) int64 {
    
    the_time, err := time.Parse("2006-01-02", tt)
    if err != nil {
	return -1
    }
    y,m,d := the_time.Date()
    newtime := time.Date(y,m,d,0,0,0,0,time.Local)
    return newtime.Unix()
}
func getTimeDetail(tt string) int64{
    the_time, err := time.Parse("2006-01-02 15:04:05", tt)
    if err != nil {
        return -1
    }
    return the_time.Unix()-8*60*60
}

func Str2Mac(mc string) uint64{
        bmac,err:=net.ParseMAC(mc)
        if err!=nil{
		fmt.Println(bmac)
                return 0
        }
        return BytesToMac(bmac)
}
func BytesToMac(b[] byte) uint64{
        var r uint64
        var x uint8

        for i:=0;i<6;i++{
                //x = uint64(binary.LittleEndian.Uint8(b[i:i+1]))
		x=uint8(b[i])
                r = r<<8 + uint64(x)
        }
	//fmt.Println(b,r)
    	//r = uint64(binary.LittleEndian.Uint64(b))
        return r
}

/*
func main(){
	//fmt.Println(time.Parse("2006-01-02","2015-03-12"))
	fmt.Println(Dateplus("2015-03-12",15))
	v2:= Dateplus("2015-03-12",15)
	fmt.Println(Int2date(v2),Int2date(1427385600))
	fmt.Println(Today(),getTime("2015-03-27"))
	fmt.Println(Datebyoff(0), getTimeDetail("2015-03-27 23:23:46"),time.Now().Unix(),time.Now().String())
}*/
