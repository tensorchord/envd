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

package envutil

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetDurationWithDefault parses a duration string from environment variable with the given key
// or uses the given default duration. The duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func GetDurationWithDefault(key string, o time.Duration) time.Duration {
	v, found := os.LookupEnv(key)
	if found && v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
		log.WithField(key, v).WithError(err).Panic("failed to parse")
	}
	return o
}
