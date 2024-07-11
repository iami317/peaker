package peaker

import (
	"fmt"
	"github.com/iami317/peaker/plugins"
)

type IpAddr struct {
	Ip       string
	Port     uint
	Protocol string
}

func (i IpAddr) String() string {
	return fmt.Sprintf("ip:%s port:%d protocol:%v", i.Ip, i.Port, i.Protocol)
}

type Crack struct {
	User   string
	Pass   string
	Class  uint
	Result bool
}

func (c Crack) String() string {
	return fmt.Sprintf("user:%v pass:%v class:%v result:%v", c.User, c.Pass, plugins.ClassMap[plugins.Class(c.Class)], c.Result)
}

type ResultOut struct {
	Addr  IpAddr
	Crack []Crack
}
