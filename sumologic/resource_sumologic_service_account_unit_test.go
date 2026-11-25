package sumologic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// These tests run the service account resource and data source against a stateful fake of the service account API,
// so they need no credentials. The acceptance tests call the real API.

type fakeServiceAccountApi struct {
	mu       sync.Mutex
	accounts map[string]map[string]interface{}
	calls    []string
	nextID   int
}

func newFakeServiceAccountApi() *fakeServiceAccountApi {
	return &fakeServiceAccountApi{accounts: map[string]map[string]interface{}{}}
}

func (f *fakeServiceAccountApi) callLog() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.calls...)
}

func (f *fakeServiceAccountApi) client(t *testing.T) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(server.Close)
	client, err := NewClient("id", "key", "", "us2", server.URL, false)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	return client
}

func (f *fakeServiceAccountApi) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.RequestURI())
	w.Header().Set("Content-Type", "application/json")

	body, _ := io.ReadAll(r.Body)
	var parsed map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &parsed); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	path := strings.TrimPrefix(r.URL.Path, "/v1/serviceAccounts")
	path = strings.TrimPrefix(path, "/")
	switch {
	case path == "" && r.Method == http.MethodPost:
		f.nextID++
		id := fmt.Sprintf("SA%d", f.nextID)
		account := map[string]interface{}{
			"id":       id,
			"name":     parsed["name"],
			"email":    parsed["email"],
			"roleIds":  parsed["roleIds"],
			"isActive": true,
		}
		f.accounts[id] = account
		json.NewEncoder(w).Encode(account)
	case path == "" && r.Method == http.MethodGet:
		list := []interface{}{}
		for _, account := range f.accounts {
			list = append(list, account)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"data": list})
	default:
		account, found := f.accounts[path]
		if !found {
			w.WriteHeader(http.StatusNotFound)
			// Code and message as returned by the real API for a missing service account.
			w.Write([]byte(`{"id":"ABC","errors":[{"code":"serviceaccount:service_account_not_found","message":"Service Account with given ID does not exist."}]}`))
			return
		}
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(account)
		case http.MethodPut:
			for key, value := range parsed {
				account[key] = value
			}
			json.NewEncoder(w).Encode(account)
		case http.MethodDelete:
			delete(f.accounts, path)
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func serviceAccountResourceData(t *testing.T, config map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceSumologicServiceAccount().Schema, config)
}

func TestServiceAccountCreateSendsOnlyCreateFields(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":     "ci-bot",
		"email":    "ci-bot@example.com",
		"role_ids": []interface{}{"R1"},
	})

	if err := resourceSumologicServiceAccountCreate(d, client); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if d.Id() != "SA1" {
		t.Fatalf("expected id SA1, got %q", d.Id())
	}
	if calls := fake.callLog(); len(calls) < 1 || calls[0] != "POST /v1/serviceAccounts" {
		t.Fatalf("expected POST first, got %v", calls)
	}
	for _, call := range fake.callLog() {
		if strings.HasPrefix(call, "PUT") {
			t.Fatalf("an active account must not be updated after create, got %v", fake.callLog())
		}
	}
	if d.Get("is_active") != true {
		t.Fatalf("expected is_active true")
	}
}

func TestServiceAccountCreateInactiveUpdatesAfterCreate(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":      "ci-bot",
		"email":     "ci-bot@example.com",
		"role_ids":  []interface{}{"R1"},
		"is_active": false,
	})

	if err := resourceSumologicServiceAccountCreate(d, client); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Put fetches the ETag with a GET first, so only the writes are compared.
	writes := []string{}
	for _, call := range fake.callLog() {
		if !strings.HasPrefix(call, "GET") {
			writes = append(writes, call)
		}
	}
	if len(writes) != 2 || writes[0] != "POST /v1/serviceAccounts" || writes[1] != "PUT /v1/serviceAccounts/SA1" {
		t.Fatalf("expected POST then PUT, got %v", writes)
	}
	if d.Get("is_active") != false {
		t.Fatalf("expected is_active false in state")
	}
}

func TestServiceAccountReadRemovesMissingFromState(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":     "ci-bot",
		"email":    "ci-bot@example.com",
		"role_ids": []interface{}{"R1"},
	})
	d.SetId("GONE")

	if err := resourceSumologicServiceAccountRead(d, client); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected the resource to be removed from state, id=%q", d.Id())
	}
}

func TestServiceAccountDeleteSendsOnlySetQueryParams(t *testing.T) {
	cases := []struct {
		name   string
		config map[string]interface{}
		want   string
	}{
		{"none", map[string]interface{}{}, "DELETE /v1/serviceAccounts/SA1"},
		{"transfer", map[string]interface{}{"transfer_to": "U9"}, "DELETE /v1/serviceAccounts/SA1?transferTo=U9"},
		{"delete content", map[string]interface{}{"delete_content": true}, "DELETE /v1/serviceAccounts/SA1?deleteContent=true"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := newFakeServiceAccountApi()
			client := fake.client(t)
			config := map[string]interface{}{
				"name":     "ci-bot",
				"email":    "ci-bot@example.com",
				"role_ids": []interface{}{"R1"},
			}
			for k, v := range tc.config {
				config[k] = v
			}
			d := serviceAccountResourceData(t, config)
			if err := resourceSumologicServiceAccountCreate(d, client); err != nil {
				t.Fatalf("create failed: %v", err)
			}

			if err := resourceSumologicServiceAccountDelete(d, client); err != nil {
				t.Fatalf("delete failed: %v", err)
			}

			calls := fake.callLog()
			if got := calls[len(calls)-1]; got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestServiceAccountDeleteIgnoresAlreadyDeleted(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":     "ci-bot",
		"email":    "ci-bot@example.com",
		"role_ids": []interface{}{"R1"},
	})
	d.SetId("GONE")

	if err := resourceSumologicServiceAccountDelete(d, client); err != nil {
		t.Fatalf("deleting an already deleted account must succeed, got %v", err)
	}
}

func TestServiceAccountUpdateChangesEmailInPlace(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":     "ci-bot",
		"email":    "ci-bot@example.com",
		"role_ids": []interface{}{"R1"},
	})
	if err := resourceSumologicServiceAccountCreate(d, client); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	d.Set("email", "new@example.com")
	if err := resourceSumologicServiceAccountUpdate(d, client); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if got := fake.accounts["SA1"]["email"]; got != "new@example.com" {
		t.Fatalf("expected email updated on server, got %v", got)
	}
}

func TestServiceAccountSchemaRejectsConflictingDeleteOptions(t *testing.T) {
	s := resourceSumologicServiceAccount().Schema
	if len(s["transfer_to"].ConflictsWith) == 0 && len(s["delete_content"].ConflictsWith) == 0 {
		t.Fatalf("transfer_to and delete_content must conflict")
	}
	if s["role_ids"].MinItems != 1 {
		t.Fatalf("role_ids must require at least one item")
	}
}

func TestServiceAccountDataSourceByEmail(t *testing.T) {
	fake := newFakeServiceAccountApi()
	client := fake.client(t)
	d := serviceAccountResourceData(t, map[string]interface{}{
		"name":     "ci-bot",
		"email":    "ci-bot@example.com",
		"role_ids": []interface{}{"R1"},
	})
	if err := resourceSumologicServiceAccountCreate(d, client); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	ds := schema.TestResourceDataRaw(t, dataSourceSumologicServiceAccount().Schema, map[string]interface{}{
		"email": "ci-bot@example.com",
	})
	if err := dataSourceSumologicServiceAccountRead(ds, client); err != nil {
		t.Fatalf("data source read failed: %v", err)
	}
	if ds.Id() != "SA1" || ds.Get("name") != "ci-bot" {
		t.Fatalf("unexpected data source result id=%q name=%v", ds.Id(), ds.Get("name"))
	}

	missing := schema.TestResourceDataRaw(t, dataSourceSumologicServiceAccount().Schema, map[string]interface{}{
		"email": "nobody@example.com",
	})
	if err := dataSourceSumologicServiceAccountRead(missing, client); err == nil {
		t.Fatalf("expected an error when no account matches the email")
	}
}

func TestIsNotFoundError(t *testing.T) {
	cases := []struct {
		name string
		body string
		want bool
	}{
		{"service account missing", `{"errors":[{"code":"serviceaccount:service_account_not_found","message":"Service Account with given ID does not exist."}]}`, true},
		{"user missing", `{"errors":[{"code":"user:user_not_found","message":"User with given ID does not exist."}]}`, true},
		{"access denied", `{"errors":[{"code":"serviceaccount:operation_supported_only_for_admins","message":"Operation supported only for an Administrator!"}]}`, false},
		{"invalid transfer target", `{"errors":[{"code":"user:invalid_request","message":"The user you selected to transfer content to was not found"}]}`, false},
		{"not json", `boom`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isNotFoundError(fmt.Errorf("%s", tc.body)); got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
	if isNotFoundError(nil) {
		t.Fatalf("nil must not be a not found error")
	}
}
