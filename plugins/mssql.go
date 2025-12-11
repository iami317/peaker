package plugins

import (
	_ "github.com/denisenkom/go-mssqldb"

	"database/sql"
	"fmt"
)

func ScanMssql(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}

	db, err := sql.Open("mssql", fmt.Sprintf(
		"server=%v;port=%v;user id=%v;password=%v;database=%v",
		s.Ip,
		s.Port,
		s.Username,
		s.Password,
		"master"))
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return result
}

func UnauthorizedMssql(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	result.Single = s
	db, err := sql.Open("mssql", fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v", s.Ip, s.Port, "", "", "master"))
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return result
}
