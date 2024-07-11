package plugins

import (
	gotelnet "github.com/iami317/peaker/util/go-telnet"
	"strconv"
	"time"
)

const (
	// 经过测试，linux下，延时需要大于100ms
	TimeDelayAfterWrite = 200 // 200ms
)

func UnauthorizedTelnet(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	client := gotelnet.TelNet{
		IP:               s.Ip,
		Port:             strconv.Itoa(int(s.Port)),
		IsAuthentication: true,
		UserName:         "",
		Password:         "",
		TimeOut:          50 * time.Second,
	}
	bl := client.Login()
	if bl {
		result.Result = true
		result.Class = Unauthorized
	}
	return result
}

func ScanTelnet(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	client := gotelnet.TelNet{
		IP:               s.Ip,
		Port:             strconv.Itoa(int(s.Port)),
		IsAuthentication: true,
		UserName:         s.Username,
		Password:         s.Password,
		TimeOut:          50 * time.Second,
	}
	bl := client.Login()
	if bl {
		result.Result = true
		result.Class = WeakPass
	}
	return result
}
