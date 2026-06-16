package sumologic

import (
	"encoding/json"
	"fmt"
)

type DataMaskRule struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name"`
	Pattern           string   `json:"pattern"`
	PiiType           string   `json:"piiType"`
	Replacement       string   `json:"replacement"`
	Scope             string   `json:"scope"`
	ScopeTargetOrgIds []string `json:"scopeTargetOrgIds,omitempty"`
	Enabled           bool     `json:"enabled"`
	Description       string   `json:"description,omitempty"`
	IsActive          bool     `json:"isActive,omitempty"`
}

type ListDataMaskRuleResponse struct {
	Data []DataMaskRule `json:"data"`
	Next string         `json:"next"`
}

func (s *ListDataMaskRuleResponse) Reset() {
	s.Data = nil
	s.Next = ""
}

func (c *Client) GetDataMaskRule(id string) (*DataMaskRule, error) {
	data, err := c.Get(fmt.Sprintf("v1/dataMaskRules/%s", id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var rule DataMaskRule
	if err = json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

func (c *Client) CreateDataMaskRule(rule DataMaskRule) (*DataMaskRule, error) {
	responseBody, err := c.Post("v1/dataMaskRules", rule)
	if err != nil {
		return nil, err
	}

	var createdRule DataMaskRule
	if err = json.Unmarshal(responseBody, &createdRule); err != nil {
		return nil, err
	}

	return &createdRule, nil
}

func (c *Client) UpdateDataMaskRule(rule DataMaskRule) (*DataMaskRule, error) {
	url := fmt.Sprintf("v1/dataMaskRules/%s", rule.ID)
	responseBody, err := c.Put(url, rule)
	if err != nil {
		return nil, err
	}

	var updatedRule DataMaskRule
	if err = json.Unmarshal(responseBody, &updatedRule); err != nil {
		return nil, err
	}

	return &updatedRule, nil
}

func (c *Client) DeleteDataMaskRule(id string) error {
	_, err := c.Delete(fmt.Sprintf("v1/dataMaskRules/%s", id))
	return err
}

func (c *Client) ListDataMaskRules() ([]DataMaskRule, error) {
	var listResponse ListDataMaskRuleResponse

	data, err := c.Get("v1/dataMaskRules")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &listResponse); err != nil {
		return nil, err
	}

	rules := listResponse.Data

	for listResponse.Next != "" {
		data, err = c.Get("v1/dataMaskRules?token=" + listResponse.Next)
		if err != nil {
			return nil, err
		}

		listResponse.Reset()

		if err = json.Unmarshal(data, &listResponse); err != nil {
			return nil, err
		}

		rules = append(rules, listResponse.Data...)
	}

	return rules, nil
}


