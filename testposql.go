package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/bmizerany/pq"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

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
	var users []*SampleData
	db = db.Table("samples")
	//db.AutoMigrate(&SampleData{})
	db.Limit(10).Find(&users)
	for i, d := range users {
		fmt.Println(i, d)
		fmt.Println(d.Ts)
	}
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
