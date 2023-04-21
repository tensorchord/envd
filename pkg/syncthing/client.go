// Copyright 2023 The envd Authors
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

package syncthing

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewAPIClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
}

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

type Client struct {
	ApiKey   string
	Client   *http.Client
	BasePath string
}

func (s *Syncthing) NewClient() *Client {
	return &Client{
		ApiKey:   s.ApiKey,
		Client:   NewAPIClient(),
		BasePath: fmt.Sprintf("http://127.0.0.1:%s", s.Port),
	}
}

// Makes API calls to the syncthing instance's rest api
func (c *Client) SendRequest(method string, url string, params map[string]string, body []byte) ([]byte, error) {
	logrus.Debug("calling syncthing API: ", url)
	// TODO: can implement retry logic

	var urlPath = c.BasePath + url

	req, err := http.NewRequest(method, urlPath, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize syncthing API request: %w", err)
	}

	req.Header.Set("X-API-Key", c.ApiKey)

	q := req.URL.Query()

	for key, value := range params {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call syncthing [%s]: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to call syncthing [%s]: %s", url, resp.Status)
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from syncthing [%s]: %w", url, err)
	}

	return body, nil
}
