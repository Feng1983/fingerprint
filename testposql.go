mackage main

import (
	"database/sql"
	"fmt"
	//_ "github.com/bmizerany/pq"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/albrow/zoom"
	"log"
	"time"
	"errors"
	"knn2"
)
var (
	au = []int{5,6,8,}
	ch = []int{9,12,13,}
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
type PostgresqlStorage{
	DB	*gorm.DB	
}
//init sqlStorage
func NewPostgresqlStorage() *PostgresqlStorage{
	
}
func getTime(tt string) int64 {
    the_time, err := time.Parse("2006-01-02 15:04:05", tt)
    if err != nil {
        unix_time := the_time.Unix()
        fmt.Println(unix_time)
    }
    return the_time.Unix()
}

func (b SampleData) TableName() string {
	return "samples"
}

func openDB() *sql.DB {
	//db, err := sql.Open("postgres", "host=192.168.0.8 user=postgres password=111111 dbname=new_structure01 sslmode=disable")
	//db, err := sql.Open("postgres", "host=123.57.254.158 user=postgres password=111111 dbname=ailink_wifi sslmode=disable")
	db, err := sql.Open("postgres", "host=123.57.254.158 user=postgres password=111111 dbname=ailink_wifi sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func closeDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return err
}

func testdb(db *sql.DB) error {
	//db := openDB()
	//defer closeDB(db)

	rows, err := db.Query("select * from samples limit 10")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		//var id, ts,inframac,devmac,rssi, freq int
		var id int32
		var ts, devmac int64
		var inframac, rssi, freq int16

		err := rows.Scan(&id, &ts, &inframac, &devmac, &rssi, &freq)
		if err != nil {
			log.Fatal(err)
		}
		/*fmt.Println("id:", id)
		fmt.Println("timestamp:", ts)
		fmt.Println("infrastramac:", inframac)
		fmt.Println("devicemac: ", devmac)
		fmt.Println("rssi:", rssi)
		fmt.Println("freq:", freq)*/
		fmt.Println(id, ts, inframac, devmac, rssi, freq)
	}

	defer rows.Close()
	return nil
}

func ormInit() *gorm.DB {
	//db, err := gorm.Open("postgres", "host=192.168.0.8 user=postgres password=111111 dbname=new_structure01 sslmode=disable")
	db, err := gorm.Open("postgres", "host=123.57.254.158 user=postgres password=111111 dbname=ailink_wifi sslmode=disable")
	if err != nil {
		fmt.Println("connection error...", err)
	}
	dbsql := db.DB()
	dbsql.Ping()
	//dbsql.Close()
	return &db
}
func closeORdb(db *gorm.DB) error {
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
func fectchItem(db *gorm.DB,params interface{},dd1 ,dd2 int64){

	var users [] *SampleData
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
}

func loadFingers(loc int, db *gorm.DB){
	var findata []*knn2.FingerOri
	var id int
	res,_ := zoom.NewQuery("Finger2").Filter("Id= ",loc).Count()
	if res !=0{
	 return
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
		zoom.save(fin)
	}
}

func saveDate(db *gorm.DB)error{
	return nil
}

func SaveFingerData(db* gorm.DB, datas[]* knn2.ProcessData, userid storeid int) error{
	var err error
	if len(datas)==0{
		err.New("None data found")
		return err
	}
	for _, d:=range datas{
		ds:= FinMatchData{
			StoreId:storeid,
			UserId:userid,
			Ts:d.Timestamp,
			Mac:d.Mac,
			X: d.X,
			Y: d.Y,		
		}	
		db.Create(&ds)
	}
	
	return nil
}

func RedisInit(){
    conf:= &zoom.Configuration{
        Address:"localhost:59999",
        Network:"tcp",
    }
    zoom.Init(conf)
    //zoom.Register(&Person{})
    //zoom.Register(&Rssiample{})
    zoom.Register(&knn2.Finger{})
    zoom.Register(&knn2.ProcessData{})
	zoom.Register(&knn2.Finger2{})
}


func main() {
	var dbs *sql.DB
	dbs = openDB()
	testdb(dbs)
	closeDB(dbs)
	db := ormInit()
	//doIter(db)
	dd1:=getTime("2015-03-17 00:00:00")
	dd2:=getTime("2015-03-18 00:00:00")
	fectchItem(db,[]int{5,6,8},dd1,dd2)
	loadFingers(1,db)
	closeORdb(db)
}
