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
