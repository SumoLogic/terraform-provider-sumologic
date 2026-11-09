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

	policies := &Policies{
		Audit:                              auditPolicy,
		DataAccessLevel:                    dataAccessLevelPolicy,
		MaxUserSessionTimeout:              maxUserSessionTimeoutPolicy,
		SearchAudit:                        searchAuditPolicy,
		ShareDashboardsOutsideOrganization: shareDashboardsOutsideOrganizationPolicy,
		UserConcurrentSessionsLimit:        userConcurrentSessionsLimitPolicy,
	}
	if err := s.getOptionalPolicies(policies); err != nil {
		return nil, err
	}
	return policies, nil
}

// getOptionalPolicies reads the policies that older deployments may not serve. A policy whose endpoint is
// missing (404) is left nil instead of failing the whole read.
func (s *Client) getOptionalPolicies(policies *Policies) error {
	var oAuthCimdPolicy OAuthCimdPolicy
	var disableUnusedAccessKeysPolicy DisableUnusedAccessKeysPolicy
	var accessKeysLifetimePolicy AccessKeysLifetimePolicy

	optionalPolicies := []struct {
		endpoint string
		response interface{}
		assign   func()
	}{
		{"oAuthCimd", &oAuthCimdPolicy, func() { policies.OAuthCimd = &oAuthCimdPolicy }},
		{"disableUnusedAccessKeys", &disableUnusedAccessKeysPolicy, func() { policies.DisableUnusedAccessKeys = &disableUnusedAccessKeysPolicy }},
		{"accessKeysLifetime", &accessKeysLifetimePolicy, func() { policies.AccessKeysLifetime = &accessKeysLifetimePolicy }},
	}
	for _, optional := range optionalPolicies {
		found, err := s.getPolicyIfPresent(optional.endpoint, optional.response)
		if err != nil {
			return err
		}
		if found {
			optional.assign()
		}
	}
	return nil
}

func (s *Client) UpdatePolicies(policies Policies) (*Policies, error) {
	var auditPolicy AuditPolicy
	var dataAccessLevelPolicy DataAccessLevelPolicy
	var maxUserSessionTimeoutPolicy MaxUserSessionTimeoutPolicy
	var searchAuditPolicy SearchAuditPolicy
	var shareDashboardsOutsideOrganizationPolicy ShareDashboardsOutsideOrganizationPolicy
	var userConcurrentSessionsLimitPolicy UserConcurrentSessionsLimitPolicy
	var oAuthCimdPolicy OAuthCimdPolicy
	var disableUnusedAccessKeysPolicy DisableUnusedAccessKeysPolicy
	var accessKeysLifetimePolicy AccessKeysLifetimePolicy

	policyRequests := []PolicyRequest{
		PolicyRequest{"audit", &policies.Audit, &auditPolicy},
		PolicyRequest{"dataAccessLevel", &policies.DataAccessLevel, &dataAccessLevelPolicy},
		PolicyRequest{"maxUserSessionTimeout", &policies.MaxUserSessionTimeout, &maxUserSessionTimeoutPolicy},
		PolicyRequest{"searchAudit", &policies.SearchAudit, &searchAuditPolicy},
		PolicyRequest{"shareDashboardsOutsideOrganization", &policies.ShareDashboardsOutsideOrganization, &shareDashboardsOutsideOrganizationPolicy},
		PolicyRequest{"userConcurrentSessionsLimit", &policies.UserConcurrentSessionsLimit, &userConcurrentSessionsLimitPolicy},
	}
	// Optional policies are only sent when set (see resourceToPolicies).
	if policies.OAuthCimd != nil {
		policyRequests = append(policyRequests, PolicyRequest{"oAuthCimd", policies.OAuthCimd, &oAuthCimdPolicy})
	}
	if policies.DisableUnusedAccessKeys != nil {
		policyRequests = append(policyRequests, PolicyRequest{"disableUnusedAccessKeys", policies.DisableUnusedAccessKeys, &disableUnusedAccessKeysPolicy})
	}
	if policies.AccessKeysLifetime != nil {
		policyRequests = append(policyRequests, PolicyRequest{"accessKeysLifetime", policies.AccessKeysLifetime, &accessKeysLifetimePolicy})
	}

	var err error
	for _, policyRequest := range policyRequests {
		err = s.putPolicy(policyRequest.Endpoint, policyRequest.Body, policyRequest.Response)
		if err != nil {
			return nil, err
		}
	}

	updated := &Policies{
		Audit:                              auditPolicy,
		DataAccessLevel:                    dataAccessLevelPolicy,
		MaxUserSessionTimeout:              maxUserSessionTimeoutPolicy,
		SearchAudit:                        searchAuditPolicy,
		ShareDashboardsOutsideOrganization: shareDashboardsOutsideOrganizationPolicy,
		UserConcurrentSessionsLimit:        userConcurrentSessionsLimitPolicy,
	}
	// Read back every optional policy (including the ones just written) so the returned state is complete.
	if err := s.getOptionalPolicies(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Client) getPolicy(endpoint string, policy interface{}) error {
	url := fmt.Sprintf("v1/policies/%s", endpoint)

	data, err := s.Get(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &policy)
}

// getPolicyIfPresent is like getPolicy but reports found=false when the endpoint does not exist (404).
func (s *Client) getPolicyIfPresent(endpoint string, policy interface{}) (bool, error) {
	data, err := s.Get(fmt.Sprintf("v1/policies/%s", endpoint))
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil
	}
	return true, json.Unmarshal(data, policy)
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

type OAuthCimdPolicy struct {
	OAuthCimdPolicy string `json:"oAuthCimdPolicy"`
}

type DisableUnusedAccessKeysPolicy struct {
	UnusedAccessKeysDisableAfterInDays int `json:"unusedAccessKeysDisableAfterInDays"`
}

type AccessKeysLifetimePolicy struct {
	AccessKeysLifetimeInDays string `json:"accessKeysLifetimeInDays"`
}

type Policies struct {
	Audit                              AuditPolicy
	DataAccessLevel                    DataAccessLevelPolicy
	MaxUserSessionTimeout              MaxUserSessionTimeoutPolicy
	SearchAudit                        SearchAuditPolicy
	ShareDashboardsOutsideOrganization ShareDashboardsOutsideOrganizationPolicy
	UserConcurrentSessionsLimit        UserConcurrentSessionsLimitPolicy
	// The policies below are optional: nil means "leave the server-side value untouched".
	OAuthCimd               *OAuthCimdPolicy
	DisableUnusedAccessKeys *DisableUnusedAccessKeysPolicy
	AccessKeysLifetime      *AccessKeysLifetimePolicy
}
