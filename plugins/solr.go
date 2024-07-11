package plugins

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"strings"
)

func ScanSolr(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	req.SetBasicAuth("solr", "SolrRocks")
	res, err := req.Get(fmt.Sprintf("http://%v:%v/solr/admin/cores?wt=json&indexInfo=false", s.Ip, s.Port))
	if err == nil {
		body, err := res.Body()
		if err == nil &&
			strings.Contains(string(body), "instanceDir") &&
			strings.Contains(string(body), "dataDir") {
			result.Class = WeakPass
			result.Result = true
		}
	}
	return result
}

func UnauthorizedSolr(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	req := HttpRequest.NewRequest()
	req.SetTimeout(s.TimeOut)
	res, err := req.Get(fmt.Sprintf("http://%v:%v/solr/admin/cores?wt=json&indexInfo=false", s.Ip, s.Port))
	if err == nil {
		body, err := res.Body()
		if err == nil &&
			strings.Contains(string(body), "instanceDir") &&
			strings.Contains(string(body), "dataDir") {
			result.Class = Unauthorized
			result.Result = true
		}
	}
	return result
}
