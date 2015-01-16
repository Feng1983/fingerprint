
package main
import(
	"fmt"
	"os/user"
	"os"
	"os/signal"
)


var d = make(chan int, 10)
var ds string

func f(){
	ds="hello, world!"
	d<- 0
} 

func main(){
	var err user.UnknownUserIdError
	err=2
	
	fmt.Println(err.Error())
	user,_ := user.Current()
	fmt.Println(user.Uid)
	fmt.Println(user.Gid)
	fmt.Println(user.Name)
	fmt.Println(user.Username)
	fmt.Println(user.HomeDir)

	c:= make(chan os.Signal,1)
	signal.Notify(c,os.Interrupt,os.Kill)

	s:=  <-c
	fmt.Println("get signal .." ,s)

	go f()
	<- d
	fmt.Println(ds)
}
