// Copyright 2022 The envd Authors
// Copyright 2022 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	bitSize = 4096
)

// KeyExists returns true if the okteto key pair exists
func KeyExists(public, private string) bool {
	// public, private := getKeyPaths()
	publicKeyExists, _ := fileutil.FileExists(public)
	if !publicKeyExists {
		logrus.Debugf("%s doesn't exist", public)
		return false
	}

	logrus.Debugf("%s already present", public)

	privateKeyExists, _ := fileutil.FileExists(private)
	if !privateKeyExists {
		logrus.Debugf("%s doesn't exist", private)
		return false
	}

	logrus.Debugf("%s already present", private)
	return true
}

// GenerateKeys generates a SSH key pair on path
func GenerateKeys() error {
	publicKeyPath, privateKeyPath, err := getDefaultKeyPaths()
	if err != nil {
		return err
	}
	return generateKeys(publicKeyPath, privateKeyPath, bitSize)
}

func generateKeys(public, private string, bitSize int) error {
	if KeyExists(public, private) {
		return nil
	}

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return errors.Wrap(err, "failed to generate private SSH key")
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to generate public SSH key")
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	if err := os.WriteFile(public, publicKeyBytes, 0600); err != nil {
		return errors.Wrap(err, "failed to write public SSH key")
	}

	if err := os.WriteFile(private, privateKeyBytes, 0600); err != nil {
		return errors.Wrap(err, "failed to write private SSH key")
	}

	logrus.Debugf("created ssh keypair at  %s and %s", public, private)
	return nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

func getDefaultKeyPaths() (string, string, error) {
	public, err := fileutil.ConfigFile(config.PublicKeyFile)
	if err != nil {
		return "", "", errors.Wrap(err, "Cannot get public key path")
	}

	private, err := fileutil.ConfigFile(config.PrivateKeyFile)
	if err != nil {
		return "", "", errors.Wrap(err, "Cannot get private key path")
	}
	return public, private, nil
}

func DefaultKeyExists() (bool, error) {
	pub, pri, err := getDefaultKeyPaths()
	if err != nil {
		return false, err
	}
	return KeyExists(pub, pri), nil
}

// GetPublicKey returns the path to the public key
func GetPublicKey() (string, error) {
	pub, _, err := getDefaultKeyPaths()
	if err != nil {
		return "", err
	}
	return pub, nil
}

// GetPrivateKey returns the path to the private key
func GetPrivateKey() (string, error) {
	_, pri, err := getDefaultKeyPaths()
	if err != nil {
		return "", err
	}
	return pri, nil
}

// GetPublicKeyOrPanic returns the path to the public key or panic.
func GetPublicKeyOrPanic() string {
	pub, _, err := getDefaultKeyPaths()
	if err != nil {
		logrus.Fatal("Cannot get public key path")
	}
	return pub
}

// GetPrivateKeyOrPanic returns the path to the private key or panic.
func GetPrivateKeyOrPanic() string {
	_, pri, err := getDefaultKeyPaths()
	if err != nil {
		logrus.Fatal("Cannot get private key path")
	}
	return pri
}
