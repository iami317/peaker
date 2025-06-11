package peaker

import (
	"bufio"
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/iami317/hubur"
	"github.com/iami317/logx"
	"github.com/iami317/peaker/plugins"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Weak struct {
	Config           Config
	StartTime        time.Time       //开始时间
	SupportProtocols map[string]bool //执行的协议
	ResultC          chan plugins.ScanResult
	Wg               *sync.WaitGroup
}

type Config struct {
	TimeOut      time.Duration //单个目标ip执行超时时间
	Ts           time.Duration //单个协议账号密码超时时间
	Thread       int           //并发目标的数量
	ThreadSingle int           //单个协议执行的并发数量
	Rate         int           //每秒钟发包速率
	Logger       *logx.Logger
	CheckAlive   bool //检测ip是否存活
	DebugMode    bool
	ResultFile   string
}

type RunIpData struct {
	Done       chan struct{}
	TimeOut    time.Duration
	Addr       IpAddr
	UserDict   []string
	PassDict   []string
	ResultChan chan interface{}
}

var (
	AliveAddr []IpAddr
	mutex     sync.Mutex
)

func init() {
	AliveAddr = make([]IpAddr, 0)
}

func NewWeak(c Config) *Weak {
	supportProtocols := make(map[string]bool)
	for proto := range plugins.ScanMap {
		p := string(proto)
		supportProtocols[strings.ToUpper(p)] = true
	}
	return &Weak{
		Config:           c,
		Wg:               &sync.WaitGroup{},
		SupportProtocols: supportProtocols,
	}
}

// RunTask 以ip为单位并发执行
func (w *Weak) RunTask(ipList []IpAddr, usersDict []string, passDict []string, resultChan chan interface{}) {
	defer func() {
		if e := recover(); e != nil {
			w.Config.Logger.Warnf(fmt.Sprintf("RunTask ERROR:%#v;stack=%s\n", e, string(debug.Stack())))
		}
		close(resultChan)
	}()
	w.Config.Logger.Verbosef("Start scanning %v targets", len(ipList))
	//是否检查ip是否存活
	if w.Config.CheckAlive {
		ipList = w.CheckAlive(ipList)
	}
	sema := hubur.NewSizedWaitGroup(w.Config.Thread)
	for _, addr := range ipList {
		if len(addr.Ip) > 0 && addr.Port > 0 && len(addr.Protocol) > 0 {
			sema.Add()
			go w.RunIpWithTimeout(addr, usersDict, passDict, resultChan, &sema)
		}
	}
	sema.Wait()
	w.Config.Logger.Verbosef("Finish,Latency %v seconds", time.Since(w.StartTime).Seconds())
	return
}

func (w *Weak) RunIpWithTimeout(addr IpAddr, usersDict []string, passDict []string, resultChan chan interface{}, sema *hubur.SizedWaitGroup) {
	defer func() {
		if e := recover(); e != nil {
			w.Config.Logger.Warnf(fmt.Sprintf("RunIpWithTimeout ERROR:%#v;stack=%s\n", e, string(debug.Stack())))
		}
		sema.Done()
	}()
	done := make(chan struct{}, 1)
	defer func() {
		if e := recover(); e != nil {
			w.Config.Logger.Warnf(fmt.Sprintf("RunIpWithTimeout111 ERROR:%#v;stack=%s\n", e, string(debug.Stack())))
		}
		close(done)
	}()
	param := RunIpData{
		Addr:       addr,
		UserDict:   usersDict,
		PassDict:   passDict,
		Done:       done,
		ResultChan: resultChan,
	}
	go w.RunIp(param)

	select {
	case <-time.After(w.Config.TimeOut):
		w.Config.Logger.Errorf("ip:%v-%v-%v,执行超时,单个ip限制 %v 秒", param.Addr.Ip, param.Addr.Port, param.Addr.Protocol, w.Config.TimeOut.Seconds())
		return
	case <-done:
		return
	}
}

func (w *Weak) RunIp(i interface{}) {
	defer func() {
		if e := recover(); e != nil {
			w.Config.Logger.Warnf(fmt.Sprintf("RunIp ERROR:%#v;stack=%s\n", e, string(debug.Stack())))
		}
	}()
	input := i.(RunIpData)
	protocol := plugins.Protocol(strings.ToUpper(input.Addr.Protocol))
	s, ok := plugins.ScanMap[protocol]
	if !ok {
		w.Config.Logger.Errorf("The current protocol %v is not supported", input.Addr.Protocol)
		input.Done <- struct{}{}
		return
	}

	if net.ParseIP(input.Addr.Ip) == nil || input.Addr.Ip == "0.0.0.0" {
		w.Config.Logger.Errorf("【%v】 not a valid IP address", input.Addr.Ip)
		input.Done <- struct{}{}
		return
	}

	if input.Addr.Port < 0 || input.Addr.Port > 65535 {
		w.Config.Logger.Errorf("【%v】 not a valid network port", input.Addr.Port)
		input.Done <- struct{}{}
		return
	}

	if !strings.Contains(input.Addr.Ip, "[") && hubur.IsIPv6(input.Addr.Ip) {
		input.Addr.Ip = fmt.Sprintf("[%v]", input.Addr.Ip)
	}

	rsOut := &ResultOut{
		Addr: input.Addr,
	}
	//获取需要并发执行的数量
	var thread int
	var timeout time.Duration
	if w.Config.ThreadSingle > 0 {
		thread = w.Config.ThreadSingle
	} else if s.Thread > 0 {
		thread = s.Thread
	} else {
		thread = plugins.DefaultThread
	}

	//获取协议的超时时间
	if w.Config.TimeOut.Seconds() > 0 {
		timeout = w.Config.Ts
	} else if s.Ts.Seconds() > 0 {
		timeout = s.Ts
	} else {
		timeout = plugins.DefaultTs
	}
	unauthorizedFuncPool := tunny.NewFunc(thread, s.UnauthorizedFunc)
	defer unauthorizedFuncPool.Close()
	//先执行未授权
	param := plugins.Single{
		TimeOut:  timeout,
		Ip:       input.Addr.Ip,
		Port:     input.Addr.Port,
		Protocol: string(input.Addr.Protocol),
		Username: "",
		Password: "",
	}
	if w.Config.DebugMode {
		w.Config.Logger.Verbosef("START=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v", param.Ip, param.Port, param.Protocol, "空", "空")
	}
	r, err := unauthorizedFuncPool.ProcessTimed(param, param.TimeOut)
	if err != nil && w.Config.DebugMode {
		w.Config.Logger.Errorf("TIMEOUT=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v ERR:%v", param.Ip, param.Port, param.Protocol, "空", "空", err.Error())
	}
	if r != nil {
		rs := r.(plugins.ScanResult)
		if rs.Result {
			rsOut.Crack = append(rsOut.Crack, Crack{
				User:   rs.Single.Username,
				Pass:   rs.Single.Password,
				Class:  uint(rs.Class),
				Result: rs.Result,
			})
			w.Config.Logger.Verbosef("SUCCESS=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v", rs.Single.Ip, rs.Single.Port, rs.Single.Protocol, "空", "空")
			input.ResultChan <- rsOut
			input.Done <- struct{}{}
			return
		}
	}

	if len(input.UserDict) > 0 && len(input.PassDict) > 0 {
		// 设置速率限制 - 每秒最多执行10次扫描
		rateLimit := 10
		if w.Config.Rate > 0 {
			rateLimit = w.Config.Rate
		}
		fmt.Println("rateLimit", rateLimit)
		ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
		defer ticker.Stop()

		scanFuncPool := tunny.NewFunc(thread, s.ScanFunc)
		defer scanFuncPool.Close()
		sema := hubur.NewSizedWaitGroup(thread)

		for _, user := range input.UserDict {
			if len(rsOut.Crack) > 0 && input.Addr.Protocol == rsOut.Addr.Protocol {
				break
			}
			for _, pass := range input.PassDict {
				// 等待速率限制器允许执行
				<-ticker.C
				sema.Add()
				paramScan := plugins.Single{
					TimeOut:  timeout,
					Ip:       input.Addr.Ip,
					Port:     input.Addr.Port,
					Protocol: string(input.Addr.Protocol),
					Username: user,
					Password: pass,
				}
				if len(rsOut.Crack) > 0 && input.Addr.Protocol == rsOut.Addr.Protocol {
					sema.Done()
					break
				}
				go func(param plugins.Single, rsOut *ResultOut, mode bool, l *logx.Logger) {
					defer func() {
						sema.Done()
					}()
					if w.Config.DebugMode {
						l.Verbosef("START=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v", param.Ip, param.Port, param.Protocol, param.Username, param.Password)
					}
					r, err := scanFuncPool.ProcessTimed(param, param.TimeOut)
					if err != nil && mode {
						l.Verbosef("TIMEOUT=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v", param.Ip, param.Port, param.Protocol, param.Username, param.Password)
						return
					}
					if r != nil {
						rs := r.(plugins.ScanResult)
						if rs.Result {
							rsOut.Crack = append(rsOut.Crack, Crack{
								User:   rs.Single.Username,
								Pass:   rs.Single.Password,
								Class:  uint(rs.Class),
								Result: rs.Result,
							})
							l.Verbosef("SUCCESS=> ip:%v,端口:%v,协议:%v,用户名:%v,密码:%v", rs.Single.Ip, rs.Single.Port, rs.Single.Protocol, rs.Single.Username, rs.Single.Password)
						}
					}
					return
				}(paramScan, rsOut, w.Config.DebugMode, w.Config.Logger)
			}
			if protocol == plugins.REDIS {
				break
			}
		}
		sema.Wait()
	}

	input.ResultChan <- rsOut
	input.Done <- struct{}{}
	return
}

// ReadUserDict 获取用户名字典文件
func (w *Weak) ReadUserDict(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		return nil, fmt.Errorf("Open user dict file err, %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users, err
}

// ReadPasswordDict 获取密码字典文件
func (w *Weak) ReadPasswordDict(passDict string) (password []string, err error) {
	file, err := os.Open(passDict)
	if err != nil {
		return nil, fmt.Errorf("Open password dict file err, %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			password = append(password, passwd)
		}
	}
	password = append(password, "")
	return password, err
}

// 获取目标
func (w *Weak) ReadIpList(fileName string) (ipList []IpAddr, err error) {
	ipListFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Open ip List file err, %v", err)
	}
	defer ipListFile.Close()
	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ipPort := strings.TrimSpace(line)
		t := strings.Split(ipPort, ":")
		if len(t) >= 2 {
			ip := t[0]
			portProtocol := t[1]
			tmpPort := strings.Split(portProtocol, "|")
			// ip列表中指定了端口对应的服务
			if len(tmpPort) == 2 {
				port, _ := strconv.Atoi(tmpPort[0])
				protocol := strings.ToUpper(tmpPort[1])
				if w.SupportProtocols[protocol] {
					addr := IpAddr{Ip: ip, Port: uint(port), Protocol: protocol}
					ipList = append(ipList, addr)
				} else {
					w.Config.Logger.Infof("Not support %v, ignore: %v:%v", protocol, ip, port)
				}
			}
		}

	}
	if len(ipList) == 0 {
		return nil, fmt.Errorf("No valid target was obtained")
	}
	return ipList, nil
}

func (w *Weak) CheckAlive(ipList []IpAddr) []IpAddr {
	w.Config.Logger.Debugf("checking target active...")
	wg := &sync.WaitGroup{}
	wg.Add(len(ipList))
	for _, addr := range ipList {
		go func(addr IpAddr) {
			defer wg.Done()
			var err error
			alive := false
			_, err = net.DialTimeout("udp", fmt.Sprintf("%v:%v", addr.Ip, addr.Port), time.Second*3)
			if err == nil {
				alive = true
			} else {
				_, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", addr.Ip, addr.Port), time.Second*3)
				if err == nil {
					alive = true
				}
			}
			if alive {
				mutex.Lock()
				AliveAddr = append(AliveAddr, addr)
				mutex.Unlock()
			}

		}(addr)
	}
	wg.Wait()
	w.Config.Logger.Debugf("Scan to %v targets alive", len(AliveAddr))
	return AliveAddr
}
