package plugins

import (
	"github.com/gosnmp/gosnmp"
	"github.com/spf13/cast"
)

func UnauthorizedSnmp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	conn := &gosnmp.GoSNMP{
		Target:             s.Ip,
		Port:               cast.ToUint16(s.Port),
		Community:          "",
		Version:            gosnmp.Version2c,
		Timeout:            s.TimeOut,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err := conn.Connect()
	if err == nil {
		result.Result = true
	}
	return result
}

/*
oids列表参照：https://blog.csdn.net/qq_29752857/article/details/120223993
*/
func ScanSNMP(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	conn := &gosnmp.GoSNMP{
		Target:             s.Ip,
		Port:               cast.ToUint16(s.Port),
		Community:          s.Password,
		Version:            gosnmp.Version2c,
		Timeout:            s.TimeOut,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err := conn.Connect()
	if err == nil {
		result.Result = true
	}
	return result
}
