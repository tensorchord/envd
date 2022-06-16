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
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

type Client interface {
	Attach() error
	ExecWithOutput(cmd string) ([]byte, error)
	Close() error
}

type generalClient struct {
	cli *ssh.Client
}

func NewClient(server, user string,
	port int, auth bool, privateKeyPath, privateKeyPwd string) (Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	}

	var cli *ssh.Client

	if auth {
		// read private key file
		pemBytes, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return nil, errors.Wrapf(err, "reading private key %s failed", privateKeyPath)
		}
		// create signer
		signer, err := signerFromPem(pemBytes, []byte(privateKeyPwd))
		if err != nil {
			return nil, errors.Wrap(err, "creating signer from private key failed")
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	host := fmt.Sprintf("%s:%d", server, port)
	// open connection
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, errors.Wrap(err, "dialing failed")
	}
	cli = conn

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
		// XXX - might want to handle reconnects appending multiple callbacks
		auth := ssh.PublicKeysCallback(forwardingAgent.Signers)
		config.Auth = append(config.Auth, auth)
		if err := agent.ForwardToAgent(cli, forwardingAgent); err != nil {
			return nil, errors.Wrap(err, "forwarding agent to client failed")
		}
	}

	return &generalClient{
		cli: cli,
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

	if err := agent.RequestAgentForwarding(session); err != nil {
		return nil, errors.Wrap(err, "requesting agent forwarding failed")
	}

	return session.Output(cmd)
}

func (c generalClient) Attach() error {
	// open session
	session, err := c.cli.NewSession()
	if err != nil {
		return errors.Wrap(err, "creating session failed")
	}
	defer session.Close()

	if err := agent.RequestAgentForwarding(session); err != nil {
		return errors.Wrap(err, "requesting agent forwarding failed")
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
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stdin for session: %v", err)
	}
	Copy(os.Stdin, stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stdout for session: %v", err)
	}

	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			logrus.Debugf("error while writing to stdOut: %s", err)
		}
	}()

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stderr for session: %v", err)
	}

	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			logrus.Debugf("error while writing to stdOut: %s", err)
		}
	}()

	// TODO(gaocegege): Refactor it to avoid direct access to DefaultGraph
	cmd := shellescape.QuoteCommand([]string{ir.DefaultGraph.Shell})
	logrus.Debugf("executing command over ssh: '%s'", cmd)
	err = session.Run(cmd)
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "status 130") || strings.Contains(err.Error(), "4294967295") {
		return nil
	}
	if strings.Contains(err.Error(), "exit code 137") || strings.Contains(err.Error(), "exit status 137") {
		logrus.Warn(`Insufficient memory.`)
	}

	logrus.Debugf("command failed: %s", err)
	return errors.Wrap(err, "command failed")
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
	//nolint
	if x509.IsEncryptedPEMBlock(pemBlock) {
		// decrypt PEM
		//nolint
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("decrypting PEM block failed %v", err)
		}

		// get RSA, EC or DSA key
		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		// generate signer instance from key
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("creating signer from encrypted key failed %v", err)
		}

		return signer, nil
	} else {
		// generate signer instance from plain key
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing plain private key failed %v", err)
		}

		return signer, nil
	}
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing PKCS private key failed %v", err)
		} else {
			return key, nil
		}
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing EC private key failed %v", err)
		} else {
			return key, nil
		}
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing DSA private key failed %v", err)
		} else {
			return key, nil
		}
	default:
		return nil, fmt.Errorf("Parsing private key failed, unsupported key type %q", block.Type)
	}
}
