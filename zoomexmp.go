package main

import(
	"github.com/albrow/zoom"
	"github.com/deckarep/golang-set"
	"fmt"
	"strconv"
	"strings"
	"time"
	"os"
	//"io/ioutil"
	//"encoding/csv"
	"bufio"
	"knn2"
)

type Person struct{
	Name string
	Age  int
	zoom.DefaultData
}

type Rssiample struct{
	Id int
	Ts int		`zoom:"index"`
	Imac int64	`zoom:"index"`
    Dmac int64	`zoom:"index"`
	Rss  int
	Frq  int
	zoom.DefaultData
} 

func load(path string )error{
	file, err:= os.Open(path)
	if err!=nil{
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
    //if err:= zoom.Register(&Rssiample{});err!=nil{
    //    fmt.Println("regiser Error")
    //}

	for scanner.Scan(){
		line := strings.Replace(scanner.Text(), " ", "\t", -1)
		tks := strings.Split(line, ";")
		//fmt.Println(line)
		//fmt.Println(tks)
		if tks[0]=="id"{
			continue
		}
		fmt.Println(tks[0],tks[1])
		//sample=&RSSSample{}
		id1,_:=strconv.Atoi(tks[0])
		if err!=nil{
			fmt.Println(err)
		}
		ts1,err:=strconv.Atoi(tks[1])
		if err!=nil{
			
		}
		imac1,err:= strconv.ParseInt(tks[2],10,64)
		if err!=nil{
		}
		dmac1,err:= strconv.ParseInt(tks[3],10,64)
		if err!=nil{}
		rss1,_:= strconv.Atoi(tks[4])
		frq1,_:=strconv.Atoi(tks[5])
		sample:=&Rssiample{Id:id1,Ts:ts1,Imac:imac1,Dmac:dmac1,Rss:rss1,Frq:frq1}
		fmt.Println(sample)
		zoom.Save(sample)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	//zoom.Close()
	return nil
}


func loadfinger(path string) error{
	file, err:= os.Open(path)
    if err!=nil{
        return err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := strings.Replace(scanner.Text(), " ", "\t", -1)
        tks := strings.Split(line, ";")
		if len(tks)!=9{
			continue
		}		
		lbel,_:=strconv.Atoi(tks[0])
		x,_:=strconv.ParseFloat(tks[1],64)
		y,_:=strconv.ParseFloat(tks[2],64)
		fs:=make(map[int64]float64)
		for i:=3;i<len(tks);i++{
			rssi,_:= strconv.ParseFloat(tks[i],64)
			dmac,_:= strconv.ParseInt(tks[i+1],10,64)
			fs[dmac]=rssi
			//fmt.Println(i,rssi,dmac)
			i=i+1
		}
		fin:=&knn2.Finger{Label:lbel,X:x,Y:y,Feature:fs}
		zoom.Save(fin)
	}
	if scanner.Err() !=nil{
		return scanner.Err()
	}
	return nil
}
func Init(){
	conf:= &zoom.Configuration{
        Address:"localhost:59999",
        Network:"tcp",
    }
	zoom.Init(conf)
	zoom.Register(&Person{})
	zoom.Register(&Rssiample{})
	zoom.Register(&knn2.Finger{})
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
func getFingers()[]* knn2.Finger{
	var rs[]* knn2.Finger
	if err:=zoom.NewQuery("Finger").Order("Label").Scan(&rs);err!=nil{
		panic(err)
	}
	return rs
}
type RdataSet struct{
	Samples []* knn2.MapBaseSample
	max_lable  int
}
/*type MapBaseSample struct{
	Features    map[int64]int
	Label		int
	Prediction  float64
	Timestamp   int64
	Mac			int64
}*/

type RIter []*Rssiample
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
type Info struct{
	x map[int64]int
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


func main(){
	Init()
	starttime:=time.Now()

	p := &Person{Name: "Alice", Age: 27}
	//for i:=0;i<1000000;i++{
	//	tmp:= []string{"name",strconv.Itoa(i)}
	//	zoom.Save(&Person{Name:strings.Join(tmp,"|"),Age:i+2})
	//}

	result,_ := zoom.FindById("Person",p.Id)
	person, ok := result.(*Person)
	if !ok{
		fmt.Println("Error Person type")
	}
	fmt.Println(result,person.Name,person.Age)
	var persons  []*Person
	if err:=zoom.NewQuery("Person").Scan(&persons);err!=nil{
		panic(err)
	}
	for _,p:= range persons{
		//p.Age =p.Age+20000
		//zoom.Save(p)
		fmt.Println(p.Age,p.Name)
	}
	fmt.Println(len(persons))
	//for i,rs := range results{
	//	if ps, ok := rs.(*Person); !ok {
	//	}
	//	fmt.Println(ps.Age, ps.Name)
	//}
	//fmt.Println(time.Now().Sub(starttime))
	defer zoom.Close()
	//err := loadcsv("./tt.csv")
	//err := load("tt2.csv")
	//if err!=nil{
	//	fmt.Println(err)
	//}
	sp:=getSamples(int64(-1))
	pp:=etldata(sp)
	cnt,rcnt:=0,0
	rset := mapset.NewSet()
	var proc_data []* knn2.MapBaseSample
	for _,m:=range pp{
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
	
	//fmt.Println(rset.String())
	/*sp:=getSamples(int64(13))
	fmt.Println(len(sp))
	res:= mapset.NewSet()
	for _, rss:= range sp{
		fmt.Println(rss)
		//zoom.DeleteById("Rssiample",rss.GetId())
		str_time := time.Unix(int64(rss.Ts), 0).Format("2006-01-02 15:04:05")
		fmt.Println(str_time)
		res.Add(rss.Dmac)
	}
	if len(sp)>0{
		start,end:=sp[0],sp[len(sp)-1]
		endTs:=time.Unix(int64(end.Ts), 0).Format("2006-01-02 15:04:05")
		stTs:=time.Unix(int64(start.Ts), 0).Format("2006-01-02 15:04:05")
		fmt.Println(len(sp),end.Ts-start.Ts,endTs, stTs)
	}
	fmt.Println(len(sp))
	fmt.Println(res.Cardinality())
	fmt.Println(res.String())
	//for i:=start.Ts;i<=end.Ts;i+=2{
	//		fmt.Println(time.Unix(int64(i), 0).Format("2006-01-02 15:04:05"))
	//}
	*/
	/*err:=loadfinger("finger.csv")
	if err!=nil{
		fmt.Println(err)
	}*/
	//ii:=0
	//for _,v:=range getFingers(){
	//	fmt.Println(v)
	//	ii++
		//zoom.DeleteById("Finger",v.GetId())
	//}
	//fmt.Println(ii)
	params:=make(map[string]string)
	params["k"]="5"
	obj:=&knn2.KNN{}
	obj.Init(params)
	fingdata:= getFingers()
	obj.Train2(fingdata)
	for _,m :=range proc_data{
		if len(m.Features)>=2{
            str_time := time.Unix(int64(m.Timestamp), 0).Format("2006-01-02 15:04:05")
            fmt.Println(m,str_time)
			mx,my:= obj.Predict2(m)
			fmt.Println(mx,"y... ",my)
			prdata:= &knn2.ProcessData{Timestamp:m.Timestamp, Mac:m.Mac,X:mx, Y:my}
			fmt.Println(prdata)
			//zoom.Save(prdata)
        }
	}
	fmt.Println(time.Now().Sub(starttime))
	defer zoom.Close()
}
