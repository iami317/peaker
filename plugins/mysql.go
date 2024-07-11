package plugins

import (
	_ "gorm.io/driver/mysql"

	"database/sql"
	"fmt"
)

func ScanMysql(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	dataSourceName := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8",
		s.Username,
		s.Password,
		s.Ip,
		s.Port, "mysql")
	db, err := sql.Open("mysql", dataSourceName)

	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Class = WeakPass
			result.Result = true
		}
	}
	return result
}

func UnauthorizedMysql(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{}
	result.Single = s
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", "root", "", s.Ip, s.Port, "mysql")
	db, err := sql.Open("mysql", dataSourceName)

	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Class = Unauthorized
			result.Result = true
		}
	}
	return result
}
