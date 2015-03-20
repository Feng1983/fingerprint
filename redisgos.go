
package main

import(
	"fmt"
	"time"
	"github.com/garyburd/redigo/redis"
	"log"
)

var MAX_POOL_SIZE =20
var redisPool chan redis.Conn

func putRedis(conn redis.Conn){
	if redisPool ==nil{
		redisPool= make(chan redis.Conn,MAX_POOL_SIZE)
	}
	if len(redisPool) > MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPool <-conn
}

func InitRedis(network,address string) redis.Conn{
	redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
	if len(redisPool) == 0{
		go func(){
			for i:=0; i<MAX_POOL_SIZE/2;i++{
				c, err:= redis.Dial(network,address)
				if err!=nil{
					panic(err)
				}
				putRedis(c)
			}	
		}()
	}
	return <- redisPool
}
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

	c := InitRedis("tcp",":59999")
	fmt.Println(time.Now())
	starttime:= time.Now()

	var Success,Failure int
	//for i:=0;i<100000;i++{
		//if ok,_ := redis.Bool(c.Do("HSET", "payVerify:session", NewV1(), "aaaa"));ok{
		//	Success++
		//}else{
			Failure++
		//}
	//}
	m,_ :=redis.Values(c.Do("HGETALL","payVerify:session"))
	fmt.Println(len(m))
	//for i:=0;i<len(m);i++{
	//		s,_:= redis.String(m[i],nil)
	//		fmt.Println(s)
	//}
	fmt.Println(time.Now())
	fmt.Println("cost: ",time.Now().Sub(starttime),"success: ",Success,"failure: ",Failure)
}
