// Copyright 2022 The envd Authors
// Copyright 2022 The Okteto Authors
// based on https://github.com/havoc-io/sshconfig
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

package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	sshConfig struct {
		source  []byte
		globals []*param
		hosts   []*host
	}
	host struct {
		comments  []string
		hostnames []string
		params    []*param
	}
	param struct {
		comments []string
		keyword  string
		args     []string
	}
)

//will use auth fields in the future
const (
	forwardAgentKeyword           = "ForwardAgent"
	pubkeyAcceptedKeyTypesKeyword = "PubkeyAcceptedKeyTypes"
	hostKeyword                   = "Host"
	hostNameKeyword               = "HostName"
	portKeyword                   = "Port"
	strictHostKeyCheckingKeyword  = "StrictHostKeyChecking"
	hostKeyAlgorithms             = "HostKeyAlgorithms"
	userKnownHostsFileKeyword     = "UserKnownHostsFile"
	identityFile                  = "IdentityFile"
)

func newHost(hostnames, comments []string) *host {
	return &host{
		comments:  comments,
		hostnames: hostnames,
	}
}

func (h *host) String() string {

	buf := &bytes.Buffer{}

	if len(h.comments) > 0 {
		for _, comment := range h.comments {
			if !strings.HasPrefix(comment, "#") {
				comment = "# " + comment
			}
			fmt.Fprintln(buf, comment)
		}
	}

	fmt.Fprintf(buf, "%s %s\n", hostKeyword, strings.Join(h.hostnames, " "))
	for _, param := range h.params {
		fmt.Fprint(buf, "  ", param.String())
	}

	return buf.String()

}

func newParam(keyword string, args, comments []string) *param {
	return &param{
		comments: comments,
		keyword:  keyword,
		args:     args,
	}
}

func (p *param) String() string {

	buf := &bytes.Buffer{}

	if len(p.comments) > 0 {
		fmt.Fprintln(buf)
		for _, comment := range p.comments {
			if !strings.HasPrefix(comment, "#") {
				comment = "# " + comment
			}
			fmt.Fprintln(buf, comment)
		}
	}

	fmt.Fprintf(buf, "%s %s\n", p.keyword, strings.Join(p.args, " "))

	return buf.String()

}

func (p *param) value() string {
	if len(p.args) > 0 {
		return p.args[0]
	}
	return ""
}

func parse(r io.Reader) (*sshConfig, error) {

	// dat state
	var (
		global = true

		p = &param{}
		h *host
	)

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := &sshConfig{
		source: data,
	}

	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {

		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if line[0] == '#' {
			p.comments = append(p.comments, line)
			continue
		}

		psc := bufio.NewScanner(strings.NewReader(line))
		psc.Split(bufio.ScanWords)
		if !psc.Scan() {
			continue
		}

		p.keyword = psc.Text()

		for psc.Scan() {
			p.args = append(p.args, psc.Text())
		}

		if p.keyword == hostKeyword {
			global = false
			if h != nil {
				config.hosts = append(config.hosts, h)
			}
			h = &host{
				comments:  p.comments,
				hostnames: p.args,
			}
			p = &param{}
			continue
		} else if global {
			config.globals = append(config.globals, p)
			p = &param{}
			continue
		}

		h.params = append(h.params, p)
		p = &param{}

	}

	if global {
		config.globals = append(config.globals, p)
	} else if h != nil {
		config.hosts = append(config.hosts, h)
	}

	return config, nil

}

func (config *sshConfig) writeTo(w io.Writer) error {
	buf := bytes.NewBufferString("")
	for _, param := range config.globals {
		if _, err := fmt.Fprint(buf, param.String()); err != nil {
			return err
		}
	}

	if len(config.globals) > 0 {
		if _, err := fmt.Fprintln(buf); err != nil {
			return err
		}
	}

	for _, host := range config.hosts {
		if _, err := fmt.Fprint(buf, host.String()); err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, buf.String())
	return err
}

func (config *sshConfig) writeToFilepath(p string) error {
	sshDir := filepath.Dir(p)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		logrus.Infof("failed to create SSH directory %s: %s", sshDir, err)
	}

	stat, err := os.Stat(p)
	var mode os.FileMode
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to get info on %s: %s", p, err)
		}

		// default for sshconfig
		mode = 0600
	} else {
		mode = stat.Mode()
	}

	dir := filepath.Dir(p)
	temp, err := os.CreateTemp(dir, "")
	if err != nil {
		return fmt.Errorf("failed to create temporary config file: %s", err)
	}

	defer os.Remove(temp.Name())

	if err := config.writeTo(temp); err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Chmod(temp.Name(), mode); err != nil {
		return fmt.Errorf("failed to set permissions to %s: %s", temp.Name(), err)
	}

	if _, err := getConfig(temp.Name()); err != nil {
		return fmt.Errorf("new config is not valid: %s", err)
	}

	if err := os.Rename(temp.Name(), p); err != nil {
		return fmt.Errorf("failed to move %s to %s: %s", temp.Name(), p, err)
	}

	return nil

}

//nolint:unused
func (config *sshConfig) getHost(hostname string) *host {
	for _, host := range config.hosts {
		for _, hn := range host.hostnames {
			if hn == hostname {
				return host
			}
		}
	}
	return nil
}

func (h *host) getParam(keyword string) *param {
	for _, p := range h.params {
		if p.keyword == keyword {
			return p
		}
	}
	return nil
}

func buildHostname(name string) string {
	return fmt.Sprintf("%s.envd", name)
}

// AddEntry adds an entry to the user's sshconfig
func AddEntry(name, iface string, port int, privateKeyPath string) error {
	return add(getSSHConfigPath(), buildHostname(name), iface, port, privateKeyPath)
}

func add(path, name, iface string, port int, privateKeyPath string) error {
	cfg, err := getConfig(path)
	if err != nil {
		return err
	}

	_ = removeHost(cfg, name)

	// TODO: Use private key to authenticate ssh
	// _, privateKey := getKeyPaths()

	host := newHost([]string{name}, []string{"entry generated by envd"})
	host.params = []*param{
		newParam(forwardAgentKeyword, []string{"yes"}, nil),
		newParam(pubkeyAcceptedKeyTypesKeyword, []string{"+ssh-rsa"}, nil),
		newParam(hostKeyAlgorithms, []string{"+ssh-rsa"}, nil),
		newParam(hostNameKeyword, []string{iface}, nil),
		newParam(portKeyword, []string{strconv.Itoa(port)}, nil),
		newParam(strictHostKeyCheckingKeyword, []string{"no"}, nil),
		newParam(userKnownHostsFileKeyword, []string{"/dev/null"}, nil),
		newParam(identityFile, []string{"\"" + privateKeyPath + "\""}, nil),
	}

	cfg.hosts = append(cfg.hosts, host)
	return save(cfg, path)
}

// RemoveEntry removes the entry to the user's sshconfig if found
func RemoveEntry(name string) error {
	return remove(getSSHConfigPath(), buildHostname(name))
}

// GetPort returns the corresponding SSH port for the dev env
func GetPort(name string) (int, error) {
	cfg, err := getConfig(getSSHConfigPath())
	if err != nil {
		return 0, err
	}

	hostname := buildHostname(name)
	i, found := findHost(cfg, hostname)
	if !found {
		return 0, fmt.Errorf("development container not found")
	}

	param := cfg.hosts[i].getParam(portKeyword)
	if param == nil {
		return 0, fmt.Errorf("port not found")
	}

	port, err := strconv.Atoi(param.value())
	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", param.value())
	}

	return port, nil
}

func remove(path, name string) error {
	cfg, err := getConfig(path)
	if err != nil {
		return err
	}

	if removeHost(cfg, name) {
		return save(cfg, path)
	}

	return nil
}

func removeHost(cfg *sshConfig, name string) bool {
	ix, ok := findHost(cfg, name)
	if ok {
		cfg.hosts = append(cfg.hosts[:ix], cfg.hosts[ix+1:]...)
		return true
	}

	return false
}

func findHost(cfg *sshConfig, name string) (int, bool) {
	for i, h := range cfg.hosts {
		for _, hn := range h.hostnames {
			if hn == name {
				p := h.getParam(portKeyword)
				s := h.getParam(strictHostKeyCheckingKeyword)
				if p != nil && s != nil {
					return i, true
				}
			}
		}
	}

	return 0, false
}

func getConfig(path string) (*sshConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &sshConfig{
				hosts: []*host{},
			}, nil
		}

		return nil, fmt.Errorf("can't open %s: %s", path, err)
	}

	defer f.Close()

	cfg, err := parse(f)
	if err != nil {
		return nil, fmt.Errorf("fail to decode %s: %s", path, err)
	}

	return cfg, nil
}

func save(cfg *sshConfig, path string) error {
	if err := cfg.writeToFilepath(path); err != nil {
		return fmt.Errorf("fail to update SSH config file %s: %w", path, err)
	}

	return nil
}

func getSSHConfigPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	return filepath.Join(dirname, ".ssh", "config")
}
