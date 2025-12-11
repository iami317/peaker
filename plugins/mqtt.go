package plugins

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func UnauthorizedMqtt(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  Unauthorized,
		Result: false,
	}
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", s.Ip, s.Port))
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		result.Result = false
		return result
	}

	result.Result = true
	return result
}

func ScanMqtt(i interface{}) interface{} {
	s := i.(Single)
	result := ScanResult{
		Single: s,
		Class:  WeakPass,
		Result: false,
	}
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", s.Ip, s.Port)).
		SetUsername(s.Username).
		SetPassword(s.Password)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return result
	}
	result.Result = true
	return result
}
