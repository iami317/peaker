package plugins

import (
	"github.com/gosnmp/gosnmp"
	"time"
)

func UnauthorizedSnmp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	gosnmp.Default.Target = s.Ip
	gosnmp.Default.Port = uint16(s.Port)
	gosnmp.Default.Community = result.Single.Password
	gosnmp.Default.Timeout = s.TimeOut

	err := gosnmp.Default.Connect()
	if err == nil {
		_, err = gosnmp.Default.Get([]string{".1.3.6.1.2.1.1.1.0"})
		if err == nil {
			result.Class = Unauthorized
			result.Result = true
			return result
		}
	}
	return result
}

/*
oids列表参照：https://blog.csdn.net/qq_29752857/article/details/120223993
*/
func ScanSNMP(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	params := &gosnmp.GoSNMP{
		Target:        s.Ip,
		Port:          uint16(s.Port),
		Version:       gosnmp.Version3,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      gosnmp.AuthNoPriv,
		Timeout:       time.Duration(2) * time.Second,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 s.Username,
			AuthenticationProtocol:   gosnmp.MD5,
			AuthenticationPassphrase: s.Password,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        s.Password,
		},
	}
	err := params.Connect()
	if err == nil {
		defer params.Conn.Close()
		_, err = params.Get([]string{".1.3.6.1.2.1.1.1.0"})
		if err == nil {
			result.Class = WeakPass
			result.Result = true
			result.Single.Username = s.Username
			result.Single.Password = s.Password
		}
	}

	return result
}
