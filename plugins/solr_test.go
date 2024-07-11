package plugins

import (
	"testing"
	"time"
)

func TestScanSolr(t *testing.T) {
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
			name: "solr_001",
			args: args{
				i: Single{
					Ip:       "192.168.105.24",
					Port:     23,
					Protocol: "telnet",
					Username: "ubuntu",
					Password: "u_16_docker",
					TimeOut:  3 * time.Second,
				},
			},
			want: ScanResult{
				Result: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScanSolr(tt.args.i); got.(ScanResult).Result != tt.want.(ScanResult).Result {
				t.Errorf("ScanSolr() = %v, want %v", got, tt.want)
			}
		})
	}
}
