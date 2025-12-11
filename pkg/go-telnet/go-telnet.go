package gotelnet

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"net"
)

const (
	// IAC interpret as command
	IAC = 255
	// SB is subnegotiation of the indicated option follows
	SB = 250
	// SE is end of subnegotiation parameters
	SE = 240
	// WILL indicates the desire to begin
	WILL = 251
	// WONT indicates the refusal to perform,
	// continue performing, the indicated option
	WONT = 252
	// DO indicates the request that the other
	// party perform, or confirmation that you are
	// expecting the other party to perform, the indicated option
	DO = 253
	// DONT indicates the demand that the other
	// party stop performing, or confirmation that you
	// are no longer expecting the other party to
	// perform, the indicated option
	DONT = 254
)

var loginRe *regexp.Regexp = regexp.MustCompile("[\\w\\d-_]+ login:")
var passwordRe *regexp.Regexp = regexp.MustCompile("Password:")
var bannerRe *regexp.Regexp = regexp.MustCompile(
	"[\\w\\d-_]+@[\\w\\d-_]+:[\\w\\d/-_~]+(\\$|#)")

// TelnetClient is basic descriptor
type TelnetClient struct {
	Login     string
	Password  string
	Address   string
	Port      string
	Timeout   time.Duration
	Verbose   bool
	LogWriter *bufio.Writer
	reader    *bufio.Reader
	writer    *bufio.Writer
	conn      net.Conn
}

func (tc *TelnetClient) setDefaultParams() {
	if tc.Port == "" {
		tc.Port = "23"
	}
	if tc.Timeout == 0 {
		tc.Timeout = 10 * time.Second
	}
	if tc.Verbose && tc.LogWriter == nil {
		tc.LogWriter = bufio.NewWriter(os.Stdout)
	}
}

func (tc *TelnetClient) log(format string, params ...interface{}) {
	if tc.Verbose {
		fmt.Fprintf(tc.LogWriter, "telnet: "+format+"\n", params...)
		tc.LogWriter.Flush()
	}
}

// Dial does open connect to telnet server
func (tc *TelnetClient) Dial() (err error) {
	tc.setDefaultParams()

	tc.log("Trying connect to %s:%s", tc.Address, tc.Port)
	tc.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", tc.Address, tc.Port), tc.Timeout)
	if err != nil {
		return
	}
	_ = tc.conn.SetReadDeadline(time.Now().Add(tc.Timeout))
	_ = tc.conn.SetWriteDeadline(time.Now().Add(tc.Timeout))

	tc.reader = bufio.NewReader(tc.conn)
	tc.writer = bufio.NewWriter(tc.conn)

	err = tc.conn.SetReadDeadline(time.Now().Add(tc.Timeout))
	if err != nil {
		fmt.Println("--------err123", err)
		return
	}

	tc.log("Waiting for the first banner")
	err = tc.waitWelcomeSigns()
	fmt.Println("--------err456", err)
	return
}

func (tc *TelnetClient) Close() {
	tc.conn.Close()
}

func (tc *TelnetClient) skipSBSequence() (err error) {
	var peeked []byte

	for {
		_, err = tc.reader.Discard(1)
		if err != nil {
			return
		}

		peeked, err = tc.reader.Peek(2)
		if err != nil {
			return
		}

		if peeked[0] == IAC && peeked[1] == SE {
			_, err = tc.reader.Discard(2)
			break
		}
	}

	return
}

func (tc *TelnetClient) skipCommand() (err error) {
	var peeked []byte

	peeked, err = tc.reader.Peek(1)
	if err != nil {
		return
	}

	switch peeked[0] {
	case WILL, WONT, DO, DONT:
		_, err = tc.reader.Discard(2)
	case SB:
		err = tc.skipSBSequence()
	}

	return
}

// ReadByte 从远程服务器接收字节，跳过命令
func (tc *TelnetClient) ReadByte() (b byte, err error) {
	for {
		b, err = tc.reader.ReadByte()
		if err != nil || b != IAC {
			fmt.Println("wwwwww", err)
			break
		}

		err = tc.skipCommand()
		if err != nil {
			fmt.Println("AAAAAA", err)
			break
		}
	}

	return
}

// ReadUntil 读取字节，直到特定符号。分隔符字符将写入结果缓冲区
func (tc *TelnetClient) ReadUntil(data *[]byte, delim byte) (n int, err error) {
	var b byte

	for b != delim {
		b, err = tc.ReadByte()
		if err != nil {
			break
		}

		*data = append(*data, b)
		n++
	}

	return
}

func findNewLinePos(data []byte) int {
	var pb byte

	for i := len(data) - 1; i >= 0; i-- {
		cb := data[i]
		if pb == '\n' && cb == '\r' {
			return i
		}

		pb = cb
	}

	return -1
}

// ReadUntilPrompt 读取数据，直到进程函数停止。
// 如果 process 函数返回 true，则读取将停止 如果找到下一行分隔符，则 Process 函数给出行块，即从行首到最后一个空格或整行
func (tc *TelnetClient) ReadUntilPrompt(process func(data []byte) bool) (output []byte, err error) {
	var n int
	var delimPos int
	var linePos int
	var chunk []byte

	output = make([]byte, 0, 64*1024)

	for {
		// Usually, if system print a prompt,
		// it requires inputing data and
		// prompt has ':' or whitespace in end of line.
		// However, may be cases which have another behaviors.
		// So client may freeze
		n, err = tc.ReadUntil(&output, ' ')
		if err != nil {
			fmt.Println("---------------****1", err)
			return
		}

		delimPos += n
		n = findNewLinePos(output)
		if n != -1 {
			linePos = n + 2
		}

		chunk = output[linePos:delimPos]
		fmt.Println("------------------chunk")
		fmt.Println(string(chunk))
		fmt.Println("------------------chunk")
		if process(chunk) {
			break
		}
	}

	return
}

// ReadUntilBanner 读取直到横幅，即命令的整个输出
func (tc *TelnetClient) ReadUntilBanner() (output []byte, err error) {
	output, err = tc.ReadUntilPrompt(func(data []byte) bool {
		m := bannerRe.Find(data)
		return len(m) > 0
	})

	output = bannerRe.ReplaceAll(output, []byte{})
	output = bytes.Trim(output, " ")

	return
}

func (tc *TelnetClient) findInputPrompt(re *regexp.Regexp, responce string, buffer []byte) bool {
	match := re.Find(buffer)
	if len(match) == 0 {
		return false
	}

	tc.Write([]byte(responce + "\r\n"))

	return true
}

// waitWelcomeSigns 等待第一个横幅出现 如果检测到登录提示，它将授权
func (tc *TelnetClient) waitWelcomeSigns() (err error) {
	_, err = tc.ReadUntilPrompt(func(data []byte) bool {
		if tc.findInputPrompt(loginRe, tc.Login, data) {
			tc.log("Found login prompt")
			return false
		}
		if tc.findInputPrompt(passwordRe, tc.Password, data) {
			tc.log("Found password prompt")
			return false
		}

		m := bannerRe.Find(data)
		return len(m) > 0
	})

	return
}

// 写入发送原始数据以刷新 telnet 服务器
func (tc *TelnetClient) Write(data []byte) (n int, err error) {
	n, err = tc.writer.Write(data)
	if err == nil {
		err = tc.writer.Flush()
	}

	return
}

// Execute 在远程服务器上执行 send 命令并返回整个输出
func (tc *TelnetClient) Execute(name string, args ...string) (stdout []byte, err error) {
	_, err = tc.reader.Discard(tc.reader.Buffered())
	if err != nil {
		return
	}

	request := []byte(name + " " + strings.Join(args, " ") + "\r\n")
	tc.log("Send command: %s", request[:len(request)-2])
	tc.Write(request)

	stdout, err = tc.ReadUntilBanner()
	if err != nil {
		return
	}
	tc.log("Received data with size = %d", len(stdout))

	return
}
