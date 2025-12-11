package plugins

import (
	"fmt"
	"github.com/iami317/peaker/pkg/grdp"
	"strings"
)

func SplitUserDomain(user string) (string, string) {
	var domain string
	if strings.Contains(user, "/") {
		user = strings.Split(user, "/")[1]
		domain = strings.Split(user, "/")[0]
	}
	return user, domain
}

func ScanRdp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	user, domain := SplitUserDomain(s.Username)
	err := grdp.Login(fmt.Sprintf("%v:%v", s.Ip, s.Port), domain, user, s.Password)
	if err == nil {
		result.Result = true
	}
	return result
}

func UnauthorizedRdp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}

	return result
}
