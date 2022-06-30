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

package builder

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/containerd/console"
	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

func ImageConfigStr(labels map[string]string) (string, error) {
	pl := platforms.Normalize(platforms.DefaultSpec())
	img := v1.Image{
		Config: v1.ImageConfig{
			Labels:     labels,
			WorkingDir: "/",
			Env:        []string{"PATH=" + DefaultPathEnv(pl.OS)},
		},
		Architecture: pl.Architecture,
		// Refer to https://github.com/tensorchord/envd/issues/269#issuecomment-1152944914
		OS: "linux",
		RootFS: v1.RootFS{
			Type: "layers",
		},
	}
	data, err := json.Marshal(img)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DefaultPathEnvUnix is unix style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ':' character .
const DefaultPathEnvUnix = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/conda/bin:/usr/local/julia/bin"

// DefaultPathEnvWindows is windows style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ';' character .
const DefaultPathEnvWindows = "c:\\Windows\\System32;c:\\Windows"

func DefaultPathEnv(os string) string {
	if os == "windows" {
		return DefaultPathEnvWindows
	}
	return DefaultPathEnvUnix
}

// ParseOutput parses --output
// Refer to https://github.com/moby/buildkit/blob/master/cmd/buildctl/build/output.go#L56
func parseOutput(exports string) ([]client.ExportEntry, error) {
	var entries []client.ExportEntry
	if exports == "" {
		return entries, nil
	}

	e, err := parseOutputCSV(exports)
	if err != nil {
		return nil, err
	}
	entries = append(entries, e)
	return entries, nil
}

// parseOutputCSV parses a single --output CSV string
func parseOutputCSV(s string) (client.ExportEntry, error) {
	ex := client.ExportEntry{
		Type:  "",
		Attrs: map[string]string{},
	}
	csvReader := csv.NewReader(strings.NewReader(s))
	fields, err := csvReader.Read()
	if err != nil {
		return ex, err
	}
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return ex, errors.Errorf("invalid value %s", field)
		}
		key := strings.ToLower(parts[0])
		value := parts[1]
		switch key {
		case "type":
			logrus.WithFields(logrus.Fields{
				"type": value,
			}).Debug("Adding type into exporter entry")
			ex.Type = value
		default:
			logrus.WithFields(logrus.Fields{
				"key":   key,
				"value": value,
			}).Debug("Adding key into exporter entry")
			ex.Attrs[key] = value
		}
	}
	if ex.Type == "" {
		return ex, errors.New("--output requires type=<type>")
	}
	if v, ok := ex.Attrs["output"]; ok {
		return ex, errors.Errorf("output=%s not supported for --output, you meant dest=%s?", v, v)
	}
	ex.Output, ex.OutputDir, err = resolveExporterDest(ex.Type, ex.Attrs["dest"])
	if err != nil {
		return ex, errors.Wrap(err, "invalid output option: output")
	}
	if ex.Output != nil || ex.OutputDir != "" {
		delete(ex.Attrs, "dest")
	}
	return ex, nil
}

// resolveExporterDest returns at most either one of io.WriteCloser (single file) or a string (directory path).
func resolveExporterDest(exporter, dest string) (func(map[string]string) (io.WriteCloser, error), string, error) {
	wrapWriter := func(wc io.WriteCloser) func(map[string]string) (io.WriteCloser, error) {
		return func(m map[string]string) (io.WriteCloser, error) {
			return wc, nil
		}
	}
	switch exporter {
	case client.ExporterLocal:
		if dest == "" {
			return nil, "", errors.New("output directory is required for local exporter")
		}
		return nil, dest, nil
	case client.ExporterOCI, client.ExporterDocker, client.ExporterTar:
		if dest != "" && dest != "-" {
			fi, err := os.Stat(dest)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return nil, "", errors.Wrapf(err, "invalid destination file: %s", dest)
			}
			if err == nil && fi.IsDir() {
				return nil, "", errors.Errorf("destination file is a directory")
			}
			w, err := os.Create(dest)
			return wrapWriter(w), "", err
		}
		// if no output file is specified, use stdout
		if _, err := console.ConsoleFromFile(os.Stdout); err == nil {
			return nil, "", errors.Errorf("output file is required for %s exporter. refusing to write to console", exporter)
		}
		return wrapWriter(os.Stdout), "", nil
	default: // e.g. client.ExporterImage
		if dest != "" {
			return nil, "", errors.Errorf("output %s is not supported by %s exporter", dest, exporter)
		}
		return nil, "", nil
	}
}
