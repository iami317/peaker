package plugins

import (
	"testing"
	"time"
)

func TestScanSsh(t *testing.T) {
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
			name: "ssh001",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     22,
					Protocol: "SSH",
					Username: "root",
					Password: "root",
					TimeOut:  3 * time.Second,
				},
			},
			want: ScanResult{
				Result: true,
			},
		},
		{
			name: "ssh002",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     22,
					Protocol: "SSH",
					Username: "ubuntu",
					Password: "u_16_docker",
					TimeOut:  3 * time.Second,
				},
			},
			want: ScanResult{
				Result: true,
			},
		},
		{
			name: "ssh003",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     22,
					Protocol: "SSH",
					Username: "root1",
					Password: "root",
					TimeOut:  3 * time.Second,
				},
			},
			want: ScanResult{
				Result: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScanSsh(tt.args.i); got.(ScanResult).Result != tt.want.(ScanResult).Result {
				t.Errorf("ScanSsh() = %v, want %v", got, tt.want)
			}
		})
	}
}
