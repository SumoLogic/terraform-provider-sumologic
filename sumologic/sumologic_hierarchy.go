package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateHierarchy(hierarchy Hierarchy) (string, error) {
	urlWithoutParams := "v1/entities/hierarchies"

	data, err := s.Post(urlWithoutParams, hierarchy, false)
	if err != nil {
		return "", err
	}

	var createdHierarchy Hierarchy

	err = json.Unmarshal(data, &createdHierarchy)
	if err != nil {
		return "", err
	}

	return createdHierarchy.ID, nil

}

func (s *Client) GetHierarchy(id string) (*Hierarchy, error) {
	urlWithoutParams := "v1/entities/hierarchies/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams, false)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var hierarchy Hierarchy

	err = json.Unmarshal(data, &hierarchy)
	if err != nil {
		return nil, err
	}

	return &hierarchy, nil

}

func (s *Client) DeleteHierarchy(id string) error {
	urlWithoutParams := "v1/entities/hierarchies/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateHierarchy(hierarchy Hierarchy) error {
	urlWithoutParams := "v1/entities/hierarchies/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, hierarchy.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	hierarchy.ID = ""

	_, err := s.Put(urlWithParams, hierarchy, false)

	return err

}

type Hierarchy struct {
	Name   string                    `json:"name"`
	ID     string                    `json:"id,omitempty"`
	Filter *HierarchyFilteringClause `json:"filter"`
	Level  Level                     `json:"level"`
}

type HierarchyFilteringClause struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Level struct {
	EntityType               string               `json:"entityType"`
	NextLevelsWithConditions []LevelWithCondition `json:"nextLevelsWithConditions"`
	NextLevel                *Level               `json:"nextLevel,omitempty"`
}

type LevelWithCondition struct {
	Condition string `json:"condition"`
	Level     Level  `json:"level"`
}
