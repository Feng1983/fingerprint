
package main

import (
    "log"
    //"github.com/boltdb/bolt"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

var (
	host = flag.String("host","0.0.0.0","http server name")
	port = flag.Int("port",19890,"http server port")
	dict = flag.String("dict","","dict string")
	staticFolder = flag.String("static_folder", "static","")
)

type JsonResponse struct{
	Segments []*Segment
}
type Segment struct{
	Text string
	Pos  string
}

func jsonRpcServer(w http.ResponseWriter, req *http.Request){
	text := req.URL.Query().Get("req")
	if text == "" {
		text = req.PostFormValue("text")
	}
	
	ss := []*Segment{}

	ss = append(ss,&Segment{Text:text,Pos:"222"})
	response,_  := json.Marshal(JsonResponse{Segments:ss})
	w.Header().Set("Content-Type","application/json")
	io.WriteString(w,string(response))
}

func main() {
    // Open the my.db data file in your current directory.
    // It will be created if it doesn't exist.
    //db, err := bolt.Open("my.db", 0600, nil)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //defer db.Close()
