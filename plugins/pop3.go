package plugins

import (
	"github.com/knadh/go-pop3"
	"github.com/spf13/cast"
)

func UnauthorizedPop3(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	return result
}

func ScanPop3(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	p := pop3.New(pop3.Opt{
		Host:       s.Ip,
		Port:       cast.ToInt(s.Port),
		TLSEnabled: false,
	})
	conn, err := p.NewConn()
	if err != nil {
		return result
	}
	defer conn.Quit()
	if err := conn.Auth(s.Username, s.Password); err != nil {
		return result
	}
	result.Result = true
	return result
}
