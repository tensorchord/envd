package syncthing

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"
)

func NewApiClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
}

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

// Makes API calls to the syncthing instance's rest api
func (s *Syncthing) ApiCall(method string, url string, params map[string]string, body []byte) ([]byte, error) {
	// TODO: can implement retry logic

	var urlPath = path.Join(fmt.Sprintf("localhost:%s", s.Port), url)

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s", urlPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize syncthing API request: %w", err)
	}

	req.Header.Set("X-API-Key", s.ApiKey)

	q := req.URL.Query()

	for key, value := range params {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := s.Client.Do(req)
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
