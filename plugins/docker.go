package plugins

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"strings"
)

func UnauthorizedDocker(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	res, err := req.Get(fmt.Sprintf("http://%v:%v/version", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		defer res.Close()
		body, err := res.Body()
		if err == nil &&
			res.StatusCode() == 200 && strings.Contains(string(body), "Platform") {
			return ScanResult{
				Single: s,
				Class:  Unauthorized,
				Result: true,
			}
		}
	}
	return result
}

func ScanDocker(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetBasicAuth(s.Username, s.Password)
	res, err := req.Get(fmt.Sprintf("http://%v:%v/version", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		defer res.Close()
		body, err := res.Body()
		if err == nil && strings.Contains(string(body), "Platform") {
			return ScanResult{
				Single: s,
				Class:  WeakPass,
				Result: true,
			}
		}
	}

	return result
}
