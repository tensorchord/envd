package syncthing

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	resBody, err := s.ApiCall(GET, "/rest/events", nil, []byte{})
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
	params := map[string]string{
		"since":   strconv.FormatInt(s.latestEventId, 10),
		"timeout": "0",
	}
	resBody, err := s.ApiCall(GET, "/rest/events", params, []byte{})
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
		latestEvent := events[len(events)-1]
		s.latestEventId = latestEvent.Id
	}

	return events, nil
}
