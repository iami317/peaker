package plugins

import (
	"crypto/tls"
	"github.com/beltran/gohive"
)

/*
*
Ubuntu: sudo apt-get install libkrb5-dev
MacOS: brew install homebrew/dupes/heimdal --without-x11
Debian: yum install -y krb5-devel
*/
func ScanHive(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	configuration := gohive.NewConnectConfiguration()
	configuration.Username = s.Username
	configuration.Password = s.Password
	configuration.ConnectTimeout = s.TimeOut
	configuration.HttpTimeout = s.TimeOut
	configuration.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := gohive.Connect(s.Ip, int(s.Port), "NONE", configuration)
	if err == nil && conn != nil {
		defer conn.Close()
		result.Result = true
	}
	return result
}

func UnauthorizedHive(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}

	configuration := gohive.NewConnectConfiguration()
	configuration.ConnectTimeout = s.TimeOut
	configuration.HttpTimeout = s.TimeOut
	configuration.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := gohive.Connect(s.Ip, int(s.Port), "NONE", configuration)
	if err == nil && conn != nil {
		defer conn.Close()
		result.Result = true
	}
	return result
}
