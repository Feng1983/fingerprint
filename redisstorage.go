package main

import (
	log	"github.com/alecthomas/log4go" 
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
	"github.com/albrow/zoom"
	"github.com/Terry-Mao/gopush-cluster/ketama"
)

var (
	RedisNoConnErr       = errors.New("can't get a redis conn")
	redisProtocolSpliter = ":"
)

// RedisMessage struct encoding the composite info.



type RedisStorage struct {
	pool  map[string] *redis.Pool
	ring  *ketama.HashRing
	//delCH chan *RedisDelMessage
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
		pw := strings.Split(addr, redisProtocolSpliter)
		if len(pw) != 2 {
			log.Error("strings.Split(\"%s\", \"%s\") failed (%v)", addr, redisProtocolSpliter, err)
			panic(fmt.Sprintf("config redis.source node:\"%s\" format error", addr))
		}
		tmpProto := pw[0]
		tmpAddr := pw[1]
		// WARN: closures use
		redisPool[nw[0]] = &redis.Pool{
			MaxIdle:     Conf.RedisMaxIdle,
			MaxActive:   Conf.RedisMaxActive,
			IdleTimeout: Conf.RedisIdleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial(tmpProto, tmpAddr)
				if err != nil {
					log.Error("redis.Dial(\"%s\", \"%s\") error(%v)", tmpProto, tmpAddr, err)
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
	for _, iter:=range data{
		if err = conn.Send("ZADD", iter.Key, iter.Ts, iter.Value); err != nil {
			log.Error("conn.Send(\"ZADD\", \"%s\", %d, \"%s\") error(%v)", iter.Key, iter.Ts, iter.Value, err)
                	return err
        	}
	}
	if err = conn.Flush(); err != nil {
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
		zoom.Save(v)
	}
	return nil
}

func (s *RedisStorage) GetRedisRssi(storeid int, startTs,endTs int64,ids []int)[]*Rssiample{
	
	//conn:= getConnByNode(node)
	
	var sp []*Rssiample

	query:= zoom.NewQuery("Rssiample").Order("Ts").Filter("Ts >= ",startTs).Filter("Ts < ",endTs)
	for _,item:=range ids{
		query.Filter("Imac = ",item) 
	}
        
        //if err:= zoom.NewQuery("Rssiample").Order("Ts").Filter("").Filter().Scan(&sp);err!=nil{
	if err:= query.Scan(&sp);err!=nil{
            	panic(err)
        }
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
	
	for _, iter:=range dat{
		key:= string(storeid)+":"+iter.Mac+":"+iter.Timestamp
		
		conn.Do("HMSET",key,"Mac",int64(iter.Mac), "storeid",storeid,"timestamp",iter.Timestamp, "X",iter.X,"Y",iter.Y)
		conn.Do("ZADD",iter.Timestamp, key)
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


