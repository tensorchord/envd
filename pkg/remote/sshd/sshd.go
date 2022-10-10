// Copyright 2022 The envd Authors
// Copyright 2022 The okteto remote Authors
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

package sshd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/cockroachdb/errors"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
)

// LoadAuthorizedKeys loads path as an array.
// It will return nil if path doesn't exist.
func LoadAuthorizedKeys(path string) ([]ssh.PublicKey, error) {
	authorizedKeysBytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	authorizedKeys := []ssh.PublicKey{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, err
		}

		authorizedKeys = append(authorizedKeys, pubKey)
		authorizedKeysBytes = rest
	}

	if len(authorizedKeys) == 0 {
		return nil, errors.New("no keys found")
	}

	return authorizedKeys, nil
}

// Server holds the ssh server configuration.
type Server struct {
	Port  int
	Shell string

	AuthorizedKeys []ssh.PublicKey
	Hostkey        ssh.Signer
}

// ListenAndServe starts the SSH server using port
func (srv *Server) ListenAndServe() error {
	server, err := srv.getServer()
	if err != nil {
		return errors.Wrap(err, "failed to parse server configs")
	}
	return server.ListenAndServe()
}

//nolint:unparam
func (srv *Server) getServer() (*ssh.Server, error) {
	forwardHandler := &ssh.ForwardedTCPHandler{}

	server := &ssh.Server{
		Addr:    fmt.Sprintf(":%d", srv.Port),
		Handler: srv.connectionHandler,
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"direct-tcpip": ssh.DirectTCPIPHandler,
			"session":      ssh.DefaultSessionHandler,
		},
		LocalPortForwardingCallback: ssh.LocalPortForwardingCallback(func(ctx ssh.Context, dhost string, dport uint32) bool {
			log.Println("Accepted forward", dhost, dport)
			return true
		}),
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		RequestHandlers: map[string]ssh.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": sftpHandler,
		},
	}

	if srv.AuthorizedKeys != nil {
		server.PublicKeyHandler = srv.authorize
	} else {
		server.PublicKeyHandler = nil
		server.PasswordHandler = nil
	}

	if srv.Hostkey != nil {
		server.AddHostKey(srv.Hostkey)
	}

	return server, nil
}

func (srv Server) buildCmd(logger *logrus.Entry, s ssh.Session) *exec.Cmd {
	var cmd *exec.Cmd

	if len(s.RawCommand()) == 0 {
		cmd = exec.Command(srv.Shell)
	} else {
		args := []string{"-c", s.RawCommand()}
		cmd = exec.Command(srv.Shell, args...)
	}

	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, s.Environ()...)

	logger.Debugf("ssh server command: %s", cmd.String())
	return cmd
}

func (srv *Server) connectionHandler(s ssh.Session) {
	// Set session ID as the logger field.
	sessionID := uuid.New().String()
	l := logrus.New()
	l.SetLevel(logrus.GetLevel())
	logger := l.WithField("session.id", sessionID)

	defer func() {
		s.Close()
		logger.Info("session closed")
	}()

	logger.Infof("starting ssh session with command '%+v'", s.RawCommand())

	cmd := srv.buildCmd(logger, s)

	if ssh.AgentRequested(s) {
		logger.Info("agent requested")
		l, err := ssh.NewAgentListener()
		if err != nil {
			logger.WithError(err).Error("failed to start agent")
			return
		}

		defer l.Close()
		go ssh.ForwardAgentConnections(l, s)
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", "SSH_AUTH_SOCK", l.Addr().String()))
	}

	ptyReq, winCh, isPty := s.Pty()
	if isPty {
		logger.Infoln("handling PTY session")
		if err := handlePTY(logger, cmd, s, ptyReq, winCh); err != nil {
			sendErrAndExit(logger, s, err)
			return
		}

		err := s.Exit(0)
		if err != nil {
			logger.Warningln("exit session with error:", err)
		}
		return
	}

	logger.Infoln("handling non PTY session")
	if err := handleNoTTY(logger, cmd, s); err != nil {
		sendErrAndExit(logger, s, err)
		return
	}

	err := s.Exit(0)
	if err != nil {
		logger.Warningln("exit session with error:", err)
	}
}

func handlePTY(logger *logrus.Entry, cmd *exec.Cmd, s ssh.Session, ptyReq ssh.Pty, winCh <-chan ssh.Window) error {
	if len(ptyReq.Term) > 0 {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
	}

	f, err := pty.Start(cmd)
	if err != nil {
		logger.WithError(err).Error("failed to start pty session")
		return err
	}

	go func() {
		for win := range winCh {
			setWinSize(f, win.Width, win.Height)
		}
	}()

	go func() {
		_, err := io.Copy(f, s) // stdin
		if err != nil {
			logger.WithError(err).Warningln("failed to copy stdin")
		}
	}()

	waitCh := make(chan struct{})
	go func() {
		defer close(waitCh)
		_, err := io.Copy(s, f) // stdout
		if err != nil {
			logger.WithError(err).Warningln("failed to copy stdin")
		}
	}()

	if err := cmd.Wait(); err != nil {
		logger.WithError(err).Errorf("pty command failed while waiting")
		return err
	}

	select {
	case <-waitCh:
		logger.Info("stdout finished")
	case <-time.NewTicker(1 * time.Second).C:
		logger.Info("stdout didn't finish after 1s")
	}

	return nil
}

func setWinSize(f *os.File, w, h int) {
	// TODO(gaocegege): Should we use syscall or docker resize?
	// Refer to https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-docker/docker.go#L99
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
	if err != 0 {
		logrus.WithError(err).Error("failed to set winsize")
	}
}

func sendErrAndExit(logger *logrus.Entry, s ssh.Session, err error) {
	msg := strings.TrimPrefix(err.Error(), "exec: ")
	if _, err := s.Stderr().Write([]byte(msg)); err != nil {
		logger.WithError(err).Errorf("failed to write error back to session")
	}

	if err := s.Exit(getExitStatusFromError(err)); err != nil {
		logger.WithError(err).Errorf("pty session failed to exit")
	}
}

func getExitStatusFromError(err error) int {
	if err == nil {
		return 0
	}

	var exitErr exec.ExitError
	if ok := errors.As(err, &exitErr); !ok {
		return 1
	}

	waitStatus, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		if exitErr.Success() {
			return 0
		}

		return 1
	}

	return waitStatus.ExitStatus()
}

func handleNoTTY(logger *logrus.Entry, cmd *exec.Cmd, s ssh.Session) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.WithError(err).Errorf("couldn't get StdoutPipe")
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.WithError(err).Errorf("couldn't get StderrPipe")
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.WithError(err).Errorf("couldn't get StdinPipe")
		return err
	}

	if err = cmd.Start(); err != nil {
		logger.WithError(err).Errorf("couldn't start command '%s'", cmd.String())
		return err
	}

	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, s); err != nil {
			logger.WithError(err).Errorf("failed to write session to stdin.")
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(s, stdout); err != nil {
			logger.WithError(err).Errorf("failed to write stdout to session.")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(s.Stderr(), stderr); err != nil {
			logger.WithError(err).Errorf("failed to write stderr to session.")
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		logger.WithError(err).Errorf("command failed while waiting")
		return err
	}

	return nil
}

func (srv *Server) authorize(ctx ssh.Context, key ssh.PublicKey) bool {
	for _, k := range srv.AuthorizedKeys {
		if ssh.KeysEqual(key, k) {
			logrus.Debugf("authorized key: %s", k.Type())
			return true
		}
	}

	logrus.Debugf("access denied")
	return false
}

func sftpHandler(sess ssh.Session) {
	debugStream := io.Discard
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(debugStream),
	}
	server, err := sftp.NewServer(
		sess,
		serverOptions...,
	)
	if err != nil {
		logrus.Infof("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); errors.Is(err, io.EOF) {
		server.Close()
		logrus.Infoln("sftp client exited session.")
	} else if err != nil {
		logrus.Infoln("sftp server completed with error:", err)
	}
}
