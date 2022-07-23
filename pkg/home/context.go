// Copyright 2022 The envd Authors
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

package home

import (
	"encoding/gob"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type contextManager interface {
	ContextFile() string
	ContextList() (types.EnvdContext, error)
	ContextUse(name string) error
	ContextGetCurrent() (types.BuilderType, string, error)
	ContextCreate(name string,
		builder types.BuilderType, socket string, use bool) error
	ContextRemove(name string) error
}

func (m *generalManager) initContext() error {
	contextFile, err := fileutil.ConfigFile("contexts")
	if err != nil {
		return errors.Wrap(err, "failed to get context file")
	}
	m.contextFile = contextFile

	// Create default context.

	_, err = os.Stat(m.contextFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("filename", m.contextFile).Debug("Creating file")
			file, err := os.Create(m.contextFile)
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			err = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to close file")
			}
			if err := m.dumpContext(); err != nil {
				return errors.Wrap(err, "failed to dump context")
			}
		} else {
			return errors.Wrap(err, "failed to stat file")
		}
	}

	file, err := os.Open(m.contextFile)
	if err != nil {
		return errors.Wrap(err, "failed to open context file")
	}
	defer file.Close()
	e := gob.NewDecoder(file)
	if err := e.Decode(&m.context); err != nil {
		return errors.Wrap(err, "failed to decode context file")
	}
	return nil
}

func (m generalManager) ContextFile() string {
	return m.contextFile
}

func (m generalManager) ContextGetCurrent() (types.BuilderType, string, error) {
	var driver types.BuilderType
	var socket string
	for _, c := range m.context.Contexts {
		if m.context.Current == c.Name {
			driver = c.Builder
			socket = c.BuilderSocket
			return driver, socket, nil
		}
	}
	return "", "", errors.New("no current context")
}

func (m *generalManager) ContextCreate(
	name string, builder types.BuilderType, socket string, use bool) error {
	for _, c := range m.context.Contexts {
		if c.Name == name {
			return errors.Newf("context \"%s\" already exists", name)
		}
	}
	switch builder {
	case types.BuilderTypeDocker, types.BuilderTypeKubernetes, types.BuilderTypeTCP:
		m.context.Contexts = append(m.context.Contexts, types.Context{
			Name:          name,
			Builder:       builder,
			BuilderSocket: socket,
		})
	default:
		return errors.New("unknown builder type")
	}
	if use {
		return m.ContextUse(name)
	}
	return m.dumpContext()
}

func (m *generalManager) ContextRemove(name string) error {
	for i, c := range m.context.Contexts {
		if c.Name == name {
			if m.context.Current == name {
				return errors.Newf("cannot remove current context \"%s\"", name)
			}
			m.context.Contexts = append(
				m.context.Contexts[:i], m.context.Contexts[i+1:]...)
			return m.dumpContext()
		}
	}
	return errors.Newf("cannot find context \"%s\"", name)
}

func (m generalManager) ContextList() (types.EnvdContext, error) {
	return m.context, nil
}

func (m *generalManager) ContextUse(name string) error {
	for _, c := range m.context.Contexts {
		if c.Name == name {
			m.context.Current = name
			return m.dumpContext()
		}
	}
	return errors.Newf("context \"%s\" does not exist", name)
}

func (m *generalManager) dumpContext() error {
	file, err := os.Create(m.contextFile)
	if err != nil {
		return errors.Wrap(err, "failed to create cache status file")
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	if err := e.Encode(m.context); err != nil {
		return errors.Wrap(err, "failed to encode cache map")
	}
	return nil
}
