package plugins

import (
	_ "github.com/gotomicro/dmgo"

	"database/sql"
	"fmt"
)

func ScanDm(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	db, err := sql.Open("dm", fmt.Sprintf("dm://%v:%v@%v:%v", s.Username, s.Password, s.Ip, s.Port))

	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}
	return result
}

func UnauthorizedDm(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	result.Single = s
	db, err := sql.Open("dm", fmt.Sprintf("dm://%v:%v@%v:%v", "root", "", s.Ip, s.Port))
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}
	return result
}
