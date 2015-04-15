package main
import (
    log "github.com/alecthomas/log4go"
    //"github.com/garyburd/redigo/redis"
    "runtime"
)

func main(){
    log.Info("message ver: \"%s\" start")
    if err := InitConfig(); err != nil {
        panic(err)
    }
    // Set max routine
    runtime.GOMAXPROCS(Conf.MaxProc)
    // init log                                                                                                                                          
    //log.LoadConfiguration(Conf.Log)
    //log.Info(Conf.RedisSource,Conf.PostgreSQLSource)
    //log.Info("message ver: \"%s\" a2 start")

    //obj:= NewRedisStorage()
    InitLog()
    InitRedis()
    //conn := RedisInst.getConnByNode("node1")
    //values, err := redis.String(conn.Do("INFO"))
    //if err != nil {
   //		log.Error("conn.Do info", err)
    //}
    mv := RedisInst.GetRedisRssi(0,1428706745,1428755679,[]int{15})
	
    for _,v :=range mv{
	log.Info(v)
    }
    //log.Info(values)
    //conn.Close()
    defer func(){
	log.Close()
	RedisInst.Clear()
    }()
    // start pprof http
    log.Info("message stop")
}

func InitLog(){
    log.Info("message ver: \"%s\" start")
    if err := InitConfig(); err != nil {
        panic(err)
    }
    // Set max routine
    runtime.GOMAXPROCS(Conf.MaxProc)
    // init log                                                                                                                                          
    log.LoadConfiguration(Conf.Log)
    log.Info(Conf.RedisSource,Conf.PostgreSQLSource)
    log.Info("message ver: \"%s\" a2 start")

    //defer log.Close()
    // start pprof http
    //log.Info("message stop")
}
