package plugins

import (
	"github.com/kirinlabs/HttpRequest"
	"strings"
)

func UnauthorizedHadoop(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	res, err := req.Get("http://%v:%v", s.Ip, s.Port)
	if err == nil {
		body, err := res.Body()
		defer res.Close()
		if err == nil &&
			strings.Contains(string(body), "Applications") &&
			strings.Contains(string(body), "hadoop") {
			result.Result = true
		}
	}
	return result
}

func ScanHadoop(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetBasicAuth(s.Username, s.Password)
	res, err := req.Get("http://%v:%v", s.Ip, s.Port)
	if err == nil {
		body, err := res.Body()
		defer res.Close()
		if err == nil && strings.Contains(string(body), "Applications") && strings.Contains(string(body), "hadoop") {
			result.Result = true
		}
	}
	return result
}
