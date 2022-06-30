// Copyright 2022 The envd Authors
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
package builder

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_parseOutput(t *testing.T) {
	type args struct {
		output string
	}

	tests := []struct {
		name            string
		args            args
		outputType      string
		expectedEntries int
		wantErr         bool
	}{{
		"parsing output successfully",
		args{
			output: "type=tar,dest=test.tar",
		},
		"tar",
		1,
		false,
	}, {
		"output without type",
		args{
			output: "type=,dest=test.tar",
		},
		"",
		0,
		true,
	}, {
		"output without dest",
		args{
			output: "type=tar,dest=",
		},
		"tar",
		1,
		false,
	}, {
		"no output",
		args{
			output: "",
		},
		"",
		0,
		false,
	}}

	logrus.SetLevel(logrus.DebugLevel)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := parseOutput(tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(entries) != tt.expectedEntries {
				t.Errorf("parseOutput() expectedEntries = %v, got %v", tt.expectedEntries, len(entries))
				return
			}
			if len(entries) == 0 {
				return
			}
			if entries[0].Type != tt.outputType {
				t.Errorf("parseOutput() outputType = %s, want %s", entries[0].Type, tt.outputType)
			}
		})
	}
}
