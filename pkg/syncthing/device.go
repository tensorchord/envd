package syncthing

import (
	"fmt"

	"github.com/syncthing/syncthing/lib/config"
)

func ConnectDevices(s1 *Syncthing, s2 *Syncthing) error {
	err := s1.PullLatestConfig()
	if err != nil {
		return err
	}

	err = s2.PullLatestConfig()
	if err != nil {
		return err
	}

	fmt.Println(s1.Config.Devices, s2.Config.Devices)

	devices := append(s1.Config.Devices, s2.Config.Devices...)
	connectedDevices := []config.DeviceConfiguration{}

	for _, device := range devices {
		device.Addresses = []string{DefaultDeviceAddress}
		connectedDevices = append(connectedDevices, device)
	}

	s1.Config.Devices = connectedDevices
	s2.Config.Devices = connectedDevices

	fmt.Println(connectedDevices)

	fmt.Println("Applying config for s1")
	err = s1.ApplyConfig()
	if err != nil {
		return err
	}

	fmt.Println("Applying config for s2")
	err = s2.ApplyConfig()
	if err != nil {
		return err
	}

	return nil
}
