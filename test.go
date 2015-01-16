package main                                                                                                                                             
import(
    "fmt"
    "time"
    "strconv"
)

func makeCakeAndSend(cs chan string){
    for i:=1;i<=3;i++{
        cakeName:= "Strewberry Cake "+ strconv.Itoa(i)
        fmt.Println("Make a cake and send...",cakeName)
        cs <- cakeName
    }
}

func receivCakeAndPack(cs chan string){
    for i:=1;i<=3;i++{
        s:= <-cs
        fmt.Println("Packing received cake:" ,s)
    }
}

func main(){
    cs:= make(chan string)
    go makeCakeAndSend(cs)
    go receivCakeAndPack(cs)

    time.Sleep(4*1e9)
}
