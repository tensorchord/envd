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

package v1

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (g generalGraph) CacheID(filename string) string {
	var cacheID string
	if g.CUDA != nil {
		cacheID = fmt.Sprintf("%s/%s-gpu", filename, g.EnvironmentName)
	} else {
		cacheID = fmt.Sprintf("%s/%s-cpu", filename, g.EnvironmentName)
	}
	logrus.Debugf("apt/pypi calculated cacheID: %s", cacheID)
	return cacheID
}
