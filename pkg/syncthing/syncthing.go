package syncthing

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type Syncthing struct {
	cmd           *exec.Cmd
	Config        *config.Configuration
	HomeDirectory string
	Port          string
}

func Main() {
	// Configure local syncthing

	// Configure remote syncthing

	// Configure folders

	// Configure Device Connections
}

func InitializeRemoteSyncthing() (*Syncthing, error) {
	return nil, nil
}

// Writes the default configuration to the home directory
func InitializeLocalSyncthing() (*Syncthing, error) {
	initConfig := InitConfig()
	homeDirectory := DefaultHomeDirectory()
	s := &Syncthing{
		Config:        initConfig,
		HomeDirectory: homeDirectory,
	}

	port, err := parsePortFromAddress(initConfig.GUI.Address())
	if err != nil {
		return nil, fmt.Errorf("failed to parse port from address: %w", err)
	}
	s.Port = port

	configBytes, err := GetConfigByte(s.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to get syncthing config bytes: %w", err)
	}

	err = fileutil.CreateDirIfNotExist(homeDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to get syncthing config file path: %w", err)
	}

	configFilePath := GetConfigFilePath(homeDirectory)
	if err = fileutil.CreateIfNotExist(configFilePath); err != nil {
		return nil, fmt.Errorf("failed to get syncthing config file path: %w", err)
	}

	if err = os.WriteFile(configFilePath, configBytes, 0666); err != nil {
		return nil, fmt.Errorf("failed to write syncthing config file: %w", err)
	}

	return s, nil
}

func (s *Syncthing) Start() error {
	if !IsInstalled() {
		InstallSyncthing()
	}

	logrus.Debug("Starting syncthing...")
	cmd := exec.Command(GetSyncthingBinPath(), "-no-browser", "-no-restart", "-home", s.HomeDirectory)

	// TODO: Configure custom home path or default?
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to run syncthing executable: %w", err)
	}
	s.cmd = cmd
	logrus.Info("Syncthing started!")

	return nil
}

func (s *Syncthing) Restart() error {
	// TODO: use api endpoint to restart
	return nil
}

func (s *Syncthing) IsRunning() bool {
	if s.cmd == nil {
		return false
	}
	return s.cmd.Process.Signal(syscall.Signal(0)) == nil

}

func (s *Syncthing) Stop() error {
	if s.cmd == nil {
		return fmt.Errorf("syncthing is not running")
	}

	err := s.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill syncthing process: %w", err)
	}

	return nil
}
