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

package data

// hashString computes the hash of s.
func hashString(s string) uint32 {
	// if len(s) >= 12 {
	// Call the Go runtime's optimized hash implementation,
	// which uses the AESENC instruction on amd64 machines.
	// 	return uint32(goStringHash(s, 0))
	// }
	return softHashString(s)
}

////go:linkname goStringHash runtime.stringHash
// func goStringHash(s string, seed uintptr) uintptr

// softHashString computes the 32-bit FNV-1a hash of s in software.
func softHashString(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}
