
package main

import (
    "fmt"
    "runtime"
    "sync"
)

var once = new(sync.Once)


func f(){
	fmt.Println("one do...")
}
func greeting(wg *sync.WaitGroup) {
    //once.Do(func() {
    //    fmt.Println("one do...")
    //})
	once.Do(f)

    fmt.Println("greeting")
    wg.Done()
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    defer fmt.Println("start....")

    wg := new(sync.WaitGroup)
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go greeting(wg)
    }
    wg.Wait()
}
