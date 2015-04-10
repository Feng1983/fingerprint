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
/*
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
*/
/*
type DwellProc struct{
	Mac	int64
	ST	int64
	Dwell	int64	
}*/

func getData()[]* knn2.ProcessData{
	var sp []*knn2.ProcessData
	if err:= zoom.NewQuery("ProcessData").Order("Timestamp").Scan(&sp);err!=nil{
            panic(err)
    	}
	return sp
}

func iter_process_dewell(dat []*Rssiample) []*DwellProc{
	//obj := NewPostgresqlStorage()
	//dat,_:= obj.GetSampleFromdb(0,"2015-03-17")
	var  ret []*DwellProc
	tmpMap  := make(map[int64][]int64)
	size:= len(dat)
	if size<=0{
		fmt.Println("no size")
		return nil
	}

	start,end:= dat[0],dat[size-1]
	var samples1= RIter(dat)
	it:=samples1.Iterator()
	for i:=start.Ts+2;i<=end.Ts+1;i+=2{
		mapr:=make(map[int64] *Info)
		for{
			val,ok := it()
			if !ok{
				break
			}
			if val.Ts >i {
				break
			}
			if v, exist := mapr[val.Dmac];!exist{
				tmp:=&Info{}
				rv:=make(map[int64]int)
				rv[val.Imac]=val.Rss
				tmp.x=rv
				mapr[int64(val.Dmac)]=tmp
				//fmt.Println(val.Rss)
			}else{
				(v.x)[val.Imac] =val.Rss
				mapr[val.Dmac]=v
			}
		}
		for s,t :=range mapr{
			//Features:t.x,Timestamp:int64(i),Mac:int64(s)
			if len(t.x) ==2{
				for _,v :=range t.x{
                   			 if v>-70 {
                        			//find value
						ar , exists:=tmpMap[int64(s)]
						if exists{
							ar=append(ar,int64(i))
							tmpMap[int64(s)] =ar	
						}else{
							tmpMap[int64(s)] =append(tmpMap[int64(s)],int64(i))
						}
                       	 			break
                    			 }
                		}
			}else if len(t.x)==3{
				tmpMap[int64(s)]=append(tmpMap[int64(s)],int64(i))
				
			}else{
				for _,v := range t.x{
					if v>-60{
						tmpMap[int64(s)]=append(tmpMap[int64(s)],int64(i))
						break
					}
				}
			}
		}
	}

	for k, v:=range tmpMap{
		sz:= len(v)
		fmt.Println(k, sz)
		st:= 0
		fmt.Print(k)
		for i:=1;i<sz;i++{
			if v[i]-v[i-1] > 600{
				if i-1-st ==0{
					fmt.Print(" || ",time.Unix(v[st],0).Format("2006-01-02 15:04:05"), " ",200)
					ret= append(ret, &DwellProc{ST:v[st],Mac:k, Dwell:200})
				}else{
					fmt.Print(" || ",time.Unix(v[st],0).Format("2006-01-02 15:04:05"), " ",v[i-1]-v[st])
					ret= append(ret, &DwellProc{ST:v[st],Mac:k, Dwell:v[i-1]-v[st]})	
				}
				st = i
			}	
		}
		if st==sz-1{
			fmt.Print(" || ",time.Unix(v[st],0).Format("2006-01-02 15:04:05"), " ",200)
			ret= append(ret, &DwellProc{ST:v[st],Mac:k, Dwell:200})
		}else{
			fmt.Print(" || ",time.Unix(v[st],0).Format("2006-01-02 15:04:05"), " ",v[sz-1]-v[st])
			ret= append(ret, &DwellProc{ST:v[st],Mac:k, Dwell:v[sz-1]-v[st]})
		}
		fmt.Println()
	}
	fmt.Println(len(tmpMap))
	return ret
	
}
func procdwell(date string , storeid int){
	rset := mapset.NewSet()
	startTime:= time.Now()
	obj := NewPostgresqlStorage()
        dat,_:= obj.GetSampleFromdb(storeid,date)
        savdat := iter_process_dewell(dat)
        fmt.Println(len(savdat))
        for _,d :=range savdat{
                rset.Add(d)
        }
        err:=obj.SaveDwellData(savdat ,storeid,  storeid, date)
	if err!=nil{
                fmt.Println(err)
        }
        fmt.Println(len(savdat),rset.Cardinality())
        fmt.Println(time.Now().Sub(startTime))
        defer obj.CloseORdb()
}
/*
func main(){
	//Init()
	rset := mapset.NewSet()
	startTime:= time.Now()

//	fmt.Println("count unique mac...", rset.Cardinality(),"  cnt.. ", cnt)

	obj := NewPostgresqlStorage()
        dat,_:= obj.GetSampleFromdb(1,"2015-03-18")
	savdat := iter_process_dewell(dat)
	fmt.Println(len(savdat))
	for _,d :=range savdat{
		rset.Add(d)
	}
	err:=obj.SaveDwellData(savdat ,1,  1, "2015-03-18")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(len(savdat),rset.Cardinality())
	fmt.Println(time.Now().Sub(startTime))
	defer obj.CloseORdb()
}*/
func main(){
	procdwell("2015-03-20",1)
}

