package plugins

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"strings"
)

func ScanKibana(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetBasicAuth(s.Username, s.Password)
	res2, err2 := req.Get(fmt.Sprintf("http://%v:%v/app/kibana#/", s.Ip, s.Port))
	if err2 != nil {
		return result
	}
	body2, err2 := res2.Body()
	if err2 != nil {
		return result
	} else if res2.StatusCode() == 200 &&
		strings.Contains(string(body2), "kibana") &&
		strings.Contains(string(body2), "Management") {
		result.Result = true
	}

	return result
}

func UnauthorizedKibana(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)

	// 访问url1
	res, err := req.Get(fmt.Sprintf("http://%v:%v/app/kibana#", s.Ip, s.Port))
	if err != nil {
		return result
	}
	body, err := res.Body()
	if err != nil {
		return result
	} else if res.StatusCode() == 200 &&
		strings.Contains(string(body), "kibana") &&
		strings.Contains(string(body), "Management") {
		result.Result = true
		return result
	}
	return result

}
