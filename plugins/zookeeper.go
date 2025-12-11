package plugins

import (
	"fmt"
	"github.com/go-zookeeper/zk"
)

func UnauthorizedZookeeper(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", s.Ip, s.Port)}, s.TimeOut)
	if err != nil {
		return result
	}
	defer conn.Close()
	result.Result = true
	return result
}

func ScanZookeeper(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: true,
	}
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", s.Ip, s.Port)}, s.TimeOut)
	if err != nil {
		return result
	}
	defer conn.Close()
	err = conn.AddAuth("digest", []byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	if err != nil {
		return result
	}
	result.Result = true

	return result
}
