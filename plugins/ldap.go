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
	var conn *ldap.Conn
	var err error
	// LDAP连接配置
	ldap.DefaultTimeout = s.TimeOut
	conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err == nil {
		defer conn.Close()
		// LDAP绑定（登录）
		err = conn.Bind(s.Username, s.Password)
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
	ldap.DefaultTimeout = s.TimeOut
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err == nil {
		defer conn.Close()
		// LDAP绑定
		err = conn.Bind(s.Username, s.Password)
		if err == nil {
			result.Result = true
		}
	}
	return result
}
