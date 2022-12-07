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
	gatewayclient "github.com/moby/buildkit/frontend/gateway/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

const (
	defaultFile = "build.envd"
	defaultFunc = "build"
)

func ImageConfigStr(labels map[string]string, ports map[string]struct{},
	entrypoint []string, env []string, user string) (string, error) {
	pl := platforms.Normalize(platforms.DefaultSpec())
	img := v1.Image{
		Config: v1.ImageConfig{
			Labels:       labels,
			User:         user,
			WorkingDir:   "/",
			Env:          env,
			ExposedPorts: ports,
			Entrypoint:   entrypoint,
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

func parseImportCacheCSV(s string) (gatewayclient.CacheOptionsEntry, error) {
	im := gatewayclient.CacheOptionsEntry{
		Type:  "",
		Attrs: map[string]string{},
	}
	csvReader := csv.NewReader(strings.NewReader(s))
	fields, err := csvReader.Read()
	if err != nil {
		return im, err
	}
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return im, errors.Errorf("invalid value %s", field)
		}
		key := strings.ToLower(parts[0])
		value := parts[1]
		switch key {
		case "type":
			im.Type = value
		default:
			im.Attrs[key] = value
		}
	}
	if im.Type == "" {
		return im, errors.New("--import-cache requires type=<type>")
	}
	return im, nil
}

// ParseImportCache parses --import-cache
func ParseImportCache(importCaches []string) ([]gatewayclient.CacheOptionsEntry, error) {
	var imports []gatewayclient.CacheOptionsEntry
	for _, importCache := range importCaches {
		legacy := !strings.Contains(importCache, "type=")
		if legacy {
			logrus.Warn("--import-cache <ref> is deprecated. Please use --import-cache type=registry,ref=<ref>,<opt>=<optval>[,<opt>=<optval>] instead.")
			imports = append(imports, gatewayclient.CacheOptionsEntry{
				Type:  "registry",
				Attrs: map[string]string{"ref": importCache},
			})
		} else {
			im, err := parseImportCacheCSV(importCache)
			if err != nil {
				return nil, err
			}
			imports = append(imports, im)
		}
	}
	return imports, nil
}

// ParseExportCache parses --export-cache (and legacy --export-cache-opt)
// Refer to github.com/moby/buildkit/cmd/buildctl/build/exportcache.go
func ParseExportCache(exportCaches, legacyExportCacheOpts []string) ([]client.CacheOptionsEntry, error) {
	var exports []client.CacheOptionsEntry
	if len(legacyExportCacheOpts) > 0 {
		if len(exportCaches) != 1 {
			return nil, errors.New("--export-cache-opt requires exactly single --export-cache")
		}
	}
	for _, exportCache := range exportCaches {
		if len(exportCache) <= 0 {
			continue
		}
		legacy := !strings.Contains(exportCache, "type=")
		if legacy {
			logrus.Warnf("--export-cache <ref> --export-cache-opt <opt>=<optval> is deprecated. Please use --export-cache type=registry,ref=<ref>,<opt>=<optval>[,<opt>=<optval>] instead")
			attrs, err := attrMap(legacyExportCacheOpts)
			if err != nil {
				return nil, err
			}
			if _, ok := attrs["mode"]; !ok {
				attrs["mode"] = "min"
			}
			attrs["ref"] = exportCache
			exports = append(exports, client.CacheOptionsEntry{
				Type:  "registry",
				Attrs: attrs,
			})
		} else {
			if len(legacyExportCacheOpts) > 0 {
				return nil, errors.New("--export-cache-opt is not supported for the specified --export-cache. Please use --export-cache type=<type>,<opt>=<optval>[,<opt>=<optval>] instead")
			}
			ex, err := parseExportCacheCSV(exportCache)
			if err != nil {
				return nil, err
			}
			exports = append(exports, ex)
		}
	}
	return exports, nil
}

func attrMap(sl []string) (map[string]string, error) {
	m := map[string]string{}
	for _, v := range sl {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			return nil, errors.Errorf("invalid value %s", v)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func parseExportCacheCSV(s string) (client.CacheOptionsEntry, error) {
	ex := client.CacheOptionsEntry{
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
			ex.Type = value
		default:
			ex.Attrs[key] = value
		}
	}
	if ex.Type == "" {
		return ex, errors.New("--export-cache requires type=<type>")
	}
	if _, ok := ex.Attrs["mode"]; !ok {
		ex.Attrs["mode"] = "min"
	}
	return ex, nil
}

// parseOutput parses --output
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

func ParseFromStr(fromStr string) (string, string, error) {
	filename := defaultFile
	funcname := defaultFunc
	if !strings.Contains(fromStr, ":") {
		if len(fromStr) > 0 {
			filename = fromStr
		}
		return filename, funcname, nil
	}

	fromArr := strings.Split(fromStr, ":")

	if len(fromArr) != 2 {
		return "", "", errors.New("invalid from format, expected `file:func`")
	}
	if fromArr[0] != "" {
		filename = fromArr[0]
	}
	if fromArr[1] != "" {
		funcname = fromArr[1]
	}
	return filename, funcname, nil
}
