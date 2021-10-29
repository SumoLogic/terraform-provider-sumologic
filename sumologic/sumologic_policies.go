package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetPolicies() (*Policies, error) {
	var auditPolicy AuditPolicy
	var dataAccessLevelPolicy DataAccessLevelPolicy
	var maxUserSessionTimeoutPolicy MaxUserSessionTimeoutPolicy
	var searchAuditPolicy SearchAuditPolicy
	var shareDashboardsOutsideOrganizationPolicy ShareDashboardsOutsideOrganizationPolicy
	var userConcurrentSessionsLimitPolicy UserConcurrentSessionsLimitPolicy

	policyRequests := []PolicyRequest{
		PolicyRequest{"audit", nil, &auditPolicy},
		PolicyRequest{"dataAccessLevel", nil, &dataAccessLevelPolicy},
		PolicyRequest{"maxUserSessionTimeout", nil, &maxUserSessionTimeoutPolicy},
		PolicyRequest{"searchAudit", nil, &searchAuditPolicy},
		PolicyRequest{"shareDashboardsOutsideOrganization", nil, &shareDashboardsOutsideOrganizationPolicy},
		PolicyRequest{"userConcurrentSessionsLimit", nil, &userConcurrentSessionsLimitPolicy},
	}

	var err error
	for _, policyRequest := range policyRequests {
		err = s.getPolicy(policyRequest.Endpoint, policyRequest.Response)
		if err != nil {
			return nil, err
		}
	}

	return &Policies{
		Audit:                              auditPolicy,
		DataAccessLevel:                    dataAccessLevelPolicy,
		MaxUserSessionTimeout:              maxUserSessionTimeoutPolicy,
		SearchAudit:                        searchAuditPolicy,
		ShareDashboardsOutsideOrganization: shareDashboardsOutsideOrganizationPolicy,
		UserConcurrentSessionsLimit:        userConcurrentSessionsLimitPolicy,
	}, nil
}

func (s *Client) UpdatePolicies(policies Policies) (*Policies, error) {
	var auditPolicy AuditPolicy
	var dataAccessLevelPolicy DataAccessLevelPolicy
	var maxUserSessionTimeoutPolicy MaxUserSessionTimeoutPolicy
	var searchAuditPolicy SearchAuditPolicy
	var shareDashboardsOutsideOrganizationPolicy ShareDashboardsOutsideOrganizationPolicy
	var userConcurrentSessionsLimitPolicy UserConcurrentSessionsLimitPolicy

	policyRequests := []PolicyRequest{
		PolicyRequest{"audit", &policies.Audit, &auditPolicy},
		PolicyRequest{"dataAccessLevel", &policies.DataAccessLevel, &dataAccessLevelPolicy},
		PolicyRequest{"maxUserSessionTimeout", &policies.MaxUserSessionTimeout, &maxUserSessionTimeoutPolicy},
		PolicyRequest{"searchAudit", &policies.SearchAudit, &searchAuditPolicy},
		PolicyRequest{"shareDashboardsOutsideOrganization", &policies.ShareDashboardsOutsideOrganization, &shareDashboardsOutsideOrganizationPolicy},
		PolicyRequest{"userConcurrentSessionsLimit", &policies.UserConcurrentSessionsLimit, &userConcurrentSessionsLimitPolicy},
	}

	var err error
	for _, policyRequest := range policyRequests {
		err = s.putPolicy(policyRequest.Endpoint, policyRequest.Body, policyRequest.Response)
		if err != nil {
			return nil, err
		}
	}

	return &Policies{
		Audit:                              auditPolicy,
		DataAccessLevel:                    dataAccessLevelPolicy,
		MaxUserSessionTimeout:              maxUserSessionTimeoutPolicy,
		SearchAudit:                        searchAuditPolicy,
		ShareDashboardsOutsideOrganization: shareDashboardsOutsideOrganizationPolicy,
		UserConcurrentSessionsLimit:        userConcurrentSessionsLimitPolicy,
	}, nil
}

func (s *Client) getPolicy(endpoint string, policy interface{}) error {
	url := fmt.Sprintf("v1/policies/%s", endpoint)

	data, _, err := s.Get(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &policy)
}

func (s *Client) putPolicy(endpoint string, body interface{}, policy interface{}) error {
	url := fmt.Sprintf("v1/policies/%s", endpoint)

	data, err := s.Put(url, body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &policy)
}

type PolicyRequest struct {
	Endpoint string
	Body     interface{}
	Response interface{}
}

type AuditPolicy struct {
	Enabled bool `json:"enabled"`
}

type DataAccessLevelPolicy struct {
	Enabled bool `json:"enabled"`
}

type MaxUserSessionTimeoutPolicy struct {
	MaxUserSessionTimeout string `json:"maxUserSessionTimeout"`
}

type SearchAuditPolicy struct {
	Enabled bool `json:"enabled"`
}

type ShareDashboardsOutsideOrganizationPolicy struct {
	Enabled bool `json:"enabled"`
}

type UserConcurrentSessionsLimitPolicy struct {
	Enabled               bool `json:"enabled"`
	MaxConcurrentSessions int  `json:"maxConcurrentSessions"`
}

type Policies struct {
	Audit                              AuditPolicy
	DataAccessLevel                    DataAccessLevelPolicy
	MaxUserSessionTimeout              MaxUserSessionTimeoutPolicy
	SearchAudit                        SearchAuditPolicy
	ShareDashboardsOutsideOrganization ShareDashboardsOutsideOrganizationPolicy
	UserConcurrentSessionsLimit        UserConcurrentSessionsLimitPolicy
}
