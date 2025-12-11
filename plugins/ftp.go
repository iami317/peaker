package plugins

import (
	"fmt"
	"github.com/jlaffaye/ftp"
)

func UnauthorizedFtp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	input1 := Single{Ip: s.Ip, Port: s.Port, Protocol: s.Protocol, Username: "anonymous", Password: "", TimeOut: s.TimeOut}
	r1, c1 := execScanFtp(input1)
	if r1 {
		return ScanResult{
			Single: input1,
			Result: true,
			Class:  c1,
		}
	}
	input2 := Single{Ip: s.Ip, Port: s.Port, Protocol: s.Protocol, Username: "FTP", Password: "", TimeOut: s.TimeOut}
	r2, c2 := execScanFtp(input2)
	if r2 {
		return ScanResult{
			Single: input2,
			Result: true,
			Class:  c2,
		}
	}
	return result
}

func ScanFtp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
	}

	r3, c3 := execScanFtp(s)
	if r3 {
		return ScanResult{
			Single: s,
			Result: true,
			Class:  c3,
		}
	}
	return result
}

func execScanFtp(s Single) (result bool, c Class) {
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", s.Ip, s.Port), s.TimeOut)
	if err == nil {
		err = conn.Login(s.Username, s.Password)
		if err == nil {
			defer conn.Logout()
			if len(s.Password) > 0 {
				return true, WeakPass
			} else {
				return true, Unauthorized
			}
		}
	}
	return false, UnKnow
}
