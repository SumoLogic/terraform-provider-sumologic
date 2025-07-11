package sumologic

import (
	"encoding/json"
	"fmt"
)

type LocalWindowsEventLogSource struct {
	Source
	LogNames       []interface{} `json:"logNames"`
	RenderMessages bool          `json:"renderMessages"`
	EventFormat    int           `json:"eventFormat"`
	EventMessage   *int          `json:"eventMessage,omitempty"`
	// Allowlist      []string      `json:"allowlist,omitempty"`
	DenyList string `json:"denylist,omitempty"`
}

func (s *Client) CreateLocalWindowsEventLogSource(source LocalWindowsEventLogSource, collectorID int) (int, error) {

	type LocalWindowsEventLogSourceMessage struct {
		Source LocalWindowsEventLogSource `json:"source"`
	}

	request := LocalWindowsEventLogSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response LocalWindowsEventLogSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetLocalWindowsEventLogSource(collectorID, sourceID int) (*LocalWindowsEventLogSource, error) {
	body, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID))
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type LocalWindowsEventLogSourceResponse struct {
		Source LocalWindowsEventLogSource `json:"source"`
	}

	var response LocalWindowsEventLogSourceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil

}

func (s *Client) UpdateLocalWindowsEventLogSource(source LocalWindowsEventLogSource, collectorID int) error {

	type LocalWindowsEventLogMessage struct {
		Source LocalWindowsEventLogSource `json:"source"`
	}

	request := LocalWindowsEventLogMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request)

	return err
}
