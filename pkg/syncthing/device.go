package syncthing

func ConnectDevices(s1 *Syncthing, s2 *Syncthing) error {
    var err error
	connectedDevices := append(s1.Config.Devices, s2.Config.Devices...)

	s1.Config.Devices = connectedDevices
	s2.Config.Devices = connectedDevices

	err = s1.ApplyConfig()
	if err != nil {
		return err
	}

	err = s2.ApplyConfig()
	if err != nil {
		return err
	}

return nil
}

func (s *Syncthing) SetDeviceAddress(addr string) (err error) {
	s.Config.Devices[0].Addresses = []string{addr}
	return nil
}
