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

package vscode

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/ziputil"
)

const (
	cacheKeyPrefix = "vscode-plugins"
)

type Client interface {
	DownloadOrCache(plugin Plugin) (bool, error)
	PluginPath(p Plugin) string
}

type generalClient struct {
	vendor MarketplaceVendor
	logger *logrus.Entry
}

func NewClient(vendor MarketplaceVendor) (Client, error) {
	switch vendor {
	case MarketplaceVendorOpenVSX:
		return &generalClient{
			vendor: vendor,
			logger: logrus.WithField("vendor", MarketplaceVendorOpenVSX),
		}, nil
	case MarketplaceVendorVSCode:
		return &generalClient{
			vendor: vendor,
			logger: logrus.WithField("vendor", MarketplaceVendorVSCode),
		}, nil
	default:
		return nil, errors.Errorf("unknown marketplace vendor %s", vendor)
	}
}

func (c generalClient) PluginPath(p Plugin) string {
	if p.Version != nil {
		return fmt.Sprintf("%s.%s-%s/extension/", p.Publisher, p.Extension, *p.Version)

	}
	return fmt.Sprintf("%s.%s/extension/", p.Publisher, p.Extension)
}

func unzipPath(p Plugin) string {
	if p.Version != nil {
		return fmt.Sprintf("%s/%s.%s-%s", home.GetManager().CacheDir(),
			p.Publisher, p.Extension, *p.Version)
	}
	return fmt.Sprintf("%s/%s.%s", home.GetManager().CacheDir(),
		p.Publisher, p.Extension)
}

// DownloadOrCache downloads or cache the plugin.
// If the plugin is already downloaded, it returns true.
func (c generalClient) DownloadOrCache(p Plugin) (bool, error) {
	cacheKey := fmt.Sprintf("%s-%s", cacheKeyPrefix, p)
	if home.GetManager().Cached(cacheKey) {
		logrus.WithFields(logrus.Fields{
			"cache": cacheKey,
		}).Debugf("vscode plugin %s already exists in cache", p)
		return true, nil
	}

	var url, filename string
	if c.vendor == MarketplaceVendorVSCode {
		if p.Version == nil {
			return false, errors.New("version is required for vscode marketplace")
		}
		// TODO(gaocegege): Support version auto-detection.
		url = fmt.Sprintf(vendorVSCodeTemplate,
			p.Publisher, p.Publisher, p.Extension, *p.Version)
		filename = fmt.Sprintf("%s/%s.%s-%s.vsix", home.GetManager().CacheDir(),
			p.Publisher, p.Extension, *p.Version)
	} else {
		var err error
		url, err = GetLatestVersionURL(p)
		if err != nil {
			return false, errors.Wrap(err, "failed to get latest version url")
		}
		filename = fmt.Sprintf("%s/%s.%s.vsix", home.GetManager().CacheDir(),
			p.Publisher, p.Extension)
	}

	logger := logrus.WithFields(logrus.Fields{
		"publisher": p.Publisher,
		"extension": p.Extension,
		"version":   p.Version,
		"url":       url,
		"file":      filename,
	})

	logger.Debugf("downloading vscode plugin %s", p)
	out, err := os.Create(filename)

	if err != nil {
		return false, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	logger.Debugf("downloading vscode plugin")

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	_, err = ziputil.Unzip(filename, unzipPath(p))
	if err != nil {
		return false, errors.Wrap(err, "failed to unzip")
	}

	if err := home.GetManager().MarkCache(cacheKey, true); err != nil {
		return false, errors.Wrap(err, "failed to update cache status")
	}
	return false, nil
}
