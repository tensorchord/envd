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

package telemetry

import (
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	segmentio "github.com/segmentio/analytics-go/v3"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/version"
)

type TelemetryField func(*segmentio.Properties)

type Reporter interface {
	Telemetry(command string, fields ...TelemetryField)
}

type defaultReporter struct {
	client        segmentio.Client
	telemetryFile string
	enabled       bool

	UID string
}

var (
	reporter *defaultReporter
	once     sync.Once
)

func Initialize(enabled bool, token string) error {
	once.Do(func() {
		// Ref https://segment.com/docs/connections/sources/catalog/libraries/server/go/#development-settings
		c, err := segmentio.NewWithConfig(token, segmentio.Config{
			BatchSize: 1,
		})
		if err != nil {
			panic(err)
		}
		reporter = &defaultReporter{
			enabled: enabled,
			client:  c,
		}
	})
	return reporter.init()
}

func GetReporter() Reporter {
	return reporter
}

// init gets the UID from .config or create one if not existed.
func (r *defaultReporter) init() error {
	// Create $HOME/.config/envd/telemetry
	tfile, err := fileutil.ConfigFile("telemetry")
	if err != nil {
		return errors.Wrap(err, "failed to get telemetry file")
	}

	r.telemetryFile = tfile
	_, err = os.Stat(r.telemetryFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("filename", r.telemetryFile).Debug("Creating file")
			file, err := os.Create(r.telemetryFile)
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			err = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to close file")
			}
			r.UID = uuid.New().String()
			if err := r.dumpTelemetry(); err != nil {
				return errors.Wrap(err, "failed to dump auth")
			}
		} else {
			return errors.Wrap(err, "failed to stat file")
		}
	}

	file, err := os.Open(r.telemetryFile)
	if err != nil {
		return errors.Wrap(err, "failed to open telemetry file")
	}

	// Read uid to file.
	uid, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read telemetry file")
	}
	file.Close()
	r.UID = string(uid)

	r.Identify()
	return nil
}

func (r *defaultReporter) dumpTelemetry() error {
	file, err := os.Create(r.telemetryFile)
	if err != nil {
		return errors.Wrap(err, "failed to create cache telemetry file")
	}
	defer file.Close()

	// Write uid to file.
	_, err = file.Write([]byte(r.UID))
	return err
}

func (r *defaultReporter) Identify() {
	logrus.WithField("UID", r.UID).Debug("telemetry initialization")
	if r.enabled {
		logrus.Debug("sending telemetry")
		v := version.GetVersion()
		if err := r.client.Enqueue(segmentio.Identify{
			UserId: r.UID,
			Context: &segmentio.Context{
				OS: segmentio.OSInfo{
					Name:    runtime.GOOS,
					Version: runtime.GOARCH,
				},
				App: segmentio.AppInfo{
					Name:    "envd-cli",
					Version: v.Version,
				},
			},
			Timestamp: time.Now(),
			Traits:    segmentio.NewTraits(),
		}); err != nil {
			logrus.Warn("telemetry failed")
			return
		}
	}
}

func AddField(name string, value interface{}) TelemetryField {
	return func(p *segmentio.Properties) {
		p.Set(name, value)
	}
}

func (r *defaultReporter) Telemetry(command string, fields ...TelemetryField) {
	if r.enabled {
		logrus.WithFields(logrus.Fields{
			"UID":     r.UID,
			"command": command,
		}).Debug("sending telemetry track event")
		t := segmentio.Track{
			UserId:     r.UID,
			Event:      command,
			Properties: segmentio.NewProperties(),
		}
		for _, field := range fields {
			field(&t.Properties)
		}
		if err := r.client.Enqueue(t); err != nil {
			logrus.Warn(err)
		}
		// make sure the msg can be sent out
		r.client.Close()
	}
}
