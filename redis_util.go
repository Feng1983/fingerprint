
package main
import(
	"time"
	"strings"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var defaultTimeout = 1 * time.Second

func GetRedisStat(addr string) (map[string]string, error) {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	ret, err := redis.String(c.Do("INFO"))
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	lines := strings.Split(ret, "\n")
	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			m[k] = v
		}
	}

	var reply []string

	reply, err = redis.Strings(c.Do("config", "get", "maxmemory"))
	if err != nil {
		return nil, err
	}
	// we got result
	if len(reply) == 2 {
		if reply[1] != "0" {
			m["maxmemory"] = reply[1]
		} else {
			m["maxmemory"] = "8"
		}
	}

	return m, nil
}
func GetRedisConfig(addr string, configName string) (string, error) {
	//c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	//if err != nil {
	//}
	return "", nil
}

func main(){
	p,_:=GetRedisStat(":6379")
	for k, v:=range p{
		fmt.Println(k,v)
	}
}
