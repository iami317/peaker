package plugins

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
)

func UnauthorizedOracle(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	return result
}

func ScanOracle(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	connStr := fmt.Sprintf(
		"oracle://%s:%s@%s:%s/%s?connection_timeout=%d&connection_pool_timeout=%d",
		s.Username,
		s.Password,
		s.Ip,
		s.Port,
		"orcl",
		s.TimeOut,
		s.TimeOut)

	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		defer conn.Close()
		result.Result = true
	}
	return result
}
