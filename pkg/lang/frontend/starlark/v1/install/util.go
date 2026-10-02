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
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v1/builtin"
)

// resolveExisting resolves symlinks for the deepest existing ancestor of
// candidate and appends the remaining (non-existent) components lexically,
// because filepath.EvalSymlinks fails when any component does not exist.
func resolveExisting(candidate string) string {
	if r, err := filepath.EvalSymlinks(candidate); err == nil {
		return r
	}
	parent := filepath.Dir(candidate)
	if parent == candidate {
		return filepath.Clean(candidate)
	}
	return filepath.Join(resolveExisting(parent), filepath.Base(candidate))
}

// resolvePathInBuildContext resolves path against the build context directory
// and returns it relative to the build context, so that the path can be reused
// inside the image where the context is mounted. It returns an error if the
// path escapes the build context (via an absolute path, `..` traversal, or a
// symlink). This validation does not pin filesystem objects; host reads must
// independently enforce containment through a directory-rooted open.
func resolvePathInBuildContext(thread *starlark.Thread, path string) (string, error) {
	contextDir, _ := thread.Local(builtin.BuildContextDir).(string)
	if contextDir == "" {
		// The build context is unknown (e.g. in unit tests); only accept
		// paths that are already relative and stay inside the context.
		if filepath.IsAbs(path) {
			return "", errors.Errorf("path %s must be relative to the build context, absolute paths are not allowed", path)
		}
		cleaned := filepath.Clean(path)
		if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
			return "", errors.Errorf("path %s escapes the build context", path)
		}
		return path, nil
	}

	contextAbs, err := filepath.Abs(contextDir)
	if err != nil {
		return "", errors.Wrapf(err, "failed to resolve the build context %s", contextDir)
	}
	contextAbs, err = filepath.EvalSymlinks(contextAbs)
	if err != nil {
		return "", errors.Wrapf(err, "failed to resolve the build context %s", contextDir)
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(contextAbs, candidate)
	}

	// Resolve symlinks so a symlink pointing outside the context is rejected,
	// even when the final components of the path do not exist yet.
	resolved := resolveExisting(candidate)
	if resolved != contextAbs && !strings.HasPrefix(resolved, contextAbs+string(filepath.Separator)) {
		return "", errors.Errorf("path %s is outside the build context %s", path, contextDir)
	}

	rel, err := filepath.Rel(contextAbs, resolved)
	if err != nil {
		return "", errors.Wrapf(err, "failed to relativize path %s", path)
	}
	return rel, nil
}
