package sumologic

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// These tests run the resource and data source against a fake of the usage attribution config API, so they need no
// credentials. The acceptance tests call the real API.

type fakeUsageAttributionApi struct {
	mu       sync.Mutex
	puts     []map[string]interface{}
	response string
}

func newFakeUsageAttributionApi(response string) *fakeUsageAttributionApi {
	return &fakeUsageAttributionApi{response: response}
}

func (f *fakeUsageAttributionApi) lastPut(t *testing.T) map[string]interface{} {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.puts) == 0 {
		t.Fatalf("expected a PUT to the config API, got none")
	}
	return f.puts[len(f.puts)-1]
}

func (f *fakeUsageAttributionApi) putCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.puts)
}

func (f *fakeUsageAttributionApi) client(t *testing.T) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/account/usages/granular/config" {
			http.Error(w, "unexpected path "+r.URL.Path, http.StatusNotFound)
			return
		}
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			var parsed map[string]interface{}
			if err := json.Unmarshal(body, &parsed); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			f.mu.Lock()
			f.puts = append(f.puts, parsed)
			f.mu.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(f.response))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient("id", "key", "", "us2", server.URL, false)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	return client
}

const usageAttributionConfigResponse = `{
  "enabled": true,
  "dashboardId": "dash123",
  "templateVersion": 1,
  "standardTags": [{"name": "sourceCategory", "enabled": true}, {"name": "collector", "enabled": true}],
  "customTags": [{"name": "department", "enabled": true}, {"name": "cost_center", "enabled": true}],
  "lookupColumns": [],
  "breakdownByOptions": ["sourceCategory", "collector", "department", "cost_center"],
  "enabledAt": "2026-11-05T11:03:11.896846427Z"
}`

func usageAttributionResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceSumologicUsageAttributionConfig().Schema, raw)
}

func TestUsageAttributionConfigCreate_sendsEnabledAndCustomTags(t *testing.T) {
	api := newFakeUsageAttributionApi(usageAttributionConfigResponse)
	d := usageAttributionResourceData(t, map[string]interface{}{
		"custom_tags": []interface{}{"department", "cost_center"},
	})

	if err := resourceSumologicUsageAttributionConfigCreate(d, api.client(t)); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	put := api.lastPut(t)
	if put["enabled"] != true {
		t.Errorf("expected enabled to default to true, got %v", put["enabled"])
	}
	tags, _ := put["customTags"].([]interface{})
	if len(tags) != 2 {
		t.Fatalf("expected 2 custom tags, got %v", put["customTags"])
	}
	first, _ := tags[0].(map[string]interface{})
	if first["name"] != "department" || first["enabled"] != true {
		t.Errorf("expected {department, enabled}, got %v", first)
	}
	if d.Id() != "usageAttributionConfig" {
		t.Errorf("expected the fixed id, got %q", d.Id())
	}
}

func TestUsageAttributionConfigCreate_setsComputedAttributesFromResponse(t *testing.T) {
	api := newFakeUsageAttributionApi(usageAttributionConfigResponse)
	d := usageAttributionResourceData(t, map[string]interface{}{})

	if err := resourceSumologicUsageAttributionConfigCreate(d, api.client(t)); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if d.Get("dashboard_id") != "dash123" {
		t.Errorf("dashboard_id: got %v", d.Get("dashboard_id"))
	}
	if d.Get("enabled_at") != "2026-11-05T11:03:11.896846427Z" {
		t.Errorf("enabled_at: got %v", d.Get("enabled_at"))
	}
	standard := d.Get("standard_tags").([]interface{})
	if len(standard) != 2 || standard[0] != "sourceCategory" {
		t.Errorf("standard_tags: got %v", standard)
	}
	options := d.Get("breakdown_by_options").([]interface{})
	if len(options) != 4 {
		t.Errorf("breakdown_by_options: got %v", options)
	}
}

func TestUsageAttributionConfigCreate_emptyCustomTagsAreSentAsEmptyList(t *testing.T) {
	api := newFakeUsageAttributionApi(`{"enabled": true, "standardTags": [], "customTags": [], "lookupColumns": [], "breakdownByOptions": []}`)
	d := usageAttributionResourceData(t, map[string]interface{}{})

	if err := resourceSumologicUsageAttributionConfigCreate(d, api.client(t)); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// null would be read by the API as "not given"; an empty list is how a user clears the custom tags
	tags, ok := api.lastPut(t)["customTags"].([]interface{})
	if !ok || len(tags) != 0 {
		t.Errorf("expected customTags to be an empty list, got %#v", api.lastPut(t)["customTags"])
	}
}

func TestUsageAttributionConfigRead_setsStateFromApi(t *testing.T) {
	api := newFakeUsageAttributionApi(usageAttributionConfigResponse)
	d := usageAttributionResourceData(t, map[string]interface{}{})

	if err := resourceSumologicUsageAttributionConfigRead(d, api.client(t)); err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if d.Get("enabled") != true {
		t.Errorf("enabled: got %v", d.Get("enabled"))
	}
	custom := d.Get("custom_tags").([]interface{})
	if len(custom) != 2 || custom[0] != "department" || custom[1] != "cost_center" {
		t.Errorf("custom_tags keep the API order, got %v", custom)
	}
	if api.putCount() != 0 {
		t.Errorf("read must not write, got %d PUTs", api.putCount())
	}
}

func TestUsageAttributionConfigDelete_turnsItOffAndResendsTags(t *testing.T) {
	api := newFakeUsageAttributionApi(`{"enabled": false, "standardTags": [], "customTags": [], "lookupColumns": [], "breakdownByOptions": []}`)
	d := usageAttributionResourceData(t, map[string]interface{}{
		"enabled":     true,
		"custom_tags": []interface{}{"department"},
	})
	d.SetId("usageAttributionConfig")

	if err := resourceSumologicUsageAttributionConfigDelete(d, api.client(t)); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	put := api.lastPut(t)
	if put["enabled"] != false {
		t.Errorf("destroy must disable usage attribution, got enabled=%v", put["enabled"])
	}
	// destroy resends the tags the resource holds, so that the request is the same shape as an update
	tags, _ := put["customTags"].([]interface{})
	if len(tags) != 1 {
		t.Errorf("expected destroy to resend the custom tags, sent %v", put["customTags"])
	}
}

func TestUsageAttributionConfigSchema_rejectsInvalidTagNames(t *testing.T) {
	validate := resourceSumologicUsageAttributionConfig().Schema["custom_tags"].Elem.(*schema.Schema).ValidateFunc

	for _, name := range []string{"department", "_internal", "cost_center", "Tag9"} {
		if _, errs := validate(name, "custom_tags"); len(errs) != 0 {
			t.Errorf("expected %q to be valid, got %v", name, errs)
		}
	}
	for _, name := range []string{"", "9lives", "cost center", "a-b", "tag.name"} {
		_, errs := validate(name, "custom_tags")
		if len(errs) == 0 {
			t.Errorf("expected %q to be rejected", name)
		} else if !strings.Contains(errs[0].Error(), "letter or underscore") {
			t.Errorf("expected a helpful message for %q, got %v", name, errs[0])
		}
	}
}

func TestUsageAttributionConfigDataSource_readsWithoutWriting(t *testing.T) {
	api := newFakeUsageAttributionApi(usageAttributionConfigResponse)
	d := schema.TestResourceDataRaw(t, dataSourceSumologicUsageAttributionConfig().Schema, map[string]interface{}{})

	if err := dataSourceSumologicUsageAttributionConfigRead(d, api.client(t)); err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if d.Id() != "usageAttributionConfig" || d.Get("dashboard_id") != "dash123" {
		t.Errorf("unexpected state: id=%q dashboard_id=%v", d.Id(), d.Get("dashboard_id"))
	}
	if api.putCount() != 0 {
		t.Errorf("a data source must not write, got %d PUTs", api.putCount())
	}
}
