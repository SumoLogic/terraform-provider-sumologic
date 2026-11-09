package sumologic

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These tests guard backward compatibility of sumologic_policies for configurations written before the optional
// oauth_cimd_policy, unused_access_keys_disable_after_in_days and access_keys_lifetime_in_days were added.

func policiesResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceSumologicPolicies().Schema, raw)
}

func legacyPoliciesConfig() map[string]interface{} {
	return map[string]interface{}{
		"audit":                                 true,
		"data_access_level":                     true,
		"max_user_session_timeout":              "1h",
		"search_audit":                          true,
		"share_dashboards_outside_organization": true,
		"user_concurrent_sessions_limit": []interface{}{
			map[string]interface{}{"enabled": true, "max_concurrent_sessions": 70},
		},
	}
}

func TestPoliciesResourceToPolicies_legacyConfigLeavesOptionalPoliciesUnset(t *testing.T) {
	policies := resourceToPolicies(policiesResourceData(t, legacyPoliciesConfig()))

	if policies.OAuthCimd != nil || policies.DisableUnusedAccessKeys != nil || policies.AccessKeysLifetime != nil {
		t.Fatalf("expected optional policies to be nil for a legacy config, got %+v", policies)
	}
}

func TestPoliciesResourceToPolicies_explicitZeroIsSet(t *testing.T) {
	config := legacyPoliciesConfig()
	config["unused_access_keys_disable_after_in_days"] = 0
	config["access_keys_lifetime_in_days"] = "0"
	config["oauth_cimd_policy"] = "disabled"

	policies := resourceToPolicies(policiesResourceData(t, config))

	if policies.DisableUnusedAccessKeys == nil || policies.DisableUnusedAccessKeys.UnusedAccessKeysDisableAfterInDays != 0 {
		t.Errorf("expected unused access keys policy set to 0, got %+v", policies.DisableUnusedAccessKeys)
	}
	if policies.AccessKeysLifetime == nil || policies.AccessKeysLifetime.AccessKeysLifetimeInDays != "0" {
		t.Errorf("expected access keys lifetime policy set to \"0\", got %+v", policies.AccessKeysLifetime)
	}
	if policies.OAuthCimd == nil || policies.OAuthCimd.OAuthCimdPolicy != "disabled" {
		t.Errorf("expected oauth cimd policy set to disabled, got %+v", policies.OAuthCimd)
	}
}

func TestPoliciesSchema_rejectsInvalidValues(t *testing.T) {
	cases := map[string]interface{}{
		"oauth_cimd_policy":                        "sometimes",
		"unused_access_keys_disable_after_in_days": 366,
		"access_keys_lifetime_in_days":             "42",
	}
	for key, value := range cases {
		config := legacyPoliciesConfig()
		config[key] = value
		diags := resourceSumologicPolicies().Validate(terraform.NewResourceConfigRaw(config))
		if !diags.HasError() {
			t.Errorf("expected validation error for %s=%v", key, value)
		}
	}
}

// fakePoliciesServer records every request and serves the policy endpoints with fixed server-side values.
type fakePoliciesServer struct {
	mu       sync.Mutex
	requests []string
	bodies   map[string]string
	server   *httptest.Server
}

func newFakePoliciesServer(t *testing.T) *fakePoliciesServer {
	t.Helper()
	f := &fakePoliciesServer{bodies: map[string]string{}}
	responses := map[string]string{
		"audit":                              `{"enabled":true}`,
		"dataAccessLevel":                    `{"enabled":true}`,
		"maxUserSessionTimeout":              `{"maxUserSessionTimeout":"1h"}`,
		"searchAudit":                        `{"enabled":true}`,
		"shareDashboardsOutsideOrganization": `{"enabled":true}`,
		"userConcurrentSessionsLimit":        `{"enabled":true,"maxConcurrentSessions":70}`,
		"oAuthCimd":                          `{"oAuthCimdPolicy":"enabled"}`,
		"disableUnusedAccessKeys":            `{"unusedAccessKeysDisableAfterInDays":30}`,
		"accessKeysLifetime":                 `{"accessKeysLifetimeInDays":"60"}`,
	}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.requests = append(f.requests, r.Method+" "+endpoint)
		f.bodies[r.Method+" "+endpoint] = string(body)
		f.mu.Unlock()
		resp, ok := responses[endpoint]
		if !ok {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodPut {
			// Echo what was written, like the real API.
			resp = string(body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(resp))
	}))
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakePoliciesServer) client(t *testing.T) *Client {
	t.Helper()
	c, err := NewClient("id", "key", "", "", f.server.URL+"/api/", false)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func (f *fakePoliciesServer) saw(request string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, r := range f.requests {
		if r == request {
			return true
		}
	}
	return false
}

var optionalPolicyEndpoints = []string{"oAuthCimd", "disableUnusedAccessKeys", "accessKeysLifetime"}

func TestUpdatePolicies_unsetOptionalPoliciesAreNotWrittenButAreReadBack(t *testing.T) {
	f := newFakePoliciesServer(t)
	policies := resourceToPolicies(policiesResourceData(t, legacyPoliciesConfig()))

	updated, err := f.client(t).UpdatePolicies(policies)
	if err != nil {
		t.Fatal(err)
	}

	for _, endpoint := range optionalPolicyEndpoints {
		if f.saw("PUT " + endpoint) {
			t.Errorf("unset optional policy %s must not be written", endpoint)
		}
		if !f.saw("GET " + endpoint) {
			t.Errorf("unset optional policy %s should be read back", endpoint)
		}
	}
	if !f.saw("PUT audit") || !f.saw("PUT userConcurrentSessionsLimit") {
		t.Errorf("required policies must still be written, saw %v", f.requests)
	}
	if updated.OAuthCimd == nil || updated.OAuthCimd.OAuthCimdPolicy != "enabled" ||
		updated.DisableUnusedAccessKeys == nil || updated.DisableUnusedAccessKeys.UnusedAccessKeysDisableAfterInDays != 30 ||
		updated.AccessKeysLifetime == nil || updated.AccessKeysLifetime.AccessKeysLifetimeInDays != "60" {
		t.Errorf("server values were not read into the returned policies: %+v", updated)
	}
}

func TestUpdatePolicies_setOptionalPoliciesAreWritten(t *testing.T) {
	f := newFakePoliciesServer(t)
	config := legacyPoliciesConfig()
	config["oauth_cimd_policy"] = "enabled-pre-registered-only"
	config["unused_access_keys_disable_after_in_days"] = 0
	config["access_keys_lifetime_in_days"] = "180"

	_, err := f.client(t).UpdatePolicies(resourceToPolicies(policiesResourceData(t, config)))
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"oAuthCimd":               `{"oAuthCimdPolicy":"enabled-pre-registered-only"}`,
		"disableUnusedAccessKeys": `{"unusedAccessKeysDisableAfterInDays":0}`,
		"accessKeysLifetime":      `{"accessKeysLifetimeInDays":"180"}`,
	}
	for endpoint, want := range expected {
		got, ok := f.bodies["PUT "+endpoint]
		if !ok {
			t.Errorf("expected PUT %s", endpoint)
			continue
		}
		if !jsonEqual(t, got, want) {
			t.Errorf("PUT %s body = %s, want %s", endpoint, got, want)
		}
	}
}

func TestDeletePolicies_doesNotTouchOptionalPolicies(t *testing.T) {
	f := newFakePoliciesServer(t)
	d := policiesResourceData(t, legacyPoliciesConfig())

	if err := resourceSumologicPoliciesDelete(d, f.client(t)); err != nil {
		t.Fatal(err)
	}

	for _, endpoint := range optionalPolicyEndpoints {
		if f.saw("PUT " + endpoint) {
			t.Errorf("destroy must not reset %s", endpoint)
		}
	}
}

func jsonEqual(t *testing.T, a, b string) bool {
	t.Helper()
	var x, y interface{}
	if err := json.Unmarshal([]byte(a), &x); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &y); err != nil {
		return false
	}
	xb, _ := json.Marshal(x)
	yb, _ := json.Marshal(y)
	return string(xb) == string(yb)
}

// priorPoliciesState mimics state after a refresh: the server holds values for the optional policies.
func priorPoliciesState() *terraform.InstanceState {
	return &terraform.InstanceState{ID: "org-policies", Attributes: map[string]string{
		"audit":                                    "true",
		"data_access_level":                        "true",
		"max_user_session_timeout":                 "1h",
		"search_audit":                             "true",
		"share_dashboards_outside_organization":    "true",
		"user_concurrent_sessions_limit.#":         "1",
		"user_concurrent_sessions_limit.0.enabled": "true",
		"user_concurrent_sessions_limit.0.max_concurrent_sessions": "70",
		"oauth_cimd_policy":                        "enabled",
		"unused_access_keys_disable_after_in_days": "30",
		"access_keys_lifetime_in_days":             "60",
	}}
}

func policiesOnUpdate(t *testing.T, config map[string]interface{}) Policies {
	t.Helper()
	r := resourceSumologicPolicies()
	state := priorPoliciesState()
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatal(err)
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatal(err)
	}
	return resourceToPolicies(d)
}

// Writing an access keys policy can have side effects even when the value is unchanged, so an unrelated update
// must never re-send the optional policies.
func TestPoliciesResourceToPolicies_unrelatedUpdateDoesNotResendOptionalPolicies(t *testing.T) {
	config := legacyPoliciesConfig()
	config["audit"] = false

	policies := policiesOnUpdate(t, config)

	if policies.OAuthCimd != nil || policies.DisableUnusedAccessKeys != nil || policies.AccessKeysLifetime != nil {
		t.Fatalf("optional policies must not be re-sent when unchanged, got cimd=%v unused=%v lifetime=%v",
			policies.OAuthCimd, policies.DisableUnusedAccessKeys, policies.AccessKeysLifetime)
	}
}

func TestPoliciesResourceToPolicies_updateSendsOnlyChangedOptionalPolicies(t *testing.T) {
	config := legacyPoliciesConfig()
	config["access_keys_lifetime_in_days"] = "90"

	policies := policiesOnUpdate(t, config)

	if policies.AccessKeysLifetime == nil || policies.AccessKeysLifetime.AccessKeysLifetimeInDays != "90" {
		t.Errorf("changed lifetime should be sent, got %v", policies.AccessKeysLifetime)
	}
	if policies.OAuthCimd != nil || policies.DisableUnusedAccessKeys != nil {
		t.Errorf("unchanged policies must not be sent, got cimd=%v unused=%v", policies.OAuthCimd, policies.DisableUnusedAccessKeys)
	}
}

func TestPoliciesResourceToPolicies_updateToExplicitZeroIsSent(t *testing.T) {
	config := legacyPoliciesConfig()
	config["unused_access_keys_disable_after_in_days"] = 0

	policies := policiesOnUpdate(t, config)

	if policies.DisableUnusedAccessKeys == nil || policies.DisableUnusedAccessKeys.UnusedAccessKeysDisableAfterInDays != 0 {
		t.Errorf("changing 30 -> 0 must be sent, got %v", policies.DisableUnusedAccessKeys)
	}
}

// Older deployments may not serve the newer policy endpoints; reading must still work for existing users.
func TestGetPolicies_toleratesMissingOptionalEndpoints(t *testing.T) {
	f := newFakePoliciesServer(t)
	inner := f.server.Config.Handler
	f.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, endpoint := range optionalPolicyEndpoints {
			if strings.HasSuffix(r.URL.Path, "/"+endpoint) {
				http.NotFound(w, r)
				return
			}
		}
		inner.ServeHTTP(w, r)
	})

	policies, err := f.client(t).GetPolicies()
	if err != nil {
		t.Fatalf("GetPolicies must tolerate missing optional endpoints: %v", err)
	}
	if policies.OAuthCimd != nil || policies.DisableUnusedAccessKeys != nil || policies.AccessKeysLifetime != nil {
		t.Errorf("missing endpoints should yield nil policies, got %+v", policies)
	}
	if !policies.Audit.Enabled {
		t.Errorf("required policies should still be read")
	}
}
