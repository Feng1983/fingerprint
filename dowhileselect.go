
package main

import(
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const NUM_OF_QUIT int =100
func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := make(chan bool)
	receive_chan := make(chan chan bool)
	finish:= make(chan bool)
	
	for i:=0; i<NUM_OF_QUIT; i++{
		go do_while_select(i,receive_chan,finish)
	}

	go handle_exit(done, receive_chan,finish)
	<- done
	os.Exit(0)
}

func do_while_select(num int, rece chan chan bool, done chan bool){
	quit:= make(chan bool)
	rece <- quit
	for{
		select{
		case <- quit:
			done <- true
			runtime.Goexit()
		default:
			log.Println("the",num,"is running")
		}
	}
}

func handle_exit(done chan bool, receive_chan chan chan bool, finish chan bool){
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT,syscall.SIGTERM)
	chan_slice:= make([]chan bool,0)
	for{
		select{
		case <- sigs:
			for _,v := range chan_slice{
				v <- true
			}
