package plugins

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/iami317/peaker/pkg"
	"github.com/iami317/peaker/pkg/encode"
	"github.com/spf13/cast"
	"golang.org/x/crypto/md4"
	"strconv"
	"strings"
)

func RsyncDetect(target string, timeout int) (float64, []string, error) {
	conn, err := pkg.NewSocket("tcp", target, timeout)
	if err != nil {
		return 0, nil, nil
	}
	defer conn.Close()

	rev, err := conn.Request([]byte("@RSYNCD: 31.0\n\n"), 1024)
	if err != nil {
		return 0, nil, nil
	}
	if bytes.Contains(rev, []byte("@RSYNCD: ")) && len(rev) > 13 {
		ver, err := strconv.ParseFloat(string(rev[9:13]), 32)
		if err != nil {
			return 0, nil, nil
		}
		if ss := strings.Split(string(rev), "\n"); len(ss) > 1 {
			modules := strings.Fields(ss[1])
			if len(modules) > 0 {
				return ver, modules, nil
			}
		}
		return ver, nil, nil
	}
	return 0, nil, errors.New("not rsync")
}

func RsyncUnauth(target string, ver float64, modules []string, timeout int) error {
	conn, err := pkg.NewSocket("tcp", target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Request([]byte(fmt.Sprintf("@RSYNCD: %f \n", ver)), 1024)
	if err != nil {
		return err
	}

	var data []byte
	if len(modules) > 0 {
		data = []byte(modules[0] + "\n")
	} else {
		data = []byte("\n")
	}
	rev, err := conn.Request(data, 1024)
	if err != nil {
		return err
	}
	if bytes.Contains(rev, []byte("@RSYNCD: OK")) {
		return nil
	}
	return errors.New("not unauth rsyncd")
}

func RsyncLogin(target, user, passwd string, ver float64, modules []string, timeout int) error {
	conn, err := pkg.NewSocket("tcp", target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Request([]byte(fmt.Sprintf("@RSYNCD: %f \n", ver)), 1024)
	if err != nil {
		return err
	}

	var data []byte
	if len(modules) > 0 {
		data = []byte(modules[0] + "\n")
	} else {
		data = []byte("\n")
	}
	rev2, err := conn.Request(data, 1024)
	if err != nil {
		return err
	}

	if !bytes.Contains(rev2, []byte("@RSYNCD: AUTHREQD")) {
		return errors.New("not found challenge")
	}
	if ss := strings.Fields(string(rev2)); len(ss) < 2 {
		return errors.New("not found challenge")
	} else {
		challenge := ss[2]
		c1 := passwd + challenge
		var hash []byte
		if ver >= 30 {
			md := md5.New()
			md.Write([]byte(c1))
			hash = md.Sum(nil)
		} else {
			md := md4.New()
			md.Write([]byte(c1))
			hash = md.Sum(nil)
		}

		c2 := encode.Base64Encode(hash)
		c2 = strings.Trim(c2, "=")
		rev3, err := conn.Request([]byte(user+" "+c2+"\n"), 1024)
		if err != nil {
			return err
		}
		if strings.Contains(string(rev3), "OK") {
			return nil
		}
	}

	return errors.New("rsync auth failed")
}

func UnauthorizedRsync(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	timeOut := cast.ToInt(s.TimeOut.Seconds())
	address := fmt.Sprintf("%v:%v", s.Ip, s.Port)
	ver, modules, err := RsyncDetect(address, timeOut)
	if err != nil {
		return result
	}
	err = RsyncUnauth(address, ver, modules, timeOut)
	if err != nil {
		return result
	}
	result.Result = true
	return result
}

func ScanRsync(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	timeOut := cast.ToInt(s.TimeOut.Seconds())
	address := fmt.Sprintf("%v:%v", s.Ip, s.Port)
	ver, modules, err := RsyncDetect(address, timeOut)
	if err != nil {
		return result
	}

	err = RsyncLogin(address, s.Username, s.Password, ver, modules, timeOut)
	if err != nil {
		return result
	}

	result.Result = true
	return result
}
