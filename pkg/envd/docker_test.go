package envd

import (
	"github.com/docker/docker/api/types/container"
	"reflect"
	"testing"
)

func TestDeviceRequests(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    []container.DeviceRequest
		wantErr bool
	}{
		{
			name: "device=1",
			args: args{
				value: "device=1",
			},
			want: []container.DeviceRequest{
				{
					Count:     0,
					DeviceIDs: []string{"1"},
					Driver:    "nvidia",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "device=1,3",
			args: args{
				value: "\"device=1,3\"",
			},
			want: []container.DeviceRequest{
				{
					Count:     0,
					DeviceIDs: []string{"1", "3"},
					Driver:    "nvidia",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "all",
			args: args{
				value: "all",
			},
			want: []container.DeviceRequest{
				{
					Count:  -1,
					Driver: "nvidia",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				value: "3",
			},
			want: []container.DeviceRequest{
				{
					Count:  3,
					Driver: "nvidia",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "count=5",
			args: args{
				value: "count=5",
			},
			want: []container.DeviceRequest{
				{
					Count:  5,
					Driver: "nvidia",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "device=1,3 with driver",
			args: args{
				value: "\"device=1,3\",\"driver=custom\"",
			},
			want: []container.DeviceRequest{
				{
					Count:     0,
					DeviceIDs: []string{"1", "3"},
					Driver:    "custom",
					Capabilities: [][]string{
						{"nvidia", "compute", "compat32", "graphics", "utility", "video", "display", "gpu"}},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
		{
			name: "device=1,3 with capabilities",
			args: args{
				value: "\"device=1,3\",\"capabilities=nvidia\"",
			},
			want: []container.DeviceRequest{
				{
					Count:     0,
					DeviceIDs: []string{"1", "3"},
					Driver:    "nvidia",
					Capabilities: [][]string{
						{"nvidia", "gpu"},
					},
					Options: make(map[string]string),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deviceRequests(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("deviceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deviceRequests() got = %v, want %v", got, tt.want)
			}
		})
	}
}
