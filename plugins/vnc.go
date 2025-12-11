package plugins

import (
	"fmt"
	"github.com/mitchellh/go-vnc"
	"net"
)

func UnauthorizedVnc(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}

	target := fmt.Sprintf("%v:%v", s.Ip, s.Port)
	tcpConn, err := net.DialTimeout("tcp", target, s.TimeOut)
	if err != nil {
		return result
	}
	defer tcpConn.Close()
	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: ""},
		},
	}
	vncConn, err := vnc.Client(tcpConn, &config)
	if err != nil {
		return result
	}
	defer vncConn.Close()
	result.Result = true

	return result
}

func ScanVnc(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: true,
	}
	target := fmt.Sprintf("%v:%v", s.Ip, s.Port)

	tcpConn, err := net.DialTimeout("tcp", target, s.TimeOut)
	if err != nil {
		return result
	}
	defer tcpConn.Close()
	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: s.Password},
		},
	}
	vncConn, err := vnc.Client(tcpConn, &config)
	if err != nil {
		return err
	}
	defer vncConn.Close()
	result.Result = true
	return result
}
