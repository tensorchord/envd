// Copyright 2024 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package envd

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/container"
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
