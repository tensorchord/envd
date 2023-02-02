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
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

type Syncthing struct {
	Name          string
	Cmd           *exec.Cmd
	Config        *config.Configuration
	PrevConfig    config.Configuration // Unapplied config
	HomeDirectory string
	Port          string
	DeviceID      protocol.DeviceID
	Client        *Client
	ApiKey        string
	latestEventId int64
	DeviceAddress string
}

// Initializes the remote syncthing instance
func InitializeRemoteSyncthing() (*Syncthing, error) {
	s := &Syncthing{
		Name:          "Remote Syncthing",
		Port:          DefaultRemotePort,
		HomeDirectory: "/config",
		ApiKey:        DefaultApiKey,
	}

	s.Client = s.NewClient()

	err := s.WaitForStartup(15 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for syncthing startup: %w", err)
	}

	err = s.PullLatestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to pull latest config: %w", err)
	}

	logrus.Debug("Remote syncthing connected")

	err = s.SetDeviceAddress(DefaultRemoteDeviceAddress)
	if err != nil {
		return nil, err
	}
	s.Config.Options.RawListenAddresses = []string{DefaultRemoteDeviceAddress}
	s.DeviceID = s.Config.Devices[0].DeviceID

	err = s.ApplyConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Initializes the local syncthing instance
func InitializeLocalSyncthing(name string) (*Syncthing, error) {

	initConfig := InitLocalConfig()
	homeDirectory, err := GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	s := &Syncthing{
		Name:          "Local Syncthing",
		Config:        initConfig,
		HomeDirectory: homeDirectory,
		ApiKey:        DefaultApiKey,
		Port:          ParsePortFromAddress(initConfig.GUI.Address()),
	}

	s.Client = s.NewClient()

	logrus.Debug("Port for local syncthing is: ", initConfig.GUI.Address())

	if err != nil {
		return nil, err
	}

	if err = s.WriteLocalConfig(); err != nil {
		return nil, err
	}

	err = s.StartLocalSyncthing()
	if err != nil {
		return nil, err
	}

	err = s.PullLatestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to pull latest config: %w", err)
	}

	err = s.SetDeviceAddress(DefaultLocalDeviceAddress)
	if err != nil {
		return nil, err
	}

	s.Config.Options.RawListenAddresses = []string{DefaultLocalDeviceAddress}
	s.DeviceID = s.Config.Devices[0].DeviceID

	err = s.ApplyConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Syncthing) StartLocalSyncthing() error {
	if !IsInstalled() {
		err := InstallSyncthing()
		if err != nil {
			return fmt.Errorf("failed to install syncthing: %w", err)
		}
	}

	logrus.Debug("Starting local syncthing...")
	cmd := exec.Command(GetSyncthingBinPath(), "-no-restart", "-no-browser", "-home", s.HomeDirectory)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to run syncthing executable: %w", err)
	}
	s.Cmd = cmd
	logrus.Debug("Local syncthing started!")

	err = s.WaitForStartup(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for syncthing startup: %w", err)
	}

	// Handle the SIGINT signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		<-signalChan

		err := cmd.Process.Signal(os.Interrupt)
		if err != nil {
			logrus.Errorf("Failed to send SIGINT to syncthing: %s", err)
		}

		err = cmd.Wait()
		if err != nil {
			logrus.Errorf("Failed to wait for syncthing to exit: %s", err)
		}

		os.Exit(0)
	}()

	return nil
}

func (s *Syncthing) Ping() (bool, error) {
	_, err := s.Client.SendRequest(GET, "/rest/system/ping", nil, nil)
	if err != nil {
		logrus.Debug("Failed to ping syncthing: ", err)
		return false, fmt.Errorf("failed to ping syncthing: %w", err)
	}

	return true, nil
}

func (s *Syncthing) WaitForStartup(timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			logrus.Debugf("Timeout reached for syncthing: %s", s.Name)
			return fmt.Errorf("timed out waiting for syncthing to start")
		}
		if ok, _ := s.Ping(); ok {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Syncthing) StopLocalSyncthing() {
	if s.Cmd == nil {
		logrus.Error("syncthing is not running")
	}

	err := s.Cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		logrus.Errorf("failed to kill syncthing process: %s", err)
	}

	_, err = s.Cmd.Process.Wait()
	if err != nil {
		logrus.Errorf("failed to kill syncthing process: %s", err)
	}

	if err = CleanLocalConfig(s.Name); err != nil {
		logrus.Errorf("failed to clean local syncthing config: %s", err)
	}

}

func (s *Syncthing) IsRunning() bool {
	if s.Cmd == nil {
		return false
	}
	return s.Cmd.Process.Signal(syscall.Signal(0)) == nil
}
