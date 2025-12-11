package plugins

import (
	"fmt"
	"github.com/streadway/amqp"
)

func UnauthorizedNeutron(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", "guest", "guest", s.Ip, s.Port))
	if err != nil {
		return result
	}
	defer conn.Close()
	result.Result = true
	return result
}

func ScanNeutron(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", s.Username, s.Password, s.Ip, s.Port))
	if err != nil {
		return result
	}
	defer conn.Close()
	result.Result = true
	return result
}
