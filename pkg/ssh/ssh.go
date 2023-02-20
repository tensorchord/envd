// Copyright 2022 The envd Authors
// Copyright 2022 The okteto Authors
// Copyright 2022 stefanprodan
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

// https://gist.github.com/stefanprodan/2d20d0c6fdab6f14ce8219464e8b4b9a
// Refer to okteto/pkg/ssh/exec.go

package ssh

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"

	"github.com/tensorchord/envd/pkg/ssh/config"
)

type Client interface {
	Attach() error
	ExecWithOutput(cmd string) ([]byte, error)
	LocalForward(localAddress, targetAddress string) error
	RemoteForward(localAddress, targetAddress string) error
	Close() error
}

type Options struct {
	AgentForwarding bool
	Server          string
	User            string
	Port            int
	Auth            bool
	PrivateKeyPath  string
	PrivateKeyPwd   string
}

func DefaultOptions() Options {
	return Options{
		User:            "envd",
		Auth:            true,
		PrivateKeyPwd:   "",
		AgentForwarding: true,
	}
}

func GetOptions(entry string) (*Options, error) {
	path, err := config.GetPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "getting private key failed")
	}
	port, err := config.GetPort(entry)
	if err != nil {
		return nil, errors.Wrap(err, "getting port failed")
	}
	// TODO(gaocegege): Make it configurable.
	opt := DefaultOptions()
	opt.Port = port
	opt.PrivateKeyPath = path
	return &opt, nil
}

type generalClient struct {
	cli *ssh.Client
	opt *Options
}

func NewClient(opt Options) (Client, error) {
	logrus.WithFields(logrus.Fields{
		"user":             opt.User,
		"port":             opt.Port,
		"server":           opt.Server,
		"agent-forwarding": opt.AgentForwarding,
		"auth":             opt.Auth,
	}).Debug("ssh to the environment")

	config := &ssh.ClientConfig{
		User: opt.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	}

	var cli *ssh.Client

	if opt.Auth {
		// read private key file
		pemBytes, err := os.ReadFile(opt.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrapf(
				err, "reading private key %s failed", opt.PrivateKeyPath)
		}
		// create signer
		signer, err := signerFromPem(pemBytes, []byte(opt.PrivateKeyPwd))
		if err != nil {
			return nil, errors.Wrap(err, "creating signer from private key failed")
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	host := fmt.Sprintf("%s:%d", opt.Server, opt.Port)
	// open connection
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, errors.Wrap(err, "dialing failed")
	}
	cli = conn

	if opt.AgentForwarding {
		// open connection to the local agent
		socketLocation := os.Getenv("SSH_AUTH_SOCK")
		if socketLocation != "" {
			agentConn, err := net.Dial("unix", socketLocation)
			if err != nil {
				return nil, errors.Wrap(err, "could not connect to local agent socket")
			}
			// create agent and add in auth
			forwardingAgent := agent.NewClient(agentConn)
			// add callback for forwarding agent to SSH config
			// might want to handle reconnects appending multiple callbacks
			auth := ssh.PublicKeysCallback(forwardingAgent.Signers)
			config.Auth = append(config.Auth, auth)
			if err := agent.ForwardToAgent(cli, forwardingAgent); err != nil {
				return nil, errors.Wrap(err, "forwarding agent to client failed")
			}
		} else {
			logrus.Warn("SSH Agent Forwarding is disabled. This will have no impact on your normal use if you do not use the ssh key on the host.")
		}
	}

	return &generalClient{
		cli: cli,
		opt: &opt,
	}, nil
}

func (c generalClient) Close() error {
	return c.cli.Close()
}

func (c generalClient) ExecWithOutput(cmd string) ([]byte, error) {
	defer c.cli.Close()

	// open session
	session, err := c.cli.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "creating session failed")
	}
	defer session.Close()

	if c.opt.AgentForwarding {
		if err := agent.RequestAgentForwarding(session); err != nil {
			return nil, errors.Wrap(err, "requesting agent forwarding failed")
		}
	}

	return session.CombinedOutput(cmd)
}

func (c generalClient) Attach() error {
	// open session
	session, err := c.cli.NewSession()
	if err != nil {
		return errors.Wrap(err, "creating session failed")
	}
	defer session.Close()

	if c.opt.AgentForwarding {
		if err := agent.RequestAgentForwarding(session); err != nil {
			return errors.Wrap(err, "requesting agent forwarding failed")
		}
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,      // Disable echoing
		ssh.ECHOCTL:       0,      // Don't print control chars
		ssh.IGNCR:         1,      // Ignore CR on input
		ssh.TTY_OP_ISPEED: 115200, // baud in
		ssh.TTY_OP_OSPEED: 115200, // baud out
	}

	height, width := 80, 40
	var termFD int
	var ok bool
	if termFD, ok = isTerminal(os.Stdin); ok {
		width, height, err = term.GetSize(int(os.Stdout.Fd()))
		logrus.Debugf("terminal width %d height %d", width, height)
		if err != nil {
			logrus.Debugf("request for terminal size failed: %s", err)
		}
	}

	state, err := term.MakeRaw(termFD)
	if err != nil {
		logrus.Debugf("request for raw terminal failed: %s", err)
	}

	defer func() {
		if state == nil {
			return
		}

		if err := term.Restore(termFD, state); err != nil {
			logrus.Debugf("failed to restore terminal: %s", err)
		}

		logrus.Debugf("terminal restored")
	}()

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return errors.Newf("request for pseudo terminal failed: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	logrus.Debug("starting shell")
	err = session.Shell()
	if err != nil {
		return errors.Wrap(err, "starting shell failed")
	}
	logrus.Debug("waiting for shell to exit")
	if err = session.Wait(); err != nil {
		var ee *ssh.ExitError
		if ok := errors.As(err, &ee); ok {
			switch ee.ExitStatus() {
			case 130:
				return nil
			case 137:
				logrus.Warn(`Insufficient memory.`)
			}
		}
		var emr *ssh.ExitMissingError
		if ok := errors.As(err, &emr); ok {
			logrus.Debugf("exit status missing: %s", emr)
			return nil
		}
		return errors.Wrap(err, "waiting for session failed")
	}

	logrus.Debug("shell exited")
	return nil
}

func (c generalClient) LocalForward(localAddress, targetAddress string) error {
	localListener, err := net.Listen("tcp", localAddress)
	if err != nil {
		return errors.Wrap(err, "net.Listen failed")
	}

	logrus.Debug("begin to local forward " + localAddress + " to " + targetAddress)
	for {
		localCon, err := localListener.Accept()
		if err != nil {
			return errors.Wrap(err, "listen.Accept failed")
		}

		sshConn, err := c.cli.Dial("tcp", targetAddress)
		if err != nil {
			return errors.Wrap(err, "listen.Accept failed")
		}

		// Copy local.Reader to sshConn.Writer
		go func() {
			_, err = io.Copy(sshConn, localCon)
			if err != nil {
				logrus.Debugf("io.Copy failed: %v", err)
			}
		}()

		// Copy sshConn.Reader to localCon.Writer
		go func() {
			_, err = io.Copy(localCon, sshConn)
			if err != nil {
				logrus.Debugf("io.Copy failed: %v", err)
			}
		}()
	}
}

func (c generalClient) RemoteForward(remoteAddress, targetAddress string) error {
	sshListener, err := c.cli.Listen("tcp", remoteAddress)
	if err != nil {
		return errors.Wrap(err, "cli.Listen failed")
	}

	logrus.Debug("begin to remote forward " + remoteAddress + " to " + targetAddress)
	for {
		sshCon, err := sshListener.Accept()
		if err != nil {
			return errors.Wrap(err, "listen.Accept failed")
		}

		targetCon, err := net.Dial("tcp", targetAddress)
		if err != nil {
			return errors.Wrap(err, "net.Dial failed")
		}

		// Copy sshCon.Reader to targetCon.Writer
		go func() {
			_, err = io.Copy(targetCon, sshCon)
			if err != nil {
				logrus.Debugf("io.Copy failed: %v", err)
			}
		}()

		// Copy targetCon.Reader to sshCon.Writer
		go func() {
			_, err = io.Copy(sshCon, targetCon)
			if err != nil {
				logrus.Debugf("io.Copy failed: %v", err)
			}
		}()
	}
}

func isTerminal(r io.Reader) (int, bool) {
	switch v := r.(type) {
	case *os.File:
		return int(v.Fd()), term.IsTerminal(int(v.Fd()))
	default:
		return 0, false
	}
}

func signerFromPem(pemBytes []byte, password []byte) (ssh.Signer, error) {
	// read pem block
	err := errors.New("Pem decode failed, no key found")
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, err
	}

	// handle encrypted key
	// nolint
	if x509.IsEncryptedPEMBlock(pemBlock) {
		// decrypt PEM
		// nolint
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, []byte(password))
		if err != nil {
			return nil, errors.Newf("decrypting PEM block failed %w", err)
		}

		// get RSA, EC or DSA key
		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		// generate signer instance from key
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, errors.Newf("creating signer from encrypted key failed %w", err)
		}

		return signer, nil
	} else {
		// generate signer instance from plain key
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, errors.Newf("parsing plain private key failed %w", err)
		}

		return signer, nil
	}
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Newf("Parsing PKCS private key failed %w", err)
		}
		return key, nil
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Newf("Parsing EC private key failed %w", err)
		}
		return key, nil
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Newf("Parsing DSA private key failed %w", err)
		}
		return key, nil
	default:
		return nil, errors.Newf("Parsing private key failed, unsupported key type %q", block.Type)
	}
}
