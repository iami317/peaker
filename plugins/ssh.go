package plugins

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

func ScanSsh(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
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
			result.Class = WeakPass
			result.Result = true
		}
	}
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
