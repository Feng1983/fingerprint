package main

import (
	"fmt"
	"time"	
)

const INTERVAL_PERIOD time.Duration = 15* time.Second //24 * time.Hour

const HOUR_TO_TICK int = 17
const MINUTE_TO_TICK int = 47
const SECOND_TO_TICK int = 03

func runningRoutine() {
    ticker := updateTicker()
    for {
        <-ticker.C
        fmt.Println(time.Now(), "- just ticked")
        ticker = updateTicker()
    }
}

func runingT(){
    ticker := time.NewTicker(4 * time.Second)
    for {
	select{
           case	<- ticker.C:
		fmt.Println("start...ticker ", time.Now())
		go func(){
        		//fmt.Println(time.Now(),"- just ticked")
			time.Sleep(5* time.Second)
			fmt.Println(time.Now(),"- just ticked")
		}()
	}
    }
}

func updateTicker() *time.Ticker {
    nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
    if !nextTick.After(time.Now()) {
        nextTick = nextTick.Add(INTERVAL_PERIOD)
    }
    fmt.Println(nextTick, "- next tick")
    diff := nextTick.Sub(time.Now())
    return time.NewTicker(diff)
}

func main(){
	//runningRoutine()
	runingT()
}
