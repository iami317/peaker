package plugins

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/iami317/logx"
	"net"
)

func UnauthorizedHttp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), s.TimeOut)
	if err == nil {
		defer conn.Close()
		d := &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     "",
				Password: "",
			},
		}

		smb, err := d.Dial(conn)
		if err == nil {
			defer smb.Logoff()
			result.Class = Unauthorized
			result.Result = true
		}
	}
	return result
}

func ScanHttp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port))
	if err == nil {
		defer conn.Close()
		d := &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     s.Username,
				Password: s.Password,
			},
		}

		smb, err := d.Dial(conn)
		if err == nil {
			defer smb.Logoff()
			names, err := smb.ListSharenames()
			if err == nil {
				for _, name := range names {
					logx.Verbosef("show list:%v", name)
				}
			}

			result.Class = WeakPass
			result.Result = true
		}
	}
	return result
}
