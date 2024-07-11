package plugins

import (
	"testing"
	"time"
)

func TestScanTelnet(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{
			name: "telnet_001",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     23,
					Protocol: "telnet",
					Username: "ubuntu",
					Password: "u_16_docker",
					TimeOut:  30 * time.Second,
				},
			},
			want: ScanResult{
				Result: true,
			},
		},
		{
			name: "telnet_002(err:user)",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     23,
					Protocol: "telnet",
					Username: "ubuntu1",
					Password: "u_16_docker",
					TimeOut:  30 * time.Second,
				},
			},
			want: ScanResult{
				Result: false,
			},
		},
		{
			name: "telnet_003(err:pass)",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     23,
					Protocol: "telnet",
					Username: "ubuntu",
					Password: "u_16_docker1",
					TimeOut:  30 * time.Second,
				},
			},
			want: ScanResult{
				Result: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScanTelnet(tt.args.i); got.(ScanResult).Result != tt.want.(ScanResult).Result {
				t.Errorf("ScanTelnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
