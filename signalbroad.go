package main                                                                                                                                             

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
    l := new(sync.Mutex)
    c := sync.NewCond(l)

    for i:=0;i<10;i++{
        go func(i int){
            fmt.Printf("waiting %d\n",i)
            l.Lock()
            defer l.Unlock()
            c.Wait()
            fmt.Printf("go %d\n",i)
        }(i)
    }

    for i:=0;i<3;i++{
        time.Sleep(1*time.Second)
        //c.Signal()
        fmt.Printf("%d\n",i)
    }
    c.Broadcast()
    time.Sleep(3*time.Second)

    x:=[]int{1,2,3,4}
    y:=[]int{4,5,6,7}
    x=append(x,y...)
    fmt.Println(x)

}
