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

package netutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	assert.NoError(t, err)
	assert.NotEqual(t, port, 0)
}

func TestGetHost(t *testing.T) {
	tcs := []struct {
		host     string
		expected string
		err      bool
	}{
		{
			host:     "https://localhost:8080",
			expected: "localhost",
			err:      false,
		},
		{
			host:     "localhost:8080",
			expected: "",
			err:      true,
		},
		{
			host:     "http://localhost:8080",
			expected: "localhost",
			err:      false,
		},
		{
			host:     "http://1.1.1.1:8080",
			expected: "1.1.1.1",
			err:      false,
		},
	}
	for _, tc := range tcs {
		host, err := GetHost(tc.host)
		if tc.err == true {
			if err == nil {
				t.Error("expect to get the error, but got nil")
			}
			continue
		}
		if tc.err == false {
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if host != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, host)
			}
		}
	}
}
