package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEInsightsResolution(id int) (*CSEInsightsResolutionGet, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/insight-resolutions/%d", id), false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEInsightsResolutionResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEInsightsResolutionGet, nil
}

func (s *Client) DeleteCSEInsightsResolution(id int) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/insight-resolutions/%d", id))

	return err
}

func (s *Client) CreateCSEInsightsResolution(CSEInsightsResolutionPost CSEInsightsResolutionPost) (int, error) {

	request := CSEInsightsResolutionRequestPost{
		CSEInsightsResolutionPost: CSEInsightsResolutionPost,
	}

	var response CSEInsightsResolutionResponse

	responseBody, err := s.Post("sec/v1/insight-resolutions", request, false)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return -1, err
	}

	return response.CSEInsightsResolutionGet.ID, nil
}

func (s *Client) UpdateCSEInsightsResolution(CSEInsightsResolutionPost CSEInsightsResolutionPost) error {
	url := fmt.Sprintf("sec/v1/insight-resolutions/%d", CSEInsightsResolutionPost.ID)

	request := CSEInsightsResolutionRequestUpdate{
		CSEInsightsResolutionUpdate{
			Description: CSEInsightsResolutionPost.Description,
		},
	}

	_, err := s.Put(url, request, false)

	return err
}

type CSEInsightsResolutionRequestPost struct {
	CSEInsightsResolutionPost CSEInsightsResolutionPost `json:"fields"`
}

type CSEInsightsResolutionRequestUpdate struct {
	CSEInsightsResolutionUpdate CSEInsightsResolutionUpdate `json:"fields"`
}

type CSEInsightsResolutionResponse struct {
	CSEInsightsResolutionGet CSEInsightsResolutionGet `json:"data"`
}

type CSEInsightsResolutionPost struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	ParentId    int    `json:"parentId,omitempty"`
}

type CSEInsightsResolutionGet struct {
	ID          int                         `json:"id,omitempty"`
	Description string                      `json:"description,omitempty"`
	Name        string                      `json:"name,omitempty"`
	Parent      CSEInsightsResolutionParent `json:"parent,omitempty"`
}

type CSEInsightsResolutionParent struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CSEInsightsResolutionUpdate struct {
	Description string `json:"description"`
}
