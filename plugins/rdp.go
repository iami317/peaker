package plugins

import (
	"fmt"
	"github.com/iami317/peaker/util/grdp"
	"github.com/iami317/peaker/util/grdp/glog"
)

func ScanRdp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	g := grdp.NewClient(fmt.Sprintf("%s:%d", s.Ip, s.Port), glog.NONE)
	var err error
	//RDP协议登录测试
	err = g.LoginForSSL("", s.Username, s.Password)
	if err == nil {
		result.Result = true
		result.Class = WeakPass
		return result
	} else {
		return result
	}
}

func UnauthorizedRdp(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	g := grdp.NewClient(fmt.Sprintf("%s:%d", s.Ip, s.Port), glog.NONE)
	var err error
	//SSL协议登录测试
	err = g.LoginForSSL("", "", "")
	if err == nil {
		result.Result = true
		result.Class = Unauthorized
		return result
	}
	return result
}
