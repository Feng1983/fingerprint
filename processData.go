package main

import(
	"github.com/albrow/zoom"
	"github.com/deckarep/golang-set"
	"fmt"
	//"strconv"
	//"strings"
	"time"
	//"os"
	//"bufio"
	"knn2"
)

type Rssiample struct{
	Id int
	Ts int		`zoom:"index"`
	Imac int64	`zoom:"index"`
    	Dmac int64	`zoom:"index"`
	Rss  int
	Frq  int
	zoom.DefaultData
} 

type RdataSet struct{
    Samples []* knn2.MapBaseSample
    max_lable  int
}

type Info struct{
    x map[int64]int
}

type RIter []*Rssiample

func Init(){
	conf:= &zoom.Configuration{
        	Address:"localhost:59999",
        	Network:"tcp",
        }
	zoom.Init(conf)
	//zoom.Register(&Person{})
	zoom.Register(&Rssiample{})
	zoom.Register(&knn2.Finger{})
	zoom.Register(&knn2.Finger2{})
	zoom.Register(&knn2.ProcessData{})
}

func getSamples(imac int64) []*Rssiample{
	var sp []*Rssiample
	if(imac==int64(-1)){
		if err:= zoom.NewQuery("Rssiample").Order("Ts").Scan(&sp);err!=nil{
            panic(err)
        }
	}else{
		if err:= zoom.NewQuery("Rssiample").Order("Ts").Filter("Imac =",imac).Scan(&sp);err!=nil{
			panic(err)
		}
	}
	return sp
}

func getFingers(id int)[]* knn2.Finger2{
	var rs[]* knn2.Finger2
	if err:=zoom.NewQuery("Finger2").Filter("Aid =",id).Order("Label").Scan(&rs);err!=nil{
		panic(err)
	}
	return rs
}

func (i RIter)Iterator() func()( *Rssiample,bool){
	index:=0
	return func()(val *Rssiample, ok bool){
			if index>=len(i){
				return
			}
			val, ok = i[index],true
			index++
			return
	}
}


func etldata(samples []*Rssiample)[]* knn2.MapBaseSample{
	size:= len(samples)
	if size<=0{
		return nil
	}
	start,end:= samples[0],samples[size-1]
	var samples1= RIter(samples)
	it:=samples1.Iterator()
	var ret []* knn2.MapBaseSample
	str_time := time.Unix(int64(start.Ts), 0).Format("2006-01-02 15:04:05")
	str2     := time.Unix(int64(end.Ts), 0).Format("2006-01-02 15:04:05")
	fmt.Println(str_time,str2)
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
			ret = append(ret,&knn2.MapBaseSample{Features:t.x,Timestamp:int64(i),Mac:int64(s)})
		}
	}
	return ret
}


//func procETLdata(id int, dd string , obj *PostgresqlStorage) []*knn2.MapBaseSample{
func procETLdata(id int, dd string , obj []*Rssiample) []*knn2.MapBaseSample{
	var proc_data []* knn2.MapBaseSample
	//sp := getSampleFromdb(id, dd)
	
	//sp := []*Rssiample{}
	//sp ,err := obj.GetSampleFromdb(id, dd)
	var sp []*Rssiample
	sp=obj
	fmt.Println("value",len(sp))
	pp := etldata(sp)
	cnt,rcnt:=0,0
	rset := mapset.NewSet()
	for _ , m:=range pp{
		if len(m.Features)>=2{
            //str_time := time.Unix(int64(m.Timestamp), 0).Format("2006-01-02 15:04:05")
            //fmt.Println(m,str_time)
            rcnt+=1
            if len(m.Features)==2{
                for _,v :=range m.Features{
                    if v>-70 {
                        rset.Add(m.Mac)
                        //fmt.Println(k,"...",v)
                        proc_data=append(proc_data,m)
                        break
                    }
                }
            }else if len(m.Features)==3{
                rset.Add(m.Mac)
                proc_data=append(proc_data, m)
            }else{
                rset.Add(m.Mac)
            }
            //rset.Add(m.Mac)
        }
        for _,v:=range m.Features{
            if v>-60 {
                rset.Add(m.Mac)
            }
        }
        //rset.Add(m.Mac)
        cnt+=1
	}
	fmt.Println(cnt,"num...",rcnt,"unique mac",rset.Cardinality())
	return proc_data
}

func processFingerDataById(id int,proc_data []*knn2.MapBaseSample) []*knn2.ProcessData{
    params:=make(map[string]string)
    params["k"]="5"
    obj:=&knn2.KNN{}
    obj.Init(params)
    fingdata:= getFingers(id)//getFingers
    obj.Train2(fingdata)
    var ret []*knn2.ProcessData

    for _,m :=range proc_data{
        if len(m.Features)>=2{
            str_time := time.Unix(int64(m.Timestamp), 0).Format("2006-01-02 15:04:05")
            fmt.Println(m,str_time)
            mx,my:= obj.Predict2(m)
            fmt.Println(mx,"y... ",my)
            prdata:= &knn2.ProcessData{Timestamp:m.Timestamp, Mac:m.Mac,X:mx, Y:my}
            fmt.Println(prdata)
            //zoom.Save(prdata)
	    ret = append(ret,prdata)
	    fmt.Println(prdata)
        }
    }
    return ret
}

func cron_proc(date string, storeid int) {
	starttime := time.Now()
	obj := NewPostgresqlStorage()
	obj.LoadFingers(storeid)
	dat,_:= obj.GetSampleFromdb(storeid,date)	


	proc_dat := procETLdata(storeid,date,dat)
        res:= processFingerDataById(storeid,proc_dat)
	obj.SaveFingerData(res,storeid,storeid)
	fmt.Println(time.Now().Sub(starttime))
        obj.CloseORdb()
        defer zoom.Close()
}
/*
func main(){
	Init()
	starttime:=time.Now()
	obj := NewPostgresqlStorage()	
	//obj.GetSampleFromdb(0,"2015-03-17")
	obj.LoadFingers(0)
	dat,_:= obj.GetSampleFromdb(0,"2015-03-17")
	
	//obj.IterFinData(0)
	//obj.DelteById(0)
	proc_dat := procETLdata(0,"2015-3-17",dat)
	res:= processFingerDataById(0,proc_dat)
	//fmt.Println(len(proc_dat))
	obj.SaveFingerData(res,0,0)
	fmt.Println(time.Now().Sub(starttime))
	obj.CloseORdb()
	defer zoom.Close()
}*/
func main(){
	Init()
	//cron_proc("2015-03-17",1)
	obj := NewPostgresqlStorage()
	obj.GetDeviceMap()
	obj.CloseORdb()
        defer zoom.Close()
}
