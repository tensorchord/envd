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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
)

// Status represents the status of a syncthing folder.
type Status struct {
	State      string `json:"state"`
	PullErrors int64  `json:"pullErrors"`
}

const ConfigSaved = "ConfigSaved"

type ConfigSavedEvent struct {
	Type     string               `json:"type"`
	Id       int64                `json:"id"`
	GlobalId int64                `json:"globalID"`
	Time     string               `json:"time"`
	Data     config.Configuration `json:"data"`
}

type GeneralEvent struct {
	Type     string      `json:"type"`
	Id       int64       `json:"id"`
	GlobalId int64       `json:"globalID"`
	Time     string      `json:"time"`
	Data     interface{} `json:"data"`
}

// Fetches the most recent event using the syncthing rest api
func (s *Syncthing) GetMostRecentEvent() (*GeneralEvent, error) {
	resBody, err := s.Client.SendRequest(GET, "/rest/events", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get most recent event: %w", err)
	}

	var events []*GeneralEvent
	err = json.Unmarshal(resBody, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal most recent event: %w", err)
	}

	latestEvent := events[len(events)-1]
	s.latestEventId = latestEvent.Id

	// Assuming that the events are returned in order
	return events[len(events)-1], nil
}

// Fetches the latest config saved events using the syncthing rest api starting from the latest event id
func (s *Syncthing) GetConfigSavedEvents() ([]*ConfigSavedEvent, error) {
	logrus.Debugf("Getting config saved events")
	params := map[string]string{
		"since":   strconv.FormatInt(s.latestEventId, 10),
		"timeout": "0",
	}
	resBody, err := s.Client.SendRequest(GET, "/rest/events", params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get config saved event: %w", err)
	}

	var allEvents []*ConfigSavedEvent
	err = json.Unmarshal(resBody, &allEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config saved event: %w", err)
	}

	var events []*ConfigSavedEvent
	for _, event := range allEvents {
		if event.Type == ConfigSaved {
			events = append(events, event)
		}
	}

	if len(allEvents) > 0 {
		latestEvent := allEvents[len(allEvents)-1]
		s.latestEventId = latestEvent.Id
	}

	return events, nil
}
