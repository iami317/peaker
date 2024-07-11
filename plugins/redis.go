package plugins

import (
	"fmt"
	"github.com/go-redis/redis"
)

func UnauthorizedRedis(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	input1 := Single{
		Ip:       s.Ip,
		Port:     s.Port,
		Protocol: s.Protocol,
		Password: "",
	}
	r1, c1 := execScanRedis(input1)

	if r1 {
		return ScanResult{
			Single: Single{Ip: s.Ip, Port: s.Port, Protocol: s.Protocol, Password: ""},
			Result: true,
			Class:  c1,
		}
	}
	return result
}

func ScanRedis(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	r2, c2 := execScanRedis(s)
	if r2 {
		result.Single.Username = ""
		result.Result = r2
		result.Class = c2
	}
	return result
}

func execScanRedis(s Single) (result bool, c Class) {
	opt := redis.Options{
		Addr:        fmt.Sprintf("%v:%v", s.Ip, s.Port),
		Password:    s.Password,
		DB:          0,
		DialTimeout: s.TimeOut,
	}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err := client.Ping().Result()
	if err == nil {
		if len(s.Password) > 0 {
			return true, WeakPass
		} else {
			return true, Unauthorized
		}
	}
	return false, UnKnow
}
