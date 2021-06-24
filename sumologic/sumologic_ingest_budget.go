package sumologic

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type IngestBudget struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	FieldValue  string `json:"fieldValue"`
	Capacity    int    `json:"capacityBytes"`
	Timezone    string `json:"timezone"`
	ResetTime   string `json:"resetTime"`
	Description string `json:"description,omitempty"`
	Action      string `json:"action"`
}

func (s *Client) CreateIngestBudget(budget IngestBudget) (string, error) {
	body, err := s.Post("v1/ingestBudgets", budget, false)

	if err != nil {
		return "", err
	}

	var response IngestBudget

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

func (s *Client) GetIngestBudget(id string) (*IngestBudget, error) {
	body, _, err := s.Get(fmt.Sprintf("v1/ingestBudgets/%s", id), false)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	var response IngestBudget

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Client) UpdateIngestBudget(budget IngestBudget) error {
	urlPath := fmt.Sprintf("v1/ingestBudgets/%s", budget.ID)
	_, err := s.Put(urlPath, budget, false)

	return err
}

func (s *Client) DeleteIngestBudget(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/ingestBudgets/%s", id))

	return err
}

func (s *Client) FindIngestBudget(name string) (*IngestBudget, error) {
	type IngestBudgetList struct {
		Next string         `json:"next"`
		Data []IngestBudget `json:"data"`
	}

	next := ""

	for {
		body, _, err := s.Get(fmt.Sprintf("v1/ingestBudgets?next=%s", next), false)
		if err != nil {
			return nil, err
		}

		if body == nil {
			return nil, nil
		}

		var response IngestBudgetList

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		for _, budget := range response.Data {
			if budget.Name == name {
				return &budget, nil
			}
		}

		if response.Next == "" {
			break
		}

		next = response.Next
	}

	return nil, fmt.Errorf("unable to find ingest budget '%s'", name)
}

func (s *Client) CollectorAssignedToIngestBudget(ingestBudgetId string, collectorId int) (bool, error) {
	type Response struct {
		Next string `json:"next"`
		Data []struct {
			CollectorId string `json:"id"`
		} `json:"data"`
	}
	next := ""

	for {
		body, _, err := s.Get(fmt.Sprintf("v1/ingestBudgets/%s/collectors?next=%s", ingestBudgetId, next), false)
		if err != nil {
			return false, err
		}

		if body == nil {
			return false, nil
		}

		var response Response

		err = json.Unmarshal(body, &response)
		if err != nil {
			return false, err
		}

		for _, collector := range response.Data {
			if collector.CollectorId == strconv.Itoa(collectorId) {
				return true, nil
			}
		}

		if response.Next == "" {
			break
		}

		next = response.Next
	}

	return false, fmt.Errorf("unable to find collector %d for ingest budget '%s'", collectorId, ingestBudgetId)
}

func (s *Client) AssignCollectorToIngestBudget(ingestBudgetId string, collectorId int) error {
	urlPath := fmt.Sprintf("v1/ingestBudgets/%s/collectors/%d", ingestBudgetId, collectorId)
	_, err := s.Put(urlPath, nil, false)

	return err
}

func (s *Client) UnAssignCollectorToIngestBudget(ingestBudgetId string, collectorId int) error {
	_, err := s.Delete(fmt.Sprintf("v1/ingestBudgets/%s/collectors/%d", ingestBudgetId, collectorId))

	return err
}
