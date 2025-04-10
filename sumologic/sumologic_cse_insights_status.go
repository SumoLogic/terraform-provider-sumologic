package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEInsightsStatus(id string) (*CSEInsightsStatusGet, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/insight-status/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEInsightsStatusResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEInsightsStatusGet, nil
}

func (s *Client) DeleteCSEInsightsStatus(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/insight-status/%s", id))

	return err
}

func (s *Client) CreateCSEInsightsStatus(CSEInsightsStatusPost CSEInsightsStatusPost) (string, error) {

	request := CSEInsightsStatusRequestPost{
		CSEInsightsStatusPost: CSEInsightsStatusPost,
	}

	var response CSEInsightsStatusResponse

	responseBody, err := s.Post("sec/v1/insight-status", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEInsightsStatusGet.ID, nil
}

func (s *Client) UpdateCSEInsightsStatus(CSEInsightsStatusPost CSEInsightsStatusPost) error {
	url := fmt.Sprintf("sec/v1/insight-status/%s", CSEInsightsStatusPost.ID)

	request := CSEInsightsStatusRequestUpdate{
		CSEInsightsStatusUpdate{
			Description: CSEInsightsStatusPost.Description,
			Name:        CSEInsightsStatusPost.Name,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSEInsightsStatusRequestPost struct {
	CSEInsightsStatusPost CSEInsightsStatusPost `json:"fields"`
}

type CSEInsightsStatusRequestUpdate struct {
	CSEInsightsStatusUpdate CSEInsightsStatusUpdate `json:"fields"`
}

type CSEInsightsStatusResponse struct {
	CSEInsightsStatusGet CSEInsightsStatusGet `json:"data"`
}

type CSEInsightsStatusPost struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

type CSEInsightsStatusGet struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type CSEInsightsStatusUpdate struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}
