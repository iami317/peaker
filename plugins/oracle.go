package plugins

//
//import (
//	_ "gopkg.in/rana/ora.v4"
//
//	"database/sql"
//	"fmt"
//)
//
//func ScanOracle(i interface{}) interface{} {
//	s := i.(Single)
//	result := ScanResult{
//		Single: s,
//	}
//	result.Single = s
//	dataSourceName := fmt.Sprintf(
//		"%v:%v@tcp(%v:%v)/%v",
//		s.Username,
//		s.Password,
//		s.Ip,
//		s.Port, "orcl")
//
//	db, err := sql.Open("ora", dataSourceName)
//	defer db.Close()
//	if err == nil {
//		err = db.Ping()
//		if err == nil {
//			result.Class = WeakPass
//			result.Result = true
//		}
//	}
//
//	return result
//}
