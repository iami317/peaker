package grdp

import (
	"fmt"
	"github.com/iami317/peaker/util/grdp/glog"
	"testing"
)

func testrdp(target string) {
	domain := ""
	username := "administrator"
	password := "930517"
	//target = "180.102.17.30:3389"
	var err error
	g := NewClient(target, glog.NONE)
	//SSL协议登录测试
	err = g.LoginForSSL(domain, username, password)
	if err == nil {
		fmt.Println("Login Success")
		return
	}
	if err.Error() != "PROTOCOL_RDP" {
		fmt.Println("Login Error:", err)
		return
	}
	//RDP协议登录测试
	err = g.LoginForRDP(domain, username, password)
	if err == nil {
		fmt.Println("Login Success")
		return
	} else {
		fmt.Println("Login Error:", err)
		return
	}
}

func TestName(t *testing.T) {
	targetArr := []string{
		//"50.57.49.172:3389",
		//"20.49.22.250:3389",
		"192.168.106.25:3389",
	}
	for _, target := range targetArr {
		fmt.Println(target)
		testrdp(target)
	}
}
