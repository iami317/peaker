package plugins

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func ScanLdap(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	// LDAP连接配置
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err == nil {
		defer l.Close()
		// LDAP绑定（登录）
		err = l.Bind(s.Username, s.Password)
		if err == nil {
			result.Result = true
		}
	}

	return result
}

func UnauthorizedLdap(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	// LDAP连接配置
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err == nil {
		defer l.Close()
		// LDAP绑定（登录）
		err = l.Bind(s.Username, s.Password)
		if err == nil {
			result.Result = true
		}
	}
	return result
}
