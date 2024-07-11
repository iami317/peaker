package plugins

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ScanMemcache(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	return result
}

func UnauthorizedMemcache(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), s.TimeOut)
	if err != nil {
		return result
	}
	defer conn.Close()
	_, err = conn.Write([]byte("stats\r\n"))
	if err != nil {
		return result
	}

	// 从服务器端收字符串
	reader := bufio.NewReader(conn)
	var buf [1024]byte
	n, err := reader.Read(buf[:])
	if err != nil {
		return result
	}
	resv := string(buf[:n])
	if strings.Contains(resv, "STAT version") {
		result.Class = Unauthorized
		result.Result = true
	}
	return result
}
