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

type Metrics struct {
	NCpus        uint8
	CPUUtil      int
	NetTx        int64
	NetRx        int64
	MemLimit     int64
	MemPercent   int
	MemUsage     int64
	IOBytesRead  int64
	IOBytesWrite int64
	Pids         int
}

func NewMetrics() Metrics {
	return Metrics{
		CPUUtil:      -1,
		NetTx:        -1,
		NetRx:        -1,
		MemUsage:     -1,
		MemPercent:   -1,
		IOBytesRead:  -1,
		IOBytesWrite: -1,
		Pids:         -1,
	}
}
