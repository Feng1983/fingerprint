package main

import (
	log	"github.com/alecthomas/log4go" 
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	//"time"
	"github.com/albrow/zoom"
	"github.com/Terry-Mao/gopush-cluster/ketama"
	"knn2"
)
const (
	ketamaBase       = 255
)

var (
	RedisNoConnErr       = errors.New("can't get a redis conn")
	redisProtocolSpliter = ":"
	RedisInst  *RedisStorage
)

// RedisMessage struct encoding the composite info.



type RedisStorage struct {
	pool  map[string] *redis.Pool
	ring  *ketama.HashRing
}

type RedisMac struct{
	Key  string
	Value string
	Ts   int64
	Expire int64
}


// NewRedis initialize the redis pool and consistency hash ring.
func NewRedisStorage() *RedisStorage {
	redisPool := map[string]*redis.Pool{}
	ring := ketama.NewRing(ketamaBase)
	for n, addr := range Conf.RedisSource {
		nw := strings.Split(n, ":")
		if len(nw) != 2 {
			err := errors.New("node config error, it's nodeN:W")
			log.Error("strings.Split(\"%s\", :) failed (%v)", n, err)
			panic(err)
		}
		w, err := strconv.Atoi(nw[1])
		if err != nil {
			log.Error("strconv.Atoi(\"%s\") failed (%v)", nw[1], err)
			panic(err)
		}
		// get protocol and addr
		pw := strings.Split(addr, "@")
		if len(pw) != 2 {
			log.Error("strings.Split(\"%s\", \"%s\") failed (%v)", addr, redisProtocolSpliter, err)
			panic(fmt.Sprintf("config redis.source node:\"%s\" format error", addr))
		}
		tmpProto := pw[0]
		tmpAddr := pw[1]
		// WARN: closures use
		//tmp := addr
		log.Info(tmpProto," | ",tmpAddr)
		redisPool[nw[0]] = &redis.Pool{
			MaxIdle:     Conf.RedisMaxIdle,
			MaxActive:   Conf.RedisMaxActive,
			IdleTimeout: Conf.RedisIdleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", tmpAddr)
				if err != nil {
					log.Error("redis.Dial(\"tcp\", \"%s\") error(%v)", tmpAddr, err)
					return nil, err
				}
				return conn, err
			},
		}
		// add node to ketama hash
		ring.AddNode(nw[0], w)
	}
	ring.Bake()
	s := &RedisStorage{pool: redisPool, ring: ring}
	return s
}

//Save
func (s *RedisStorage) SaveRedisMac(dat []*RedisMac,node string)error{
	if len(dat)==0{
		log.Error("no data dump to redis")
		return nil
	}
	conn := s.getConnByNode(node)
        if conn == nil {
                return RedisNoConnErr
        }
        defer conn.Close()
	for _, iter:=range dat{
		if err := conn.Send("ZADD", iter.Key, iter.Ts, iter.Value); err != nil {
			//key storeid:mac:ts
			log.Error("conn.Send(\"ZADD\", \"%s\", %d, \"%s\") error(%v)", iter.Key, iter.Ts, iter.Value, err)
			fmt.Println(err)
                	return err
        	}
	}
	if err := conn.Flush(); err != nil {
                log.Error("conn.Flush() error(%v)", err)
                return err
        }
	return nil
} 

func (s *RedisStorage) SaveRssi(dat []*Rssiample) error{
	if len(dat)==0{
		return nil
	}
	for _, v:=range dat{
		log.Info("save...", v)
		//log.Info("insert ... ",v.Imac," ",v.Dmac," ",v.Ts," ",v.Rss," ",v.Frq," ",v)
		zoom.Save(v)
	}
	return nil
}

func (s *RedisStorage) GetRedisRssi(storeid int, startTs,endTs int64,ids []int)[]*Rssiample{
	
	sp := []* Rssiample{}
	fmt.Println("ids is .. ",ids)
	log.Info("ids is .. ",ids)

	zmquery := zoom.NewQuery("Rssiample").Order("Ts").Filter("Ts >=",int(startTs))//.Filter("Ts <",int(endTs))
	for _,item:=range ids{
		zmquery.Filter("Imac =",int64(item))
		fmt.Println("Imac is ",item," query ",zmquery) 
		log.Info("Imac is ",item)
	}
	cnt,_:=zmquery.Count()
	fmt.Println(" query ",zmquery," ",cnt)
        
	//pps,_ := zoom.FindById("Rssiample","SMlyWvpyqHPqjQE1nmn50q")
	//fmt.Println(pps)
	//if ms,ok :=pps.(*Rssiample);!ok{
		
	//}else{
	//	fmt.Println(ms)
	//}
	//return sp
        //if err:= zoom.NewQuery("Rssiample").Order("Ts").Filter("").Filter().Scan(&sp);err!=nil{
	if err:= zmquery.Scan(&sp);err!=nil{
            	panic(err)
        }
	//results,err := zmquery.Run()
	//sp = results.([]*Rssiample)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//for _,rs:=range sp{
		//if person, ok := rs.(*Rssiample); !ok {
	//		fmt.Println(rs," ",rs.Imac,rs.Ts)
		//}
	//}
	fmt.Println(len(sp))
	return sp
}

func(s *RedisStorage)SaveFingerData(dat []*knn2.ProcessData , storeid int, node string) error{
	conn:= s.getConnByNode(node)
	
	if len(dat)==0 {
		return nil
	}
	
	if conn == nil {
		return RedisNoConnErr
	}
	defer conn.Close()
	zkey:= strconv.Itoa(storeid)+":fingersort"
	var key string
	for _, iter:=range dat{
		key = strconv.Itoa(storeid)+":"+strconv.FormatInt(iter.Mac,10)+":"+strconv.FormatInt(iter.Timestamp,10)
		
		//conn.Do("HMSET",key,"Mac",int64(iter.Mac), "storeid",storeid,"timestamp",iter.Timestamp, "X",iter.X,"Y",iter.Y)
		if err := conn.Send("HMSET",key,"Mac",int64(iter.Mac), "storeid",storeid,"timestamp",iter.Timestamp, "X",iter.X,"Y",iter.Y);err !=nil{
			log.Error("conn.Send(\"HMset\", \"%s\") error(%v)",zkey,err)
			return err
		}
		
		if err := conn.Send("ZADD", zkey, iter.Timestamp, key); err != nil {
                        //key storeid:mac:ts
                        log.Error("conn.Send(\"ZADD\", \"%s\", %d, \"%s\") error(%v)", zkey, iter.Timestamp, key, err)
                        fmt.Println(err)
                        return err
                }
		
		//conn.Do("ZADD",zkey, iter.Timestamp, key)
	}
	if err := conn.Flush(); err != nil {
                log.Error("conn.Flush() error(%v)", err)
                return err
        }
	return nil
}

func (s *RedisStorage) getConnByNode(node string) redis.Conn {
	p, ok := s.pool[node]
	if !ok {
		log.Warn("no node: \"%s\" in redis pool", node)
		return nil
	}

	return p.Get()
}

// getConn get the connection of matching with key using ketama hashing.
func (s *RedisStorage) getConn(key string) redis.Conn {
	if len(s.pool) == 0 {
		return nil
	}
	node := s.ring.Hash(key)
	log.Debug("user_key: \"%s\" hit redis node: \"%s\"", key, node)
	return s.getConnByNode(node)
}

func (s *RedisStorage)Clear(){
	zoom.Close()
}

func InitRedis(){
	RedisInst = NewRedisStorage()
	conf:= &zoom.Configuration{
                Address:"localhost:6379",
                Network:"tcp",
        }
        zoom.Init(conf)
        zoom.Register(&Rssiample{})
        zoom.Register(&knn2.Finger{})
        zoom.Register(&knn2.Finger2{})
        zoom.Register(&knn2.ProcessData{})	
}
