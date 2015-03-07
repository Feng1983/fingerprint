package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/bmizerany/pq"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func getTime(tt string) int64 {
	the_time, err := time.Parse("2006-01-02 15:04:05", tt)
	if err == nil {
		unix_time := the_time.Unix()
		fmt.Println(unix_time)
	}
	return the_time.Unix()
}

type SampleData struct {
	Id       int32
	Ts       int64 `gorm:"column:timestamp"`
	Inframac int16 `gorm:"column:infra_mac"`
	Devmac   int64 `gorm:"column:device_mac"`
	Rssi     int16
	Freq     int16 `gorm:"column:frequency"`
}

func (b SampleData) TableName() string {
	return "samples"
}

func openDB() *sql.DB {
	db, err := sql.Open("postgres", "host=192.168.0.8 user=postgres password=111111 dbname=new_structure01 sslmode=disable")
	//db, err := sql.Open("postgres", "host=123.57.254.158 user=postgres password=111111 dbname=ailink_wifi sslmode=disable")
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
	db, err := gorm.Open("postgres", "host=192.168.0.8 user=postgres password=111111 dbname=new_structure01 sslmode=disable")
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
	dd := getTime("2015-03-05 12:44:33")
	db.Where("timestamp >= ? ", dd).Limit(20).Find(&users2)
	for _, d := range users2 {
		fmt.Println(d)
		fmt.Println(time.Unix(d.Ts, 0).Format("2006-01-02 15:04:05"))
	}
	var psize int
	dd2 := getTime("2015-03-06 12:44:33")
	db.Where("timestamp >= ? and timestamp <?", dd, dd2).Count(&psize)
	fmt.Println(psize)

}

func main() {
	//db,err  := sql.Open("postgres","user=postgres dbname=new_structure01 sslmode=verify-full")
	//defer db.close()
	var dbs *sql.DB
	dbs = openDB()
	testdb(dbs)
	closeDB(dbs)
	db := ormInit()
	doIter(db)
	closeORdb(db)
}
