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

package v1

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLanguage(t *testing.T) {
	tcs := []struct {
		l                string
		ExpectedLanguage string
		ExpectedVersion  string
		ExpectedError    bool
	}{
		{
			l:                "python",
			ExpectedLanguage: "python",
			ExpectedVersion:  "",
			ExpectedError:    false,
		},
		{
			l:                "python3.7",
			ExpectedLanguage: "python",
			ExpectedVersion:  "3.7",
			ExpectedError:    false,
		},
		{
			l:                "python3.7.1",
			ExpectedLanguage: "python",
			ExpectedVersion:  "3.7.1",
			ExpectedError:    false,
		},
		{
			l:             "python-3.7.1",
			ExpectedError: true,
		},
		{
			l:                "r",
			ExpectedLanguage: "r",
			ExpectedVersion:  "",
			ExpectedError:    false,
		},
	}

	for _, tc := range tcs {
		language, version, err := parseLanguage(tc.l)
		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("parseLanguage(%s) returned error: %v", tc.l, err)
			}
		} else {
			if language != tc.ExpectedLanguage {
				t.Errorf("parseLanguage(%s) returned language %s, expected %s", tc.l, language, tc.ExpectedLanguage)
			}
			if version == nil {
				if tc.ExpectedVersion != "" {
					t.Errorf("parseLanguage(%s) returned version nil, expected %s", tc.l, tc.ExpectedVersion)
				}
			} else {
				if *version != tc.ExpectedVersion {
					t.Errorf("parseLanguage(%s) returned version %s, expected %s", tc.l, *version, tc.ExpectedVersion)
				}
			}
		}

	}
}

func TestIsRequirementsFileSafeToCopyContentPathEscape(t *testing.T) {
	envPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(envPath, "requirements.txt"), []byte("numpy\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// A file outside the build context that must never be read.
	outside := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(outside, []byte("secret-content\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// A symlink inside the build context pointing outside of it.
	if err := os.Symlink(outside, filepath.Join(envPath, "link.txt")); err != nil {
		t.Fatal(err)
	}

	relOutside, err := filepath.Rel(envPath, outside)
	if err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		name           string
		requirements   string
		ExpectedSafe   bool
		ExpectedLeaked bool
	}{
		{name: "normal file", requirements: "requirements.txt", ExpectedSafe: true},
		{name: "dot-dot escape", requirements: relOutside, ExpectedSafe: false},
		{name: "absolute path", requirements: outside, ExpectedSafe: false},
		{name: "symlink escape", requirements: "link.txt", ExpectedSafe: false},
	}

	for _, tc := range tcs {
		g := generalGraph{
			EnvironmentPath:  envPath,
			RequirementsFile: &tc.requirements,
		}
		dependencies, safe := g.IsRequirementsFileSafeToCopyContent()
		if safe != tc.ExpectedSafe {
			t.Errorf("%s: safe = %v, expected %v", tc.name, safe, tc.ExpectedSafe)
		}
		for _, dep := range dependencies {
			if dep == "secret-content" {
				t.Errorf("%s: content outside the build context was read", tc.name)
			}
		}
	}
}
