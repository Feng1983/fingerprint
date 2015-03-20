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
	callMap := make(map[int64][]int64) 
	for _, m :=range getData(){
		str_time := time.Unix(int64(m.Timestamp), 0).Format("2006-01-02 15:04:05")
		fmt.Println(m,str_time)
		//fmt.Println(m)
		rset.Add(m.Mac)
		//zoom.DeleteById("ProcessData",m.GetId())
		cnt+=1
		
		//find value;
		_,exists:=callMap[m.Mac]
		if exists{
				tmpp:= callMap[m.Mac]
				nz:=len(tmpp)
				if(m.Timestamp-tmpp[nz-2] - tmpp[nz-1] >300){
					callMap[m.Mac] = append(callMap[m.Mac],m.Timestamp)
					callMap[m.Mac] = append(callMap[m.Mac],0)
				}else{
					tmpp[nz-1]= m.Timestamp-callMap[m.Mac][nz-2]
				}
		}else{
			callMap[m.Mac] = append(callMap[m.Mac],int64(m.Timestamp))
			callMap[m.Mac] = append(callMap[m.Mac],0)
		}
	}
	for i,m := range callMap{
			fmt.Println(i,m)
			for i,v:= range m{
				if(i%2==0){
					fmt.Println( time.Unix(int64(v), 0).Format("2006-01-02 15:04:05"))
				}
			}
	}
	fmt.Println("count unique mac...", rset.Cardinality(),"  cnt.. ", cnt)
	fmt.Println(time.Now().Sub(startTime))
}
