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

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func ConnectDevices(s1 *Syncthing, s2 *Syncthing) error {
	logrus.Debug(fmt.Sprintf("Connecting syncthings %s and %s", s1.Name, s2.Name))
	var err error
	connectedDevices := append(s1.Config.Devices, s2.Config.Devices...)

	s1.Config.Devices = connectedDevices
	s2.Config.Devices = connectedDevices

	logrus.Debug("Adding device config for: ", s1.Name)
	err = s1.ApplyConfig()
	if err != nil {
		return err
	}

	logrus.Debug("Adding device config for: ", s2.Name)
	err = s2.ApplyConfig()
	if err != nil {
		return err
	}

	return nil
}

// This method can only be called when the devices are not connected
func (s *Syncthing) SetDeviceAddress(addr string) (err error) {
	if s.Config.Devices == nil || len(s.Config.Devices) == 0 {
		return fmt.Errorf("no devices found")
	}
	s.DeviceAddress = addr

	// Could work on better identifying which is the current device
	s.Config.Devices[0].Addresses = []string{addr}
	return nil
}

func (s *Syncthing) GetDeviceAddressPort() string {
	splitLst := strings.Split(s.DeviceAddress, ":")
	return splitLst[len(splitLst)-1]
}
