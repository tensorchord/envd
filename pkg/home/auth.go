package home

import (
	"encoding/gob"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type authManager interface {
	AuthFile() string
	AuthGetCurrent() (types.AuthConfig, error)
	AuthCreate(ac types.AuthConfig, use bool) error
	AuthUse(name string) error
}

func (m *generalManager) initAuth() error {
	// Create $HOME/.config/envd/auth.json
	auth, err := fileutil.ConfigFile("auth")
	if err != nil {
		return errors.Wrap(err, "failed to get auth file")
	}

	m.authFile = auth

	_, err = os.Stat(m.authFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("filename", m.authFile).Debug("Creating file")
			file, err := os.Create(m.authFile)
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			err = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to close file")
			}
			if err := m.dumpAuth(); err != nil {
				return errors.Wrap(err, "failed to dump auth")
			}
		} else {
			return errors.Wrap(err, "failed to stat file")
		}
	}

	file, err := os.Open(m.authFile)
	if err != nil {
		return errors.Wrap(err, "failed to open auth file")
	}
	defer file.Close()
	e := gob.NewDecoder(file)
	if err := e.Decode(&m.auth); err != nil {
		return errors.Wrap(err, "failed to decode auth file")
	}
	return nil
}

func (m *generalManager) AuthFile() string {
	return m.authFile
}

func (m *generalManager) AuthGetCurrent() (types.AuthConfig, error) {
	for _, c := range m.auth.Auth {
		if c.Name == m.auth.Current {
			return c, nil
		}
	}
	return types.AuthConfig{}, errors.New("cannot find the current auth config")
}

func (m *generalManager) AuthCreate(ac types.AuthConfig, use bool) error {
	for _, a := range m.auth.Auth {
		if a.Name == ac.Name {
			// Auth should be idempotent. Thus do not return error here.
			return nil
		}
	}
	m.auth.Auth = append(m.auth.Auth, ac)
	if use {
		return m.AuthUse(ac.Name)
	}
	return m.dumpAuth()
}

func (m *generalManager) AuthUse(name string) error {
	for _, a := range m.auth.Auth {
		if a.Name == name {
			m.auth.Current = name
			return m.dumpAuth()
		}
	}
	return errors.Newf("auth config \"%s\" does not exist", name)
}

func (m *generalManager) dumpAuth() error {
	file, err := os.Create(m.authFile)
	if err != nil {
		return errors.Wrap(err, "failed to create cache auth file")
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	if err := e.Encode(m.auth); err != nil {
		return errors.Wrap(err, "failed to encode auth file")
	}
	return nil
}
