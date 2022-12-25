package syncthing

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Syncthing struct {
	cmd *exec.Cmd
}

func (s *Syncthing) Start() error {
	if !IsInstalled() {
		InstallSyncthing()
	}

	logrus.Debug("Starting syncthing...")
	cmd := exec.Command(getSyncthingBinPath(), "-no-browser", "-no-restart")
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
