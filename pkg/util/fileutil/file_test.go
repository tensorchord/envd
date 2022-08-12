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

package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindFileAbsPath(t *testing.T) {
	type testCase struct {
		name        string
		filePath    string
		fileName    string
		expectPath  string
		expectError bool
	}

	dir, err := os.MkdirTemp("", "envd-test")
	require.Nil(t, err, "create temp dir failed")
	defer os.RemoveAll(dir)
	name := "tmpfile"
	err = os.WriteFile(filepath.Join(dir, name), []byte("test"), 0666)
	require.Nil(t, err, "write temp file failed")
	expect, err := filepath.Abs(filepath.Join(dir, name))
	require.Nil(t, err, "cannot get the abs path")

	testCases := []testCase{
		{
			"path + file with path",
			dir,
			name,
			expect,
			false,
		},
		{
			"empth path + full file",
			"",
			expect,
			expect,
			false,
		},
		{
			"empty file name",
			dir,
			"",
			"",
			true,
		},
		{
			"path + full file",
			dir,
			expect,
			expect,
			false,
		},
	}

	for _, tc := range testCases {
		res, err := FindFileAbsPath(tc.filePath, tc.fileName)
		if tc.expectError {
			require.Error(t, err)
		} else {
			require.Equal(t, res, tc.expectPath, tc.name)
		}
	}
}
