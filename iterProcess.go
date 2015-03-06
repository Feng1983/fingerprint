package main                                                                                                                                             

import(
    "github.com/albrow/zoom"
    "fmt"
    //"strconv"
    //"strings"
    "time"
    //"os"
    //"bufio"
    "knn2"
	"github.com/deckarep/golang-set"
)
func Init(){                                                                                                                                             
    conf:= &zoom.Configuration{
        Address:"localhost:59999",
        Network:"tcp",
    }
    zoom.Init(conf)
    //zoom.Register(&Person{})
    //zoom.Register(&Rssiample{})
    //zoom.Register(&knn2.Finger{})
    zoom.Register(&knn2.ProcessData{})
}

func getData()[]* knn2.ProcessData{
	var sp []*knn2.ProcessData
	if err:= zoom.NewQuery("ProcessData").Order("Timestamp").Scan(&sp);err!=nil{
            panic(err)
    }
	return sp
}
func main(){
	Init()
	rset := mapset.NewSet()
	startTime:= time.Now()
	cnt:= 0
	for _, m :=range getData(){
		str_time := time.Unix(int64(m.Timestamp), 0).Format("2006-01-02 15:04:05")
		fmt.Println(m,str_time)
		//fmt.Println(m)
		rset.Add(m.Mac)
		//zoom.DeleteById("ProcessData",m.GetId())
		cnt+=1
	}
	fmt.Println("count unique mac...", rset.Cardinality(),"  cnt.. ", cnt)
	fmt.Println(time.Now().Sub(startTime))
}
