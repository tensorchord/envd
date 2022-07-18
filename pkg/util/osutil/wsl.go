// Copyright 2022 The envd Authors
// Copyright 2022 mateors
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

package osutil

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

func IsWsl() bool {
	// Return false if meet error
	cmd := exec.Command("cat", "/proc/version")
	output, err := cmd.Output()
	if err != nil {
		logrus.Debugf("Error when check whether sys is WSL: %v", err)
		return false
	}

	return strings.Contains(strings.ToLower(string(output)), "microsoft")
}

func GetWslHostSshConfig() (string, error) {
	userCmd := exec.Command("wslvar", "USERPROFILE")
	userOutput, err := userCmd.Output()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("wslpath", string(userOutput))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	outputPath := path.Join(strings.Trim(string(output), "\n"), ".ssh", "config")
	logrus.Debugf("wsl sshconfig path: %s", outputPath)
	return outputPath, nil
}

func GetWslIp() (string, error) {
	ip, err := getInterfaceIpv4Addr("eth0")
	if err != nil {
		return "", err
	}
	return ip, nil
}

func GetWindowsEnvdConfigHome() (string, error) {

	userCmd := exec.Command("wslvar", "LOCALAPPDATA")
	userOutput, err := userCmd.Output()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("wslpath", string(userOutput))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	envdDir := filepath.Join(strings.Trim(string(output), "\n"), "envd")
	if err := os.MkdirAll(envdDir, 0755); err != nil {
		return "", err
	}
	return envdDir, nil
}

// from: https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
func getInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}

func CopyToWinEnvdHome(src string, permission os.FileMode) (string, error) {
	// Return dst path in windows format
	winhome, err := GetWindowsEnvdConfigHome()
	if err != nil {
		return "", err
	}
	filename := filepath.Base(src)
	dst := filepath.Join(winhome, filename)
	err = copy(src, dst, permission)
	if err != nil {
		return "", err
	}

	envdDirWinCmd := exec.Command("wslpath", "-w", dst)
	winDir, err := envdDirWinCmd.Output()

	if err != nil {
		return "", err
	}
	return strings.Trim(string(winDir), "\n"), nil
}

func copy(src, dst string, permission os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	err = os.Chmod(dst, permission)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
