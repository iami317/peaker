package plugins

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
)

func UnauthorizedSocks5(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	proxyURL, _ := url.Parse(fmt.Sprintf("socks5://%s:%s", s.Ip, s.Port))
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return result
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	req, err := http.NewRequest("GET", "https://baidu.com", nil)
	_, err = client.Do(req)
	if err != nil {
		return result
	}

	result.Result = true
	return result
}

func ScanSocks5(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	proxyURL, err := url.Parse(fmt.Sprintf("socks5://%s:%s@%s:%s", s.Username, s.Password, s.Ip, s.Port))
	if err != nil {
		return result
	}
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return result
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	req, err := http.NewRequest("GET", "https://baidu.com", nil)
	_, err = client.Do(req)
	if err != nil {
		return result
	}
	result.Result = true
	return result
}
