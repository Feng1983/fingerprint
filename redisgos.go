
package main

import(
	"fmt"
	"time"
	"github.com/garyburd/redigo/redis"
	"log"
)

func goes(){
	client, err := redis.DialTimeout("tcp",":59999",0,1*time.Second, 1*time.Second)
	if err!=nil{
		panic(nil)
	}
	
	size,_ :=client.Do("DBSIZE")
	fmt.Printf("db size is %d \n",size)
	_,err = client.Do("SET","user:user1",123)
	_,err = client.Do("SET","user:user0",456)
	_,err = client.Do("APPEND","user:user1",789)
	
	user0,_:=redis.Int(client.Do("GET","user:user0"))
	user1,_:=redis.Int(client.Do("GET","user:user1"))
	fmt.Printf("user0 is :%d, user1 is %d\n",user0,user1)
	defer client.Close()
}

func main(){
	log.Println("start...")
	goes()
}
