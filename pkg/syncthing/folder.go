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

package syncthing

import "github.com/syncthing/syncthing/lib/config"

func SyncFolder(s1 *Syncthing, s2 *Syncthing, dir1 string, dir2 string) error {
	baseFolder := config.FolderConfiguration{
		ID:               "default",
		RescanIntervalS:  5,
		FSWatcherEnabled: true,
		FSWatcherDelayS:  10,
		Devices: []config.FolderDeviceConfiguration{
			{
				DeviceID: s1.DeviceID,
			},
			{
				DeviceID: s2.DeviceID,
			},
		},
	}

	s1Folder := baseFolder.Copy()
	s2Folder := baseFolder.Copy()

	s1Folder.Path = dir1
	s2Folder.Path = dir2

	s1.Config.SetFolder(s1Folder)
	s2.Config.SetFolder(s2Folder)

	err := s1.ApplyConfig()
	if err != nil {
		return err
	}

	err = s2.ApplyConfig()
	if err != nil {
		return err
	}

	return nil
}
