package main

import (
	//"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/albrow/zoom"
	"time"
	"errors"
	"knn2"
)
var (
	au = []int{5,6,8,}	//1
	ch = []int{9,12,13,} //0
	ErrNoPqSQLConn     = errors.New("can't get a mysql db")
	ErrNoData		   = errors.New("NO data found.")
	ErrExist		   = errors.New("data exists ,return.")
)


type SampleData struct {
	Id       int32
	Ts       int64 `gorm:"column:timestamp"`
	Inframac int16 `gorm:"column:infra_mac"`
	Devmac   int64 `gorm:"column:device_mac"`
	Rssi     int16
	Freq     int16 `gorm:"column:frequency"`
}

type FinMatchData struct{
	StoreId	int32 `gorm:"column:store_id"`
	UserId	int32 `gorm:"column:user_id"`
	Ts		int64 
	Mac		int64
	X		float64
	Y		float64 
}

type PostgresqlStorage struct{
	DB	*gorm.DB	
}

//init sqlStorage
func NewPostgresqlStorage() *PostgresqlStorage{
	db, err := gorm.Open("postgres", "host=123.57.254.158 user=postgres password=111111 dbname=ailink_wifi sslmode=disable")
    if err != nil {
        fmt.Println("connection error...", err)
    }
	RedisInit()
    return &PostgresqlStorage{DB:&db}
}

func (b SampleData) TableName() string {
	return "samples"
}



func (pqsDb *PostgresqlStorage)CloseORdb() error {
	db:= pqsDb.DB
	if db==nil{
		return ErrNoPqSQLConn
	}
	err := db.Close()
	return err
}

func doIter(db *gorm.DB) {
	var users, users2 []*SampleData
	db = db.Table("samples")
	//db.AutoMigrate(&SampleData{})
	db.Limit(10).Find(&users)
	for i, d := range users {
		fmt.Println(i, d)
		//fmt.Println(d.Ts)
		fmt.Println(time.Unix(d.Ts, 0).Format("2006-01-02 15:04:05"))
	}
	dd := getTime("2015-03-10 12:44:33")
	db.Where("timestamp >= ? ", dd).Limit(20).Find(&users2)
	for _, d := range users2 {
		fmt.Println(d)
		fmt.Println(time.Unix(d.Ts, 0).Format("2006-01-02 15:04:05"))
	}
	var psize int
	dd2 := getTime("2015-03-12 12:44:33")
	db.Where("timestamp >= ? and timestamp <?", dd, dd2).Count(&psize)
	fmt.Println(psize)

}
func (pqsDb *PostgresqlStorage) FectchItem(params interface{},dd1 ,dd2 int64) error{
	var users [] *SampleData
	db := pqsDb.DB
	if db==nil{
		return ErrNoPqSQLConn 
	}
	db = db.Table("samples")
	//rows,err:= db.Raw("select * from samples where infra_mac in (?)",param)
	//defer rows.Close()
	db.Where("infra_mac in (?)",params).Where("timestamp >= ? and timestamp <?",dd1,dd2).Limit(200).Find(&users)
	for _ , d :=range users{
		fmt.Println(d,time.Unix(d.Ts, 0).Format("2006-01-02 15:04:05"))
	}
	var psize int
	db.Where("infra_mac in (?)",params).Where("timestamp >= ? and timestamp <?",dd1,dd2).Count(&psize)
	fmt.Println(psize)
	return nil
}

func (pqsDb  *PostgresqlStorage) LoadFingers(loc int) error{
	var findata []*knn2.FingerOri
	var id int
	db := pqsDb.DB
	if db==nil{
		return ErrNoPqSQLConn
	}
	res,_ := zoom.NewQuery("Finger2").Filter("Id= ",loc).Count()
	if res !=0{
	 return ErrNoData
	}
	if loc ==1{
		db=db.Table("au_finger_data")
		id=1
	}else{
		db=db.Table("cn_finger_data")
		id=0
	}
	db.Find(&findata)

	for _, d :=range findata{
		//fmt.Println(d)
		fs:=make(map[int64]float64)
		fs[int64(d.Ap1)]= d.Freq1
		fs[int64(d.Ap2)]= d.Freq2
		fs[int64(d.Ap3)]= d.Freq3
		
		fin:=&knn2.Finger2{Id:id,Label:d.Id,X:d.X,Y:d.Y,Feature:fs}
		fmt.Println(fin)
		zoom.Save(fin)
	}
	return nil
}

func saveDate(db *gorm.DB)error{
	return nil
}

func (pqsDb  *PostgresqlStorage) SaveFingerData(datas[]* knn2.ProcessData, userid, storeid int) error{
	db := pqsDb.DB
	if db==nil{
		return ErrNoPqSQLConn
	}
	if len(datas)==0{
		return  ErrNoData
	}
	for _, d:=range datas{
		ds:= FinMatchData{
			StoreId:int32(storeid),
			UserId:int32(userid),
			Ts:d.Timestamp,
			Mac:d.Mac,
			X: d.X,
			Y: d.Y,		
		}	
		db.Create(&ds)
	}
	
	return nil
}
func (pqsDb  *PostgresqlStorage) GetSampleFromdb(id int, dd string)([]*Rssiample,error){
	db := pqsDb.DB
    	if db==nil{
        	return nil,ErrNoPqSQLConn 
    	}
    	db = db.Table("samples")
    	dd1 := getTime(dd)
    	dd2 := Dateplus(dd,1)
    	var param []int
	var users []*SampleData
    	if id==0{
		param = ch
    	}else {
		param = au
    	}
    	db.Where("infra_mac in (?)",param).Where("timestamp >= ? and timestamp <?",dd1,dd2).Order("timestamp").Find(&users)
    	var ret []*Rssiample
    	for _, d:=range users{
		mdat := &Rssiample{Ts:int(d.Ts), Imac:int64(d.Inframac), Dmac:d.Devmac, Rss:int(d.Rssi),Frq:int(d.Freq),Id:int(d.Id)}
		fmt.Println(mdat)
		ret = append(ret, mdat)	
    	}
    	return ret,nil
	
}
func RedisInit(){
    conf:= &zoom.Configuration{
        //Address:"localhost:59999",
	Address:"localhost:6379",
        Network:"tcp",
    }
    zoom.Init(conf)
    zoom.Register(&knn2.Finger{})
    zoom.Register(&knn2.ProcessData{})
    zoom.Register(&knn2.Finger2{})
}


func main() {
	//var dbs *sql.DB
	//dbs = openDB()
	//testdb(dbs)
	//closeDB(dbs)
	//db := ormInit()
	//doIter(db)
	//dd1:=getTime("2015-03-17 00:00:00")
	//dd2:=getTime("2015-03-18 00:00:00")
	//fectchItem(db,[]int{5,6,8},dd1,dd2)
	//loadFingers(1,db)
	//closeORdb(db)
	obj := NewPostgresqlStorage()
	//db:= obj.db
	obj.GetSampleFromdb(0,"2015-03-17")
	obj.CloseORdb()
}
