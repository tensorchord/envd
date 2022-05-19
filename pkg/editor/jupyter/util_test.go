package jupyter

import (
	"testing"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

func TestGenerateCommand(t *testing.T) {
	testcases := []struct {
		graph    ir.Graph
		dir      string
		expected []string
	}{
		{
			graph: ir.Graph{
				JupyterConfig: &ir.JupyterConfig{
					Password: "",
					Port:     8888,
				},
			},
			dir: "test",
			expected: []string{
				"python3", "-m", "notebook", "--ip", "0.0.0.0", "--notebook-dir", "test",
				"--NotebookApp.password", "''", "--NotebookApp.token", "''",
				"--port", "8888",
			},
		},
		{
			graph: ir.Graph{
				JupyterConfig: &ir.JupyterConfig{
					Password: "test",
					Port:     8888,
				},
			},
			dir: "test",
			expected: []string{
				"python3", "-m", "notebook", "--ip", "0.0.0.0", "--notebook-dir", "test",
				"--NotebookApp.password", "test", "--NotebookApp.token", "''",
				"--port", "8888",
			},
		},
		{
			graph:    ir.Graph{},
			dir:      "test",
			expected: []string{},
		},
	}
	for _, tc := range testcases {
		actual := GenerateCommand(tc.graph, tc.dir)
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
