package sumologic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These tests run the real provider through Terraform against a stateful fake of the policies API, so they need
// a terraform binary (TF_ACC=1) but no Sumo Logic credentials. They check that configurations written before the
// optional access keys / OAuth policies existed are not affected by them.

// baselineProviderVersion is the last published release before the optional access keys / OAuth policies were added.
// Keep it fixed: it is the "before" side of the upgrade test, so moving it to a release that already has these
// attributes would stop the test from checking an upgrade. The sumologic_policies code is identical in 3.3.6.
const baselineProviderVersion = "3.3.5"

type statefulPoliciesAPI struct {
	mu     sync.Mutex
	values map[string]string
	writes []string
	server *httptest.Server
}

func newStatefulPoliciesAPI(t *testing.T) *statefulPoliciesAPI {
	t.Helper()
	api := &statefulPoliciesAPI{values: map[string]string{
		"audit":                              `{"enabled":false}`,
		"dataAccessLevel":                    `{"enabled":false}`,
		"maxUserSessionTimeout":              `{"maxUserSessionTimeout":"7d"}`,
		"searchAudit":                        `{"enabled":false}`,
		"shareDashboardsOutsideOrganization": `{"enabled":false}`,
		"userConcurrentSessionsLimit":        `{"enabled":false,"maxConcurrentSessions":100}`,
		// The organization already has non-default values for the newer policies, as a customer who set them in the
		// UI would have.
		"oAuthCimd":               `{"oAuthCimdPolicy":"enabled-pre-registered-only"}`,
		"disableUnusedAccessKeys": `{"unusedAccessKeysDisableAfterInDays":45}`,
		"accessKeysLifetime":      `{"accessKeysLifetimeInDays":"90"}`,
	}}
	api.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		api.mu.Lock()
		defer api.mu.Unlock()
		current, ok := api.values[endpoint]
		if !ok || !strings.Contains(r.URL.Path, "/v1/policies/") {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			api.values[endpoint] = string(body)
			api.writes = append(api.writes, endpoint)
			current = string(body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(current))
	}))
	t.Cleanup(api.server.Close)
	return api
}

func (a *statefulPoliciesAPI) value(endpoint string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.values[endpoint]
}

func (a *statefulPoliciesAPI) writtenEndpoints() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.writes...)
}

func (a *statefulPoliciesAPI) resetWrites() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.writes = nil
}

// checkWritten guards against a vacuous "nothing else was written" check: the change under test must have been sent.
func (a *statefulPoliciesAPI) checkWritten(endpoint string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		for _, written := range a.writtenEndpoints() {
			if written == endpoint {
				return nil
			}
		}
		return fmt.Errorf("expected %s to be written, writes: %v", endpoint, a.writtenEndpoints())
	}
}

// sameJSON compares two JSON documents ignoring formatting and key order.
func sameJSON(a, b string) bool {
	var x, y interface{}
	if json.Unmarshal([]byte(a), &x) != nil || json.Unmarshal([]byte(b), &y) != nil {
		return false
	}
	return reflect.DeepEqual(x, y)
}

// checkDestroyedLikeLegacyProvider verifies what destroy leaves behind: the required policies are reset to their
// defaults and the newer policies keep the organization's values.
func (a *statefulPoliciesAPI) checkDestroyedLikeLegacyProvider() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// The state passed to CheckDestroy is the one from before the destroy, so verify against the API instead.
		for endpoint, want := range map[string]string{
			"audit":                       `{"enabled":false}`,
			"maxUserSessionTimeout":       `{"maxUserSessionTimeout":"7d"}`,
			"userConcurrentSessionsLimit": `{"enabled":false,"maxConcurrentSessions":100}`,
		} {
			if got := a.value(endpoint); !sameJSON(got, want) {
				return fmt.Errorf("%s not reset on destroy: got %s, want %s", endpoint, got, want)
			}
		}
		return a.checkNewPoliciesKeepOrganizationValues()(s)
	}
}

func (a *statefulPoliciesAPI) checkNewPoliciesWereNeverWritten() resource.TestCheckFunc {
	return func(*terraform.State) error {
		for _, endpoint := range a.writtenEndpoints() {
			switch endpoint {
			case "oAuthCimd", "disableUnusedAccessKeys", "accessKeysLifetime":
				return fmt.Errorf("existing configuration must not write %s, writes: %v", endpoint, a.writtenEndpoints())
			}
		}
		return nil
	}
}

func (a *statefulPoliciesAPI) checkNewPoliciesKeepOrganizationValues() resource.TestCheckFunc {
	return func(*terraform.State) error {
		expected := map[string]string{
			"oAuthCimd":               `{"oAuthCimdPolicy":"enabled-pre-registered-only"}`,
			"disableUnusedAccessKeys": `{"unusedAccessKeysDisableAfterInDays":45}`,
			"accessKeysLifetime":      `{"accessKeysLifetimeInDays":"90"}`,
		}
		for endpoint, want := range expected {
			if got := a.value(endpoint); got != want {
				return fmt.Errorf("%s changed: got %s, want %s", endpoint, got, want)
			}
		}
		return nil
	}
}

// pointProviderAtFakeAPI makes the provider talk to the fake API through the environment variables it already reads.
func pointProviderAtFakeAPI(t *testing.T, api *statefulPoliciesAPI) {
	t.Helper()
	t.Setenv("SUMOLOGIC_BASE_URL", api.server.URL+"/api/")
	t.Setenv("SUMOLOGIC_ACCESSID", "fake")
	t.Setenv("SUMOLOGIC_ACCESSKEY", "fake")
	t.Setenv("SUMOLOGIC_ENVIRONMENT", "")
}

// legacyPoliciesTerraformConfig is a configuration written before the optional policies existed.
func legacyPoliciesTerraformConfig(audit bool) string {
	return fmt.Sprintf(`
resource "sumologic_policies" "legacy" {
  audit                                 = %t
  data_access_level                     = false
  max_user_session_timeout              = "7d"
  search_audit                          = false
  share_dashboards_outside_organization = false
  user_concurrent_sessions_limit {
    enabled                 = false
    max_concurrent_sessions = 100
  }
}`, audit)
}

// localProviderFactories serves the provider built from this branch.
func localProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"sumologic": func() (*schema.Provider, error) { return Provider(), nil },
	}
}

// A legacy configuration must not write the newer policies, must plan clean afterwards, and an unrelated change
// must still not write them.
func TestAccPoliciesCompat_legacyConfigNeverTouchesNewPolicies(t *testing.T) {
	api := newStatefulPoliciesAPI(t)
	pointProviderAtFakeAPI(t, api)

	resource.Test(t, resource.TestCase{
		ProviderFactories: localProviderFactories(),
		CheckDestroy:      api.checkDestroyedLikeLegacyProvider(),
		Steps: []resource.TestStep{
			{
				Config: legacyPoliciesTerraformConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "audit", "false"),
					// The existing organization values are read into state, not overwritten.
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "oauth_cimd_policy", "enabled-pre-registered-only"),
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "unused_access_keys_disable_after_in_days", "45"),
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "access_keys_lifetime_in_days", "90"),
					api.checkNewPoliciesWereNeverWritten(),
					api.checkNewPoliciesKeepOrganizationValues(),
				),
			},
			{
				// Same configuration again: no diff.
				Config:   legacyPoliciesTerraformConfig(false),
				PlanOnly: true,
			},
			{
				// Unrelated change: the required policies are written, the newer ones are not.
				PreConfig: api.resetWrites,
				Config:    legacyPoliciesTerraformConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "audit", "true"),
					api.checkWritten("audit"),
					api.checkNewPoliciesWereNeverWritten(),
					api.checkNewPoliciesKeepOrganizationValues(),
				),
			},
		},
	})
}

// State written by the baseline release (which does not know the newer attributes) must plan clean with this
// build and must not be rewritten by it.
func TestAccPoliciesCompat_stateFromReleasedProviderPlansCleanly(t *testing.T) {
	api := newStatefulPoliciesAPI(t)
	pointProviderAtFakeAPI(t, api)
	config := legacyPoliciesTerraformConfig(false)

	resource.Test(t, resource.TestCase{
		CheckDestroy: api.checkDestroyedLikeLegacyProvider(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"sumologic": {Source: "SumoLogic/sumologic", VersionConstraint: baselineProviderVersion},
				},
				Config: config,
			},
			{
				// Same configuration with the provider built from this branch: nothing to change.
				ProviderFactories: localProviderFactories(),
				Config:            config,
				PlanOnly:          true,
			},
			{
				// A real change after the upgrade runs an update, which must not write the newer policies.
				ProviderFactories: localProviderFactories(),
				Config:            legacyPoliciesTerraformConfig(true),
				PreConfig:         api.resetWrites,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_policies.legacy", "audit", "true"),
					api.checkWritten("audit"),
					api.checkNewPoliciesWereNeverWritten(),
					api.checkNewPoliciesKeepOrganizationValues(),
				),
			},
		},
	})
}
