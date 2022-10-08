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

	"github.com/moby/buildkit/client"
	gatewayclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestParseImportCache(t *testing.T) {
	type testCase struct {
		importCaches []string // --import-cache
		expected     []gatewayclient.CacheOptionsEntry
		expectedErr  string
	}
	testCases := []testCase{
		{
			importCaches: []string{"type=registry,ref=example.com/foo/bar", "type=local,src=/path/to/store"},
			expected: []gatewayclient.CacheOptionsEntry{
				{
					Type: "registry",
					Attrs: map[string]string{
						"ref": "example.com/foo/bar",
					},
				},
				{
					Type: "local",
					Attrs: map[string]string{
						"src": "/path/to/store",
					},
				},
			},
		},
		{
			importCaches: []string{"example.com/foo/bar", "example.com/baz/qux"},
			expected: []gatewayclient.CacheOptionsEntry{
				{
					Type: "registry",
					Attrs: map[string]string{
						"ref": "example.com/foo/bar",
					},
				},
				{
					Type: "registry",
					Attrs: map[string]string{
						"ref": "example.com/baz/qux",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		im, err := ParseImportCache(tc.importCaches)
		if tc.expectedErr == "" {
			require.EqualValues(t, tc.expected, im)
		} else {
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		}
	}
}

func TestParseExportCache(t *testing.T) {
	type testCase struct {
		exportCaches          []string // --export-cache
		legacyExportCacheOpts []string // --export-cache-opt (legacy)
		expected              []client.CacheOptionsEntry
		expectedErr           string
	}
	testCases := []testCase{
		{
			exportCaches: []string{"type=registry,ref=example.com/foo/bar"},
			expected: []client.CacheOptionsEntry{
				{
					Type: "registry",
					Attrs: map[string]string{
						"ref":  "example.com/foo/bar",
						"mode": "min",
					},
				},
			},
		},
		{
			exportCaches:          []string{"example.com/foo/bar"},
			legacyExportCacheOpts: []string{"mode=max"},
			expected: []client.CacheOptionsEntry{
				{
					Type: "registry",
					Attrs: map[string]string{
						"ref":  "example.com/foo/bar",
						"mode": "max",
					},
				},
			},
		},
		{
			exportCaches:          []string{"type=registry,ref=example.com/foo/bar"},
			legacyExportCacheOpts: []string{"mode=max"},
			expectedErr:           "--export-cache-opt is not supported for the specified --export-cache",
		},
		// TODO: test multiple exportCaches (valid for CLI but not supported by solver)

	}
	for _, tc := range testCases {
		ex, err := ParseExportCache(tc.exportCaches, tc.legacyExportCacheOpts)
		if tc.expectedErr == "" {
			require.EqualValues(t, tc.expected, ex)
		} else {
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		}
	}
}

func TestParseOutput(t *testing.T) {
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

func TestParseFromStr(t *testing.T) {
	type testCase struct {
		name        string
		from        string
		expectFile  string
		expectFunc  string
		expectError bool
	}

	testCases := []testCase{
		{
			"empty",
			"",
			"build.envd",
			"build",
			false,
		},
		{
			"without func",
			"main.envd",
			"main.envd",
			"build",
			false,
		},
		{
			"without func but has :",
			"test.envd:",
			"test.envd",
			"build",
			false,
		},
		{
			"without file",
			":test",
			"build.envd",
			"test",
			false,
		},
		{
			"all",
			"hello.envd:run",
			"hello.envd",
			"run",
			false,
		},
		{
			"more than 2 :",
			"test.envd:run:foo",
			"",
			"",
			true,
		},
	}

	for _, tc := range testCases {
		file, function, err := ParseFromStr(tc.from)
		if tc.expectError {
			require.Error(t, err)
		} else {
			require.Equal(t, file, tc.expectFile, tc.name)
			require.Equal(t, function, tc.expectFunc, tc.name)
		}
	}
}
