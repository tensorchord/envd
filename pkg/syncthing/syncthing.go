package syncthing

import (
	"fmt"
	"net/http"
	"os/exec"
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
	Client        *http.Client
	DeviceID      protocol.DeviceID
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
		Client:        NewApiClient(),
		ApiKey:        DefaultApiKey,
	}

    err := s.WaitForStartup(15 * time.Second)
    if err != nil {
        return nil, fmt.Errorf("failed to wait for syncthing startup: %w", err)
    }

	err = s.PullLatestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to pull latest config: %w", err)
	}

	logrus.Debug("Remote syncthing connected")

	s.SetDeviceAddress(DefaultRemoteDeviceAddress)
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
    // Get free port for local envd 

	initConfig := InitLocalConfig()
    homeDirectory := GetHomeDirectory(name)
	s := &Syncthing{
		Name:          "Local Syncthing",
		Config:        initConfig,
		HomeDirectory: homeDirectory,
		Client:        NewApiClient(),
		ApiKey:        DefaultApiKey,
	}

	port := ParsePortFromAddress(initConfig.GUI.Address())
	s.Port = port

    var err error
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

	s.SetDeviceAddress(DefaultLocalDeviceAddress)
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
		InstallSyncthing()
	}

	logrus.Debug("Starting local syncthing...")
	// TODO: Add no browser option later
	cmd := exec.Command(GetSyncthingBinPath(), "-no-restart", "-home", s.HomeDirectory)

	// TODO: Configure custom home path or default?
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

	return nil
}

func (s *Syncthing) Ping() (bool, error) {
	_, err := s.ApiCall(GET, "/rest/system/ping", nil, nil)
	if err != nil {
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

func (s *Syncthing) StopLocalSyncthing() error {
	if s.Cmd == nil {
		return fmt.Errorf("syncthing is not running")
	}

	err := s.Cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		return fmt.Errorf("failed to kill syncthing process: %w", err)
	}

	_, err = s.Cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("failed to kill syncthing process: %w", err)
	}

	if err = s.CleanLocalConfig(); err != nil {
		return fmt.Errorf("failed to clean local syncthing config: %w", err)
	}

	return nil
}

func (s *Syncthing) Restart() error {
	// TODO: use api endpoint to restart
	return nil
}

func (s *Syncthing) IsRunning() bool {
	if s.Cmd == nil {
		return false
	}
	return s.Cmd.Process.Signal(syscall.Signal(0)) == nil
}
