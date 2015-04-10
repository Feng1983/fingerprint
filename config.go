package main

import (
    //"encoding/json"
    "flag"
    "fmt"
    //"io/ioutil"
    "runtime"
    "time"
    "github.com/Terry-Mao/goconf"
)

var (
    Conf     *Config
    ConfFile string
)

func init() {
    flag.StringVar(&ConfFile, "c", "./trackpoint.conf", "set trackpoint config file path")
}


type Config struct {
    NodeWeight       int               `goconf:"base:node.weight"`
    User             string            `goconf:"base:user"`
    PidFile          string            `goconf:"base:pidfile"`
    Dir              string            `goconf:"base:dir"`
    Log              string            `goconf:"base:log"`
    MaxProc          int               `goconf:"base:maxproc"`
    PprofBind        []string          `goconf:"base:pprof.bind:,"`
    StorageType      string            `goconf:"storage:type"`
    RedisIdleTimeout time.Duration     `goconf:"redis:timeout:time"`
    RedisMaxIdle     int               `goconf:"redis:idle"`
    RedisMaxActive   int               `goconf:"redis:active"`
    RedisMaxStore    int               `goconf:"redis:store"`
    PostgreSQLClean       time.Duration     `goconf:"mysql:clean:time"`
    RedisSource      map[string]string `goconf:"-"`
    PostgreSQLSource      map[string]string `goconf:"-"`
}

// get a config
func InitConfig() error {
    gconf := goconf.New()
    if err := gconf.Parse(ConfFile); err != nil {
        return err
    }
    Conf = &Config{
        // base
        NodeWeight: 1,
        User:       "nobody nobody",
        PidFile:    "./trackpoint-message.pid",
        Dir:        "./",
        Log:        "./log.xml",
        MaxProc:    runtime.NumCPU(),
        PprofBind:  []string{"localhost:18170"},
        // storage
        StorageType: "redis",
        // redis
        RedisIdleTimeout: 28800 * time.Second,
        RedisMaxIdle:     50,
        RedisMaxActive:   1000,
        RedisMaxStore:    20,
        RedisSource:      make(map[string]string),
        // postgresql
        PostgreSQLSource: make(map[string]string),
        PostgreSQLClean:  1 * time.Hour,
    }
    if err := gconf.Unmarshal(Conf); err != nil {
        return err
    }
    // redis section
    redisAddrsSec := gconf.Get("redis.source")
    if redisAddrsSec != nil {
        for _, key := range redisAddrsSec.Keys() {
            addr, err := redisAddrsSec.String(key)
            if err != nil {
                return fmt.Errorf("config section: \"redis.addrs\" key: \"%s\" error(%v)", key, err)
            }
            Conf.RedisSource[key] = addr
        }
    } 
    // postgresql section
    dbSource := gconf.Get("mysql.source")
    if dbSource != nil {
        for _, key := range dbSource.Keys() {
            source, err := dbSource.String(key)
            if err != nil {
                return fmt.Errorf("config section: \"postgresql.source\" key: \"%s\" error(%v)", key, err)
            }
            Conf.PostgreSQLSource[key] = source
        }
    }
    return nil
}        
