package plugins

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/olivere/elastic/v7"
	"net/http"
)

func UnauthorizedElastic(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("https://%v:%v", s.Ip, s.Port), fmt.Sprintf("http://%v:%v", s.Ip, s.Port)),
		elastic.SetMaxRetries(3),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
	)
	if err == nil {
		_, _, err = client.Ping(fmt.Sprintf("https://%v:%v", s.Ip, s.Port)).Do(context.Background())
		if err == nil {
			result.Class = WeakPass
			result.Result = true
		}
	}
	return result
}

func ScanElastic(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("https://%v:%v", s.Ip, s.Port), fmt.Sprintf("http://%v:%v", s.Ip, s.Port)),
		elastic.SetMaxRetries(3),
		elastic.SetBasicAuth(s.Username, s.Password),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
	)
	if err == nil {
		_, _, err = client.Ping(fmt.Sprintf("https://%v:%v", s.Ip, s.Port)).Do(context.Background())
		if err == nil {
			result.Class = WeakPass
			result.Result = true
		}
	}
	return result
}
