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

package v0

import "testing"

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
