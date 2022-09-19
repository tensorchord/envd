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

package ir

import (
	"testing"
)

func TestGenerateCommand(t *testing.T) {
	testcases := []struct {
		graph    Graph
		dir      string
		expected []string
	}{
		{
			graph: Graph{
				JupyterConfig: &JupyterConfig{
					Token: "",
					Port:  8888,
				},
				uid: 1000,
			},
			dir: "test",
			expected: []string{
				"python3", "-m", "notebook", "--ip", "0.0.0.0", "--notebook-dir", "test",
				"--NotebookApp.token", "''", "--port", "8888",
			},
		},
		{
			graph: Graph{
				JupyterConfig: &JupyterConfig{
					Token: "test",
					Port:  8888,
				},
				uid: 1000,
			},
			dir: "test",
			expected: []string{
				"python3", "-m", "notebook", "--ip", "0.0.0.0", "--notebook-dir", "test",
				"--NotebookApp.token", "test", "--port", "8888",
			},
		},
		{
			graph: Graph{
				JupyterConfig: &JupyterConfig{
					Token: "test",
					Port:  8888,
				},
				uid: 0,
			},
			dir: "test",
			expected: []string{
				"python3", "-m", "notebook", "--ip", "0.0.0.0", "--notebook-dir", "test",
				"--NotebookApp.token", "test", "--port", "8888", "--allow-root",
			},
		},
		{
			graph:    Graph{},
			dir:      "test",
			expected: []string{},
		},
	}
	for _, tc := range testcases {
		actual := tc.graph.generateJupyterCommand(tc.dir)
		if !equal(actual, tc.expected) {
			t.Errorf("failed to generate the command: expected %v, got %v", tc.expected, actual)
		}
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
