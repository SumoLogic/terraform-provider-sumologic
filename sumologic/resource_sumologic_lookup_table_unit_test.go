package sumologic

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestWaitForLookupJob(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantErr     bool
		wantWarning bool
	}{
		{
			name: "success",
			body: `{"jobId":"job1","status":"Success"}`,
		},
		{
			name:        "partial success does not fail the job but returns a warning",
			body:        `{"jobId":"job1","status":"PartialSuccess","warnings":[{"message":"60 rows were dropped"}]}`,
			wantWarning: true,
		},
		{
			name:    "failed",
			body:    `{"jobId":"job1","status":"Failed","errors":[{"code":"x","message":"boom"}]}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := &http.Response{
				Status:     http.StatusText(200),
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(tt.body))),
			}
			client := newTestClient(response)

			warning, err := waitForLookupJob("job1", 5*time.Second, client)
			if (err != nil) != tt.wantErr {
				t.Errorf("waitForLookupJob() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if (warning != "") != tt.wantWarning {
				t.Errorf("waitForLookupJob() warning = %q, wantWarning = %v", warning, tt.wantWarning)
			}
		})
	}
}

func TestHashLookupContent(t *testing.T) {
	if got := hashLookupContent(""); got != "" {
		t.Errorf("hashLookupContent(\"\") = %q, want empty string", got)
	}

	a := hashLookupContent("host,team\na,b\n")
	b := hashLookupContent("host,team\na,b\n")
	if a != b {
		t.Errorf("hashLookupContent is not stable: %q != %q", a, b)
	}
	if len(a) != 64 {
		t.Errorf("hashLookupContent returned %d hex chars, want 64 (sha256)", len(a))
	}

	c := hashLookupContent("host,team\nc,d\n")
	if a == c {
		t.Errorf("hashLookupContent returned the same hash for different content: %q", a)
	}
}

func TestValidateLookupContent(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantErrors   bool
		wantWarnings bool
	}{
		{name: "empty is fine", content: ""},
		{name: "small csv is fine", content: "host,team\na,b\n"},
		{
			name:         "over warn threshold warns",
			content:      strings.Repeat("a", warnContentBytes+1),
			wantWarnings: true,
		},
		{
			name:       "over max threshold errors",
			content:    strings.Repeat("a", maxContentBytes+1),
			wantErrors: true,
		},
		{
			name:         "crlf warns",
			content:      "host,team\r\na,b\r\n",
			wantWarnings: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errs := validateLookupContent(tt.content, "content")
			if (len(errs) > 0) != tt.wantErrors {
				t.Errorf("errors = %v, wantErrors = %v", errs, tt.wantErrors)
			}
			if (len(warnings) > 0) != tt.wantWarnings {
				t.Errorf("warnings = %v, wantWarnings = %v", warnings, tt.wantWarnings)
			}
		})
	}
}

func TestValidateLookupCSV(t *testing.T) {
	fields := []LookupTableField{
		{FieldName: "host", FieldType: "string"},
		{FieldName: "team", FieldType: "string"},
	}
	primaryKeys := []string{"host"}

	tests := []struct {
		name        string
		content     string
		fields      []LookupTableField
		primaryKeys []string
		wantErr     bool
	}{
		{
			name:    "happy path",
			content: "host,team\nweb-1,team-a\n",
			fields:  fields, primaryKeys: primaryKeys,
		},
		{
			name:    "header case differs from fields",
			content: "Host,TEAM\nweb-1,team-a\n",
			fields:  fields, primaryKeys: primaryKeys,
		},
		{
			name:    "header-only is valid",
			content: "host,team\n",
			fields:  fields, primaryKeys: primaryKeys,
		},
		{
			name:    "missing column",
			content: "host\nweb-1\n",
			fields:  fields, primaryKeys: primaryKeys,
			wantErr: true,
		},
		{
			name:    "unknown extra column",
			content: "host,team,extra\nweb-1,team-a,x\n",
			fields:  fields, primaryKeys: primaryKeys,
			wantErr: true,
		},
		{
			name:    "primary key absent from header",
			content: "host,team\nweb-1,team-a\n",
			fields:  fields, primaryKeys: []string{"missing"},
			wantErr: true,
		},
		{
			name:    "ragged row",
			content: "host,team\nweb-1,team-a,extra-value\n",
			fields:  fields, primaryKeys: primaryKeys,
			wantErr: true,
		},
		{
			name:    "unbalanced quote",
			content: "host,team\n\"web-1,team-a\n",
			fields:  fields, primaryKeys: primaryKeys,
			wantErr: true,
		},
		{
			name:    "empty header name",
			content: ",team\nweb-1,team-a\n",
			fields:  fields, primaryKeys: primaryKeys,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLookupCSV(tt.content, tt.fields, tt.primaryKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLookupCSV() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// TestLookupTableCustomizeDiff exercises lookupTableCustomizeDiff itself
// (rather than validateLookupCSV directly), through the same *schema.Resource
// entry point Terraform uses at plan time - covering the NewValueKnown gating
// and the ResourceDiff-to-Go-slice reconstruction that TestValidateLookupCSV,
// operating on hand-built slices, doesn't touch.
func TestLookupTableCustomizeDiff(t *testing.T) {
	fields := []interface{}{
		map[string]interface{}{"field_name": "host", "field_type": "string"},
		map[string]interface{}{"field_name": "team", "field_type": "string"},
	}

	tests := []struct {
		name         string
		raw          map[string]interface{}
		computedKeys []string
		wantErr      bool
	}{
		{
			name: "valid content matching fields",
			raw: map[string]interface{}{
				"content":      "host,team\nweb-1,team-a\n",
				"fields":       fields,
				"primary_keys": []interface{}{"host"},
			},
		},
		{
			name: "content header does not match fields",
			raw: map[string]interface{}{
				"content":      "host,other\nweb-1,x\n",
				"fields":       fields,
				"primary_keys": []interface{}{"host"},
			},
			wantErr: true,
		},
		{
			name: "empty content skips validation",
			raw: map[string]interface{}{
				"content": "",
				"fields":  fields,
			},
		},
		{
			name: "empty fields list skips validation even with mismatched content",
			raw: map[string]interface{}{
				"content": "anything,goes\nhere,too\n",
				"fields":  []interface{}{},
			},
		},
		{
			name: "unknown content at plan time skips validation",
			raw: map[string]interface{}{
				"content": "",
				"fields":  fields,
			},
			computedKeys: []string{"content"},
		},
		{
			name: "unknown fields at plan time skips validation",
			raw: map[string]interface{}{
				"content": "host,other\nweb-1,x\n",
				"fields":  fields,
			},
			computedKeys: []string{"fields"},
		},
	}

	resource := resourceSumologicLookupTable()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(tt.raw)
			config.ComputedKeys = tt.computedKeys

			_, err := resource.Diff(context.Background(), nil, config, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
