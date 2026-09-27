// Copyright 2023 The envd Authors
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

package install

import (
	"os"
	"path/filepath"
	"testing"

	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v1/builtin"
)

func TestResolvePathInBuildContext(t *testing.T) {
	contextDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contextDir, "requirements.txt"), []byte("numpy\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(contextDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	// A file outside the build context, and a symlink inside the context
	// pointing to it.
	outside := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(outside, []byte("secret\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(contextDir, "link.txt")); err != nil {
		t.Fatal(err)
	}
	// A symlinked directory inside the context pointing outside of it.
	outsideDir := t.TempDir()
	if err := os.Symlink(outsideDir, filepath.Join(contextDir, "linkdir")); err != nil {
		t.Fatal(err)
	}

	thread := &starlark.Thread{Name: "test"}
	thread.SetLocal(builtin.BuildContextDir, contextDir)

	absInside := filepath.Join(contextDir, "requirements.txt")
	// Absolute, uncleaned paths with a non-existent target, so EvalSymlinks
	// fails and the lexical fallback path is exercised.
	absUncleanInside := filepath.Join(contextDir, "subdir", "..", "new", "requirements.txt")
	absUncleanOutside := filepath.Join(contextDir, "subdir", "..", "..", "etc", "passwd")

	tcs := []struct {
		name          string
		path          string
		expected      string
		ExpectedError bool
	}{
		{name: "relative file", path: "requirements.txt", expected: "requirements.txt"},
		{name: "dot-dot stays inside", path: "subdir/../requirements.txt", expected: "requirements.txt"},
		{name: "absolute path inside context", path: absInside, expected: "requirements.txt"},
		{name: "non-existent file inside context", path: "new/requirements.txt", expected: "new/requirements.txt"},
		{name: "uncleaned absolute path inside context", path: absUncleanInside, expected: "new/requirements.txt"},
		{name: "absolute path outside context", path: outside, ExpectedError: true},
		{name: "dot-dot escape", path: "../../etc/passwd", ExpectedError: true},
		{name: "uncleaned absolute path escaping context", path: absUncleanOutside, ExpectedError: true},
		{name: "symlink escape", path: "link.txt", ExpectedError: true},
		{name: "symlinked dir with non-existent file", path: "linkdir/new.txt", ExpectedError: true},
	}

	for _, tc := range tcs {
		resolved, err := resolvePathInBuildContext(thread, tc.path)
		if tc.ExpectedError {
			if err == nil {
				t.Errorf("%s: resolvePathInBuildContext(%s) expected error, got %s", tc.name, tc.path, resolved)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: resolvePathInBuildContext(%s) returned error: %v", tc.name, tc.path, err)
			continue
		}
		if resolved != tc.expected {
			t.Errorf("%s: resolvePathInBuildContext(%s) = %s, expected %s", tc.name, tc.path, resolved, tc.expected)
		}
	}
}

func TestResolvePathInBuildContextWithoutContext(t *testing.T) {
	thread := &starlark.Thread{Name: "test"}

	if _, err := resolvePathInBuildContext(thread, "requirements.txt"); err != nil {
		t.Errorf("resolvePathInBuildContext returned error for a relative path: %v", err)
	}
	if _, err := resolvePathInBuildContext(thread, "/etc/passwd"); err == nil {
		t.Error("resolvePathInBuildContext expected error for an absolute path without build context")
	}
	if _, err := resolvePathInBuildContext(thread, "../requirements.txt"); err == nil {
		t.Error("resolvePathInBuildContext expected error for a dot-dot escape without build context")
	}
}
