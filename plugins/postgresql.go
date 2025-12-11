package plugins

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ScanPostgres(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	dataSourceName := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		s.Username,
		s.Password,
		s.Ip,
		s.Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}
	return result
}

func UnauthorizedPostgres(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	result.Single = s
	dataSourceName := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		"",
		"",
		s.Ip,
		s.Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}
	return result
}
