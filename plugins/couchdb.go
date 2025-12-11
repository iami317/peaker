package plugins

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"strings"
	"time"
)

func UnauthorizedCouchdb(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(time.Duration(3))
	res, err := req.Get(fmt.Sprintf("http://%v:%v/_session/", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		defer res.Close()
		body, err := res.Body()
		if err == nil && strings.Contains(string(body), `"ok":true,"userCtx":{"name":"couchdb","roles":["`) {
			result.Result = true
		}
	}

	return result
}

func ScanCouchdb(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	req := HttpRequest.NewRequest()
	req.SetTimeout(time.Duration(3))
	req.SetBasicAuth(s.Username, s.Password)
	res, err := req.Get(fmt.Sprintf("http://%v:%v/_session/", s.Ip, s.Port))
	if err == nil && res != nil && res.StatusCode() == 200 {
		defer res.Close()
		body, err := res.Body()
		if err == nil && strings.Contains(string(body), `"ok":true,"userCtx":{"name":"couchdb","roles":["`) {
			var class Class
			if len(s.Password) > 0 {
				class = WeakPass
			} else {
				class = Unauthorized
			}
			return ScanResult{
				Single: s,
				Class:  class,
				Result: true,
			}
		}
	}

	return result
}
