package plugins

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"strings"
)

func ScanTomcat(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetBasicAuth(s.Username, s.Password)
	req.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Upgrade-Insecure-Requests": "1",
	})
	res, err := req.Get(fmt.Sprintf("http://%v:%v/manager/html", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		defer func() {
			err := res.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()
		body, err := res.Body()
		if err == nil && strings.Contains(string(body), "Tomcat Web Application Manager") {
			return ScanResult{
				Single: s,
				Result: true,
				Class:  WeakPass,
			}
		}
	}
	return result
}

func UnauthorizedTomcat(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetHeaders(map[string]string{
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Upgrade-Insecure-Requests": "1",
	})
	res, err := req.Get(fmt.Sprintf("http://%v:%v/manager/html", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		body, err := res.Body()
		if err == nil && strings.Contains(string(body), "Tomcat Web Application Manager") {
			return ScanResult{
				Single: s,
				Result: true,
				Class:  Unauthorized,
			}
		}
	}
	if res != nil {
		res.Close()
	}
	return result
}
