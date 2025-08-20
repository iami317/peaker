package plugins

import (
	"fmt"
	"github.com/iami317/hubur"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

var client *ssh.Client

func ScanSsh(i interface{}) interface{} {
	var err error
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout:         s.TimeOut,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "SSH-2.0-OpenSSH_9.9",
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr"}, // 显式指定常用算法
		},
	}

	_ = hubur.Retry(
		func() error {
			//fmt.Println(fmt.Sprintf("**********%v:%v %v %v", s.Ip, s.Port, s.Username, s.Password))
			client, err = ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
			if err != nil {
				//fmt.Println(fmt.Sprintf("=====%v:%v %v %v %v %v %v", s.Ip, s.Port, s.Username, s.Password, s.TimeOut, result.Result, err.Error()))
				if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "handshake failed: EOF") {
					return err
				}
				return nil
			}
			defer client.Close()

			session, err := client.NewSession()

			if err != nil {
				//fmt.Println(fmt.Sprintf("-----%v:%v %v %v %v %v", s.Ip, s.Port, s.Username, s.Password, result.Result, err.Error()))
				if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "handshake failed: EOF") {
					return err
				}
				return nil
			}
			defer session.Close()

			result.Class = WeakPass
			result.Result = true
			return nil
		},
		hubur.RetryTimes(3),
		hubur.RetryDuration(time.Second*5),
	)

	return result
}

func UnauthorizedSsh(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	config := &ssh.ClientConfig{
		User:    s.Username,
		Timeout: s.TimeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			defer session.Close()
			result.Class = Unauthorized
			result.Result = true
		}
	}
	return result
}
