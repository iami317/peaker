package gotelnet

import (
	"fmt"
	"github.com/iami317/logx"
	"log"
	"net"
	"strings"
	"time"
)

const (
	DelayAfterWrite = 200 * time.Millisecond
)

type TelNet struct {
	IP               string
	Port             string
	IsAuthentication bool
	UserName         string
	Password         string
	TimeOut          time.Duration
	Verbose          bool
}

func (t *TelNet) Login() bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", t.IP, t.Port), t.TimeOut)
	if conn != nil {
		defer conn.Close()
	}
	if nil != err {
		return false
	}
	bl, err := t.telnetProtocolHandshake(conn)
	if !bl {
		return false
	}
	if err != nil && t.Verbose {
		logx.Verbose(err)
	}

	return true

}

func (t *TelNet) telnetProtocolHandshake(conn net.Conn) (bool, error) {
	var buf [1024]byte
	var err error
	var n int
	if n, err = conn.Read(buf[0:]); err != nil {
		if t.Verbose {
			logx.Verbose(err)
		}
		return false, err
	}

	buf[1] = 252
	buf[4] = 252
	buf[7] = 252
	buf[10] = 252

	if n, err = conn.Write(buf[0:n]); err != nil {
		if t.Verbose {
			logx.Verbose(err)
		}
		return false, err
	}

	time.Sleep(DelayAfterWrite)
	if n, err = conn.Read(buf[0:]); err != nil {
		if t.Verbose {
			logx.Verbose(err)
		}
		return false, err
	}

	buf[1] = 252
	buf[4] = 251
	buf[7] = 252
	buf[10] = 254
	buf[13] = 252
	n, err = conn.Write(buf[0:n])
	if nil != err {
		return false, err
	}
	time.Sleep(DelayAfterWrite)
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false, err
	}
	if strings.Contains(string(buf[0:n]), "login:") {
		n, err = conn.Write([]byte(t.UserName + "\n"))
		if nil != err {
			return false, err
		}
		time.Sleep(DelayAfterWrite)
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Println("pkg: model, func: telnetProtocolHandshake, method: conn.Read, errInfo:", err)
		return false, err
	}
	if false == t.IsAuthentication {
		return true, nil
	}
	if strings.Contains(string(buf[0:n]), "Password:") {
		n, err = conn.Write([]byte(t.Password + "\n"))
		if nil != err {
			return false, err
		}
		time.Sleep(DelayAfterWrite)
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false, err
	}

	if strings.Contains(string(buf[0:n]), "Last login:") {
		return true, nil
	}
	return false, nil
}
