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

package metrics

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/tensorchord/envd/pkg/driver"
)

type Collector interface {
	Watch(ctx context.Context, cid string) chan Metrics
	Stop() error
}

func GetCollector(name string, handle interface{}) (Collector, error) {
	ErrUnknownCollector := errors.Newf("unknown collector: %s", name)
	ErrUnknownHandle := errors.Newf("unknown handler: %s", name)
	switch name {
	case "docker":
		client, ok := handle.(driver.Client)
		if ok {
			return NewDockerCollector(client), nil
		}
	default:
		return nil, ErrUnknownCollector
	}
	return nil, ErrUnknownHandle
}
