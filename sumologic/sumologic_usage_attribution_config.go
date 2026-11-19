package sumologic

import (
	"encoding/json"
)

const usageAttributionConfigUrl = "v1/account/usages/granular/config"

// UsageAttributionTag is a tag dimension that the usage attribution dashboard can break volume down by.
type UsageAttributionTag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// UsageAttributionConfig is the usage attribution (granular usage visibility) configuration of an organization.
type UsageAttributionConfig struct {
	Enabled            bool                  `json:"enabled"`
	DashboardId        string                `json:"dashboardId,omitempty"`
	TemplateVersion    int                   `json:"templateVersion"`
	StandardTags       []UsageAttributionTag `json:"standardTags"`
	CustomTags         []UsageAttributionTag `json:"customTags"`
	LookupTableId      string                `json:"lookupTableId,omitempty"`
	LookupKeyField     string                `json:"lookupKeyField,omitempty"`
	LookupColumns      []string              `json:"lookupColumns"`
	BreakdownByOptions []string              `json:"breakdownByOptions"`
	EnabledAt          string                `json:"enabledAt,omitempty"`
}

// UsageAttributionConfigUpdate is the request body of the update call. Only custom tags and the enabled flag can be
// set; every other field of UsageAttributionConfig is derived by Sumo Logic.
type UsageAttributionConfigUpdate struct {
	Enabled    bool                  `json:"enabled"`
	CustomTags []UsageAttributionTag `json:"customTags"`
}

func (s *Client) GetUsageAttributionConfig() (*UsageAttributionConfig, error) {
	data, err := s.Get(usageAttributionConfigUrl)
	if err != nil {
		return nil, err
	}

	var config UsageAttributionConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (s *Client) UpdateUsageAttributionConfig(update UsageAttributionConfigUpdate) (*UsageAttributionConfig, error) {
	data, err := s.Put(usageAttributionConfigUrl, update)
	if err != nil {
		return nil, err
	}

	var config UsageAttributionConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
